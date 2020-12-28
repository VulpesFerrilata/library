package app_error

import ut "github.com/go-playground/universal-translator"

type WebsocketError interface {
	error
	Message(trans ut.Translator) (string, error)
}
