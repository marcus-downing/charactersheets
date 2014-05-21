package main

import "./server"

func main() {
	server.RunTranslator("localhost:9091", 1)
}
