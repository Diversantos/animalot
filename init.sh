#!/bin/bash

PROJECT="animalot"

if [ ! -f "go.mod" ]; then
	go mod init $PROJECT
	go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5@latest
fi

go run $PROJECT


