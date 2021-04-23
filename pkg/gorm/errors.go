package gorm

import "github.com/pkg/errors"

var StaleObjectErr = errors.New("attemped to update or delete a stale object")
