package main

import (
	"net/http"

	"github.com/AkshithaLiyanage/go-web-analyzer/handler"
	"github.com/AkshithaLiyanage/go-web-analyzer/logs"
)

func main() {
	logs.Init()
	logs.Log.Info("Starting Web Analyzer...")
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.HandleFunc("/analyze", handler.AnalyzeHandler)

	logs.Log.Info("Server is running at http:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logs.Log.Fatal("Server failed:", err)
	}
}
