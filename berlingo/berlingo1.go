package main

import (
	"github.com/minaguib/berlingo"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	berlingo.Serve(&AI1{})
}
