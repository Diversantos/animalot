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
	"io"
	"io/ioutil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	ConfigDebug bool `json:"debug"`
	ConfigLogFile string `json:"logfile"`	
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

	if config.ConfigLogFile != "none" {
		fn := logOutput(config.ConfigLogFile)
		defer fn()
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


func logOutput(logfile string) func() {
	// open file read/write | create if not exist | clear file at open if exists
	f, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	// save existing stdout | MultiWriter writes to saved stdout and file
//	out := os.Stdout
	mw := io.MultiWriter(f)

	// get pipe reader and writer | writes to pipe writer come out pipe reader
	r, w, _ := os.Pipe()

	// replace stdout,stderr with pipe writer | all writes to stdout, stderr will go through pipe instead (fmt.print, log)
	os.Stdout = w
	os.Stderr = w
	
	// writes with log.Print should also write to mw
	log.SetOutput(mw)

	//create channel to control exit | will block until all copies are finished
	exit := make(chan bool)

	go func() {
		// copy all reads from pipe to multiwriter, which writes to stdout and file
		_,_ = io.Copy(mw, r)
		// when r or w is closed copy will finish and true will be sent to channel
		exit <- true
	}()

	// function to be deferred in main until program exits
	return func() {
		// close writer then block on exit channel | this will let mw finish writing before the program exits
		_ = w.Close()
		<-exit
		// close file after all writes have finished
		_ = f.Close()
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
