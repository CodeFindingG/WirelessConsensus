package main

import "math/rand"

func noiseGen() float64 {
	return rand.Float64() * 2
}

func svGen() float64 {
	return rand.Float64() * 10
}
