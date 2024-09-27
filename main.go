package main

import (
	"edot-test/api/http/handler"
	"edot-test/repository"
	"edot-test/service"
	"log"
	"net/http"
)

func main() {
	parser := repository.NewCommentParser()
	textService := service.NewTextToJsonService(parser)
	textHandler := handler.NewTextToJsonHandler(textService)

	http.HandleFunc("/api/v1/text-to-zed", textHandler.ConvertText)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
