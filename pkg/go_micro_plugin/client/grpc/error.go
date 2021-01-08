package grpc

import (
	"github.com/micro/go-micro/v2/errors"
	"google.golang.org/grpc/status"
)

func microError(err error) error {
	// no error
	switch err {
	case nil:
		return nil
	}

	if verr, ok := err.(*errors.Error); ok {
		return verr
	}

	// grpc error
	if _, ok := status.FromError(err); ok {
		return err
	}

	// fallback
	return errors.InternalServerError("go.micro.client", err.Error())
}
