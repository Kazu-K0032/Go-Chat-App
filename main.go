package main

import (
	"security_chat_app/app/controllers"
	_ "security_chat_app/config"
)

func main() {
	controllers.StartMainServer()
}
