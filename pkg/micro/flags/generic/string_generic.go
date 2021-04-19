package generic

import (
	"strings"

	"github.com/micro/cli/v2"
	"github.com/pkg/errors"
)

func NewStringGeneric(enums ...string) cli.Generic {
	return &stringGeneric{
		enums: enums,
	}
}

type stringGeneric struct {
	enums    []string
	selected string
}

func (s *stringGeneric) Set(value string) error {
	for _, enum := range s.enums {
		if enum == value {
			s.selected = value
			return nil
		}
	}
	return errors.Errorf("allowed values are %s", strings.Join(s.enums, ","))
}

func (s stringGeneric) String() string {
	return s.selected
}
