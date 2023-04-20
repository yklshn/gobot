package restapi

import (
	"encoding/json"
	"gorbot/models"
	"log"
	"net/http"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type body struct {
	Value    string `json:"value"`
	Receiver string `json:"receiver"`
}

func ListenHTTP(wg *sync.WaitGroup, ch chan models.Message) {
	defer wg.Done()

	http.HandleFunc("/sendSMS", func(w http.ResponseWriter, r *http.Request) {
		b := &body{}
		if err := json.NewDecoder(r.Body).Decode(b); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}
		msg := models.Message{
			Type:     models.Sms,
			Receiver: b.Receiver,
			Msg:      b.Value,
		}
		ch <- msg
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/sendEmail", func(w http.ResponseWriter, r *http.Request) {
		b := &body{}
		if err := json.NewDecoder(r.Body).Decode(b); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}
		msg := models.Message{
			Type:     models.Email,
			Receiver: b.Receiver,
			Msg:      b.Value,
		}
		ch <- msg
		w.WriteHeader(http.StatusOK)
	})
	type conf struct {
		Port string `yaml:"port"`
	}
	var cfg conf
	if err := cleanenv.ReadConfig("./config/config.yml", cfg); err != nil {
		log.Fatal("Whoops, something has gone wrong!\n", err)
	}

	log.Fatal(http.ListenAndServe(":3333", nil))

}
