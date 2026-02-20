package helper

import (
	"errors"
	"math/rand"
	"time"
)

func RandomDelay(min, max int) {
	delay := min + rand.Intn(max-min+1)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func RandomError(chance int) error {
	if rand.Intn(100) < chance {
		return errors.New("random error")
	}
	return nil
}
