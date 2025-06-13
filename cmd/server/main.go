package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/AkshithaLiyanage/go-web-analyzer/handler"
	"github.com/AkshithaLiyanage/go-web-analyzer/logs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logs.Init()
	logs.Log.Info("Starting Web Analyzer...")
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.HandleFunc("/analyze", handler.AnalyzeHandler)
	http.Handle("/metrics", promhttp.Handler())

	startPprof()
	logs.Log.Info("Server is running at http:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logs.Log.Fatal("Server failed:", err)
	}

}

func startPprof() {
	go func() {
		logs.Log.Info("Starting pprof server on :6060")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			logs.Log.Fatal("pprof server failed: ", err)
		}
	}()
}
