package main

import "math/rand"

func randomPanic() {
	if rand.Int()%2 == 0 {
		panic("random panic")
	}
}
