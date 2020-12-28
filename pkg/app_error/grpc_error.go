package app_error

import (
	ut "github.com/go-playground/universal-translator"
	"google.golang.org/grpc/status"
)

type GrpcError interface {
	error
	Status(trans ut.Translator) (*status.Status, error)
}
