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

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	ConfigDebug bool `json:"debug"`	
}

type Phrases struct {
	prase string
	multi bool
}

type Animal struct {
	name string
	emoji string
	phrases []Phrases
}

type Animals struct {
	animals []Animal
}

var (
	standardResponses = map[string][]string{
		"спасибо": {"Пожалуйста!", "Рад помочь!", "Не за что!"},
		"хорошо":  {"Отлично!", "Замечательно!", "Прекрасно!"},
		"окей":    {"Понял вас!", "Хорошо!", "Ясно!"},
	}


)


func main() {
	var reply string
	var config Config
	animals := Animals{}

	// Config block
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Can't read config.json file", err)
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Can't unmarshal config.json", err)
	}

	// Get Animals
	animalsFile, err := ioutil.ReadFile("animals.json")
	if err != nil {
		log.Fatal("Can't read animals.json file", err)
	}
	log.Println(animalsFile)
	err = json.Unmarshal(animalsFile, &animals)
	if err != nil {
		log.Fatal("Can't unmarshal animals.json", err)
	}

	log.Println(animals)

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
//			reply = dog + questions[rand.Intn(len(questions))]
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
