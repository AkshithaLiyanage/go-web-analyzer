package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AkshithaLiyanage/go-web-analyzer/analyzer"
	"github.com/AkshithaLiyanage/go-web-analyzer/logs"
)

// handle requests that are coming to /analyze
func AnalyzeHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		logs.Log.Error("Method not allowed : ", http.StatusMethodNotAllowed)
		writeJSONError(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		logs.Log.Error("Failed to parse form : ", err)
		writeJSONError(res, "Failed to parse form", http.StatusBadRequest)
		return
	}

	url := req.FormValue("url")
	if url == "" {
		logs.Log.Error("URL is empty")
		writeJSONError(res, "URL is required", http.StatusBadRequest)
		return
	}

	logs.Log.Info("Received URL for analysis : ", url)

	analysisResult, httpCode, err := analyzer.AnalyzeURL(url)
	if err != nil {
		logs.Log.Error("Error Analysing URL : ", err)
		writeJSONError(res, err.Error(), httpCode)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(analysisResult)
	logs.Log.Info("Analyzis Completed")
}

func writeJSONError(res http.ResponseWriter, message string, statusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(map[string]any{
		"status_code": statusCode,
		"error":       message,
	})

}
