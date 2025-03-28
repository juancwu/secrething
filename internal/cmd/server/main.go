package main

import (
	"github.com/juancwu/konbini/internal/server/config"
)

func main() {
	if err := config.Load(".env"); err != nil {
		panic(err)
	}
}
