package models

type Type string

const (
	Sms   Type = "sms"
	Email Type = "email"
)

type Message struct {
	Type     Type
	Receiver string
	Msg      string
}
