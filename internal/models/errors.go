package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching records found")
