package main

import (
	"edot-test/api/http/handler"
	"edot-test/repository"
	"edot-test/service"
	"log"
	"net/http"
)

func main() {
	// Inisialisasi repository, service, dan handler
	parser := repository.NewCommentParser()
	textService := service.NewTextToJsonService(parser)
	textHandler := handler.NewTextToJsonHandler(textService)

	// Endpoint untuk text to JSON
	http.HandleFunc("/api/v1/text-to-zed", textHandler.ConvertText)

	// Jalankan server HTTP
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
