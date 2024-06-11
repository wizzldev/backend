package utils

import (
	"math/rand"
	"time"
)

type Random struct {
	random *rand.Rand
}

func NewRandom() Random {
	source := rand.NewSource(time.Now().UnixNano())
	return Random{
		random: rand.New(source),
	}
}

func (r Random) String(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.random.Intn(len(letterRunes))]
	}
	return string(b)
}

func (r Random) Number(low int, hi int) int {
	return low + r.random.Intn(hi-low)
}
