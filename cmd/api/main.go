package main

import (
	"fmt"
	"net/http"

	"github.com/dwivedi-ritik/text-share-be/server"
)

func main() {
	var PORT string = ":8000"
	fmt.Println("Starting server at ", PORT)
	http.ListenAndServe(PORT, server.CreateServer())
}
