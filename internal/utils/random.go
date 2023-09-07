package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomChoice 随机返回
func RandomChoice[T any](list []T) T {
	return list[rand.Intn(len(list))]
}
