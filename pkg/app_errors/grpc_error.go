package app_errors

import (
	ut "github.com/go-playground/universal-translator"
	"google.golang.org/grpc/status"
)

type GrpcError interface {
	Status(trans ut.Translator) (*status.Status, error)
}
