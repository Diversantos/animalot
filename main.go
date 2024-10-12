// Simple bot that realize some dialog with animal
//

package main

import (
//	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"os"
	"encoding/json"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	ConfigDebug bool `json:"debug"`	
}

var (
	questions = []string{
		"Гав", "Ваф", "Вуф", "Тяв", "Уф", "Врррррр", "Грррррр",
	}
	
	dog = "\xF0\x9F\x90\xB6"

	standardResponses = map[string][]string{
		"спасибо": {"Пожалуйста!", "Рад помочь!", "Не за что!"},
		"хорошо":  {"Отлично!", "Замечательно!", "Прекрасно!"},
		"окей":    {"Понял вас!", "Хорошо!", "Ясно!"},
	}
)


func main() {
	go http.ListenAndServe("localhost:6060", nil)

	var reply string

	// Config block
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Can't read config.json file", err)
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Can't unmarshal config.json", err)
	}

	// Get SECRET
	secretToken := os.Getenv("ANIMALOT_SECRET_TOKEN")
	if secretToken == "" {
		log.Println("ANIMALOT_SECRET_TOKEN not set")
		return
	}

	// Init bot
	bot, err := tgbotapi.NewBotAPI(secretToken)
	if err != nil {
		log.Panic("Panic, problem with connect to telegram API: ", err)
	}
	bot.Debug = config.ConfigDebug

	log.Printf("Authorized on account %s", bot.Self.UserName)


	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		 }

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)


		// Check for standard responses
		lowercaseText := strings.ToLower(update.Message.Text)
		responses, ok := standardResponses[lowercaseText]
		if ok {
			reply = responses[rand.Intn(len(responses))]
		} else {
			// Ask a random question
			reply = dog + questions[rand.Intn(len(questions))]
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Some error with send: ", err)
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
