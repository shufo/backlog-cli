package util

import (
	"fmt"
	"math/rand"
	"time"
)

func Genrate7DigitsRandomNumber() string {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(10000000)

	firstPart := num / 1000  // take the first three digits of the number
	secondPart := num % 1000 // take the last three digits of the number

	return fmt.Sprintf("%04d-%03d", firstPart, secondPart)
}
