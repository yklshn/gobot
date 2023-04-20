package telegram

import (
	"fmt"
	"gorbot/models"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TelegramBot(wg *sync.WaitGroup, ch chan models.Message) {
	defer wg.Done()

	//idList := make([]int64, 0, 1)
	Birthday := make(map[string][]int64)
	Stuff := make(map[string][]int64)
	Memes := make(map[string][]int64)
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for {
		select {
		case update := <-updates:
			if update.Message != nil { // If we got a message
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "start":
						m := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Введи /memes или /birthday")
						bot.Send(m)
					case "memes":
						phoneNumber := update.Message.CommandArguments()
						for _, v := range Memes {
							if v[0] == update.Message.Chat.ID {
								break
							}
						}
						Memes[phoneNumber] = append(Memes[phoneNumber], update.Message.Chat.ID)
						// Create a timer that triggers every minute
						ticker := time.NewTicker(10 * time.Second)

						// Set up the timer to call the serverMonitoring function every minute
						go func() {
							for range ticker.C {
								for _, v := range Memes {
									m := sendFile(v[0], "memes")
									bot.Send(m)
								}
							}
						}()
					case "birthday":
						email := update.Message.CommandArguments()
						for _, v := range Birthday {
							if v[0] == update.Message.Chat.ID {
								break
							}
						}
						Birthday[email] = append(Birthday[email], update.Message.Chat.ID)
						// Create a timer that triggers every minute
						ticker := time.NewTicker(10 * time.Second)
						// Set up the timer to call the serverMonitoring function every minute
						go func() {
							for range ticker.C {
								for _, v := range Birthday {
									m := sendFile(v[0], "birthday")
									bot.Send(m)
								}
							}
						}()
					case "stop":
						delRecord(Memes, Birthday, update.Message.Chat.ID)
					}
				}
			}
		case msg := <-ch:
			for _, v := range Stuff[msg.Receiver] {
				m := tgbotapi.NewMessage(v, fmt.Sprintf("%s - %s", msg.Receiver, msg.Msg))
				bot.Send(m)
			}
		}
	}

}
func sendFile(chatID int64, path string) tgbotapi.PhotoConfig {
	pic := RandomPic(path)
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(pic))
	return photo
}
func delRecord(memes map[string][]int64, birthday map[string][]int64, id int64) {
	mx := sync.Mutex{}
	for k, v := range memes {
		if v[0] == id {
			mx.Lock()
			delete(memes, k)
			mx.Unlock()
		}
	}

	for k, v := range birthday {
		if v[0] == id {
			mx.Lock()
			delete(birthday, k)
			mx.Unlock()
		}
	}
}

func RandomPic(evenT string) string {
	pics := make([]string, 0, 10)
	var root string
	if evenT == "birthday" {
		root = "./data/birthday"
	} else {
		root = "./data/memes"
	}

	exts := []string{".jpg", ".jpeg", ".png", ".gif"}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, ext := range exts {

			if strings.HasSuffix(path, ext) {
				pic := path
				pics = append(pics, pic)
			}
		}

		return nil
	})

	if err != nil {
		log.Println(err.Error())
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return pics[r1.Intn(len(pics))]
}
