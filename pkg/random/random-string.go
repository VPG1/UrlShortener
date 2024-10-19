package random

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(size int) string {
	rand.Seed(time.Now().UnixNano())
	res := make([]byte, size)

	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}

	return string(res)
}
