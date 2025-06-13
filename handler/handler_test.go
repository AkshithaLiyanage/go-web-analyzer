package handler

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AkshithaLiyanage/go-web-analyzer/analyzer"
	"github.com/AkshithaLiyanage/go-web-analyzer/logs"
	"github.com/AkshithaLiyanage/go-web-analyzer/model/analysisresult"
	"github.com/stretchr/testify/assert"
)

func createMultipartFormBody(key, value string) (string, *bytes.Buffer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField(key, value)
	writer.Close()
	return writer.FormDataContentType(), body
}

func TestAnalyzeHandler_Success(t *testing.T) {
	logs.Init()

	originalFunc := analyzer.AnalyzeURL
	defer func() { analyzer.AnalyzeURL = originalFunc }()
	analyzer.AnalyzeURL = func(url string) (*analysisresult.AnalysisResult, int, error) {
		return &analysisresult.AnalysisResult{Title: "Mock Title"}, http.StatusCreated, nil
	}

	contentType, body := createMultipartFormBody("url", "http://example.com")
	req := httptest.NewRequest(http.MethodPost, "/analyze", body)
	req.Header.Set("Content-Type", contentType)

	rr := httptest.NewRecorder()
	AnalyzeHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Mock Title")
}

func TestAnalyzeHandler_MethodNotAllowed(t *testing.T) {
	logs.Init()
	req := httptest.NewRequest(http.MethodGet, "/analyze", nil)
	rr := httptest.NewRecorder()

	AnalyzeHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Contains(t, rr.Body.String(), "Method not allowed")
}

func TestAnalyzeHandler_MissingURL(t *testing.T) {
	logs.Init()
	contentType, body := createMultipartFormBody("url", "")
	req := httptest.NewRequest(http.MethodPost, "/analyze", body)
	req.Header.Set("Content-Type", contentType)

	rr := httptest.NewRecorder()
	AnalyzeHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "URL is required")
}

func TestAnalyzeHandler_AnalyzeError(t *testing.T) {
	logs.Init()

	originalFunc := analyzer.AnalyzeURL
	defer func() { analyzer.AnalyzeURL = originalFunc }()
	analyzer.AnalyzeURL = func(url string) (*analysisresult.AnalysisResult, int, error) {
		return nil, http.StatusBadGateway, errors.New("fake error")
	}

	contentType, body := createMultipartFormBody("url", "http://bad-url.com")
	req := httptest.NewRequest(http.MethodPost, "/analyze", body)
	req.Header.Set("Content-Type", contentType)

	rr := httptest.NewRecorder()
	AnalyzeHandler(rr, req)

	assert.Equal(t, http.StatusBadGateway, rr.Code)
	assert.Contains(t, rr.Body.String(), "fake error")
}
