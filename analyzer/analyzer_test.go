package analyzer

import (
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/AkshithaLiyanage/go-web-analyzer/logs"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestDetectHtmlVersion(t *testing.T) {
	logs.Init()
	tests := []struct {
		html     string
		expected string
	}{
		{"<!DOCTYPE html>", "HTML5"},
		{"<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">", "HTML 4.01"},
		{"<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 3.2//EN\">", "HTML 3.2"},
		{"<html><head><title>No Doctype</title></head></html>", "Unknown or pre-HTML 4"},
	}

	for _, test := range tests {
		version := detectHTMLVersion(docFromHTML(t, test.html))
		assert.Equal(t, test.expected, version)
	}
}

func docFromHTML(t *testing.T, html string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}
	return doc
}

func TestIsLoginForm(t *testing.T) {
	logs.Init()
	htmlWithLogin := `
		<form>
			<input type="text" name="username"/>
			<input type="password" name="password"/>
			<button type="submit">Login</button>
		</form>`
	htmlWithoutLogin := `<form><input type="text" name="q"/></form>`

	tests := []struct {
		html     string
		expected bool
	}{
		{htmlWithLogin, true},
		{htmlWithoutLogin, false},
	}

	for _, test := range tests {
		isLogin := isLoginForm(docFromHTML(t, test.html))
		assert.Equal(t, test.expected, isLogin)
	}
}

// an integration test that runs the full analysis flow against a mock HTTP server with sample HTML content.
func TestAnalyzeURL_Success(t *testing.T) {
	logs.Init()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head><title>Test Page</title></head>
			<body>
				<h1>Main Heading</h1>
				<form>
					<input type="text" name="username"/>
					<input type="password" name="password"/>
					<input type="submit" value="Login"/>
				</form>
				<a href="/internal">Internal</a>
				<a href="https://external.com">External</a>
			</body>
			</html>
		`))
	}))
	defer ts.Close()

	result, status, err := AnalyzeURL(ts.URL)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, "Test Page", result.Title)
	assert.Equal(t, "HTML5", result.HTMLVersion)
	assert.True(t, result.LoginForm)
	assert.Equal(t, 1, result.Headings["h1"])
	assert.Len(t, result.InternalLinks, 1)
	assert.Len(t, result.ExternalLinks, 1)
}
