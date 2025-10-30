package main

import (
	"fmt"
	"math/rand"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateTextId() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = alpha[rand.Intn(len(alpha))]
	}
	return fmt.Sprintf("%s%d", b, rand.Intn(1000000000))
}
