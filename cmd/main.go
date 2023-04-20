package main

import (
	"gorbot/internal/restapi"
	"gorbot/internal/telegram"
	"gorbot/models"
	"sync"
)

func main() {
	ch := make(chan models.Message)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go telegram.TelegramBot(wg, ch)
	go restapi.ListenHTTP(wg, ch)
	wg.Wait()
}
