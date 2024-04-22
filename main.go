package main

import (
	wp_api "DreamVerseCloud/wp_library/wp_http/api"
	"log"
)

func main() {
	log.Println("正在初始化……")
	log.Println("start on http://127.0.0.1:8080/")
	wp_api.StartServer()
}
