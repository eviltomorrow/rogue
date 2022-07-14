package util

import (
	"math/rand"
	"time"
)

func GenRandInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}
