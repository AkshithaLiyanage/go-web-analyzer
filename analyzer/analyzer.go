package analyzer

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/AkshithaLiyanage/go-web-analyzer/logs"
	"github.com/AkshithaLiyanage/go-web-analyzer/model/analysisresult"
	"github.com/PuerkitoBio/goquery"
)

// Analyze recieved URL
var AnalyzeURL = func(rawURL string) (*analysisresult.AnalysisResult, int, error) {

	logs.Log.Info("Starting Analysing URL : ", rawURL)
	// Validate URL
	var urlRegex = regexp.MustCompile(`^https?://([a-zA-Z0-9\-\.]+|\[[0-9a-fA-F:]+\]|\d{1,3}(\.\d{1,3}){3})(:\d+)?(/.*)?$`)

	if !urlRegex.MatchString(rawURL) {
		logs.Log.Error("invalid url")
		return nil, http.StatusBadRequest, errors.New("invalid url")
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		logs.Log.Error("invalid url")
		return nil, http.StatusBadRequest, errors.New("invalid url")
	}
	// set time out value to 8 sec
	client := &http.Client{Timeout: 8 * time.Second}

	// check availablity by fetching url
	resp, err := client.Get(parsedURL.String())
	if err != nil {
		logs.Log.Error("Error fetching URL:", err.Error())
		return nil, getErrorCode(err.Error()), getErrorDesc(err.Error())
	}
	defer resp.Body.Close()

	//read response document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Log.Error(err.Error())
		return nil, http.StatusInternalServerError, err
	}

	//find title
	title := doc.Find("title").First().Text()

	isLoginForm := isLoginForm(doc)

	htmlVersion := detectHTMLVersion(doc)

	headings := calcNumOfHeadings(doc)

	allLinks, internalLinks, externalLinks := findAllLinks(parsedURL, doc)

	inaccessibleLinks := checkAccessibility(allLinks, client)

	//build result
	result := &analysisresult.AnalysisResult{
		Title:             strings.TrimSpace(title),
		Headings:          headings,
		HTMLVersion:       htmlVersion,
		InternalLinks:     internalLinks,
		ExternalLinks:     externalLinks,
		InaccessibleLinks: inaccessibleLinks,
		LoginForm:         isLoginForm,
	}

	return result, http.StatusCreated, nil
}

// check login form is available or not
func isLoginForm(doc *goquery.Document) bool {
	logs.Log.Info("Checking web page contain login form")
	isLogin := false

	doc.Find("form").Each(func(i int, form *goquery.Selection) {
		hasUserField := form.Find("input[type='text'], input[type='email'], input[name='username']").Length() > 0
		hasPassWord := form.Find("input[type='password']").Length() > 0
		hasSubmit := form.Find("input[type='submit'], button[type='submit']").Length() > 0

		if hasUserField && hasPassWord && hasSubmit {
			isLogin = true
		}
	})

	return isLogin
}

func getErrorCode(errStr string) int {
	switch {
	case strings.Contains(errStr, "no such host"):
		return http.StatusBadGateway
	case strings.Contains(errStr, "timeout"):
		return http.StatusGatewayTimeout
	case strings.Contains(errStr, "tls"):
		return http.StatusBadGateway
	default:
		return http.StatusBadGateway
	}
}

func getErrorDesc(errStr string) error {
	switch {
	case strings.Contains(errStr, "no such host"):
		return errors.New("bad gateway: no such host")
	case strings.Contains(errStr, "timeout"):
		return errors.New("gateway timeout: request timed out")
	case strings.Contains(errStr, "tls"):
		return errors.New("bad gateway: TLS handshake failed")
	default:
		return errors.New("bad gateway: failed to fetch target URL")
	}
}

// calculate number of headings
func calcNumOfHeadings(doc *goquery.Document) map[string]int {
	logs.Log.Info("Calculating number of headings")
	const numberOfHeadings = 6
	headings := make(map[string]int)
	for i := 1; i <= numberOfHeadings; i++ {
		key := fmt.Sprintf("h%d", i)
		headings[key] = doc.Find(key).Length()
	}
	return headings
}

// find all links
func findAllLinks(parsedURL *url.URL, doc *goquery.Document) ([]string, []string, []string) {
	logs.Log.Info("Finding all links")
	baseDomain := parsedURL.Hostname()

	var allLinks []string
	var internalLinks []string
	var externalLinks []string

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || strings.HasPrefix(href, "javascript:") {
			return
		}

		linkURL, err := url.Parse(href)
		if err != nil {
			return
		}

		resolvedURL := parsedURL.ResolveReference(linkURL)
		if resolvedURL.Hostname() == baseDomain {
			internalLinks = append(internalLinks, resolvedURL.String())
		} else {
			externalLinks = append(externalLinks, resolvedURL.String())
		}
		allLinks = append(allLinks, resolvedURL.String())
	})
	return allLinks, internalLinks, externalLinks
}

// validate accessibility by checking all the links concurrently
func checkAccessibility(allLinks []string, client *http.Client) []string {
	logs.Log.Info("Validating  accessibility of all links")
	var wg sync.WaitGroup

	resultChan := make(chan string, len(allLinks))

	for _, link := range allLinks {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			res, err := client.Get(l)
			if err != nil || res.StatusCode >= 400 {
				resultChan <- l
				return
			}

		}(link)
	}
	wg.Wait()
	close(resultChan)

	var inaccessible []string
	for link := range resultChan {
		inaccessible = append(inaccessible, link)
	}
	return inaccessible
}

// detect html version
func detectHTMLVersion(doc *goquery.Document) string {
	html, err := doc.Html()
	if err != nil {
		return "Unknown"
	}
	logs.Log.Info("Detecting  HTML version")
	content := strings.ToUpper(html)

	doctypePatterns := map[string]string{
		"<!DOCTYPE HTML>":                        "HTML5",
		"HTML 4.01":                              "HTML 4.01",
		"XHTML 1.0":                              "XHTML 1.0",
		"HTML 3.2":                               "HTML 3.2",
		"-//W3C//DTD HTML 4.01 TRANSITIONAL//EN": "HTML 4.01 Transitional",
		"-//W3C//DTD XHTML 1.0 STRICT//EN":       "XHTML 1.0 Strict",
		"-//W3C//DTD XHTML 1.1//EN":              "XHTML 1.1",
	}

	for pattern, version := range doctypePatterns {
		if strings.Contains(content, pattern) {
			return version
		}
	}
	return "Unknown or pre-HTML 4"
}
