package helper

import (
	"math/rand"
	"time"
)

func GenerateRandInt(length int) string {
	var result = []rune{}
	rand.Seed(time.Now().UnixMicro())
	for i := 0; i < length; i++ {
		result = append(result, rune(rand.Intn(10)+48))
	}
	return string(result)
}
