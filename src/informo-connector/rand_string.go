package main

import (
	"math/rand"
	"time"
)

var signs = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ^!0123456789?./:;,$ù£%*µ&é'(-è_çà)='")
var onlyLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int, all bool) string {
	var letters []rune
	if all {
		letters = signs
	} else {
		letters = onlyLetters
	}
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
