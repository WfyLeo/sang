package main

import (
	"log"
	"newGoSang/config"
)

func main() {
	value := config.File.MustValue("login_server", "host", "127.0.0.1")
	log.Println(value)
}
