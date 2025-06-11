package analyzer

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/AkshithaLiyanage/go-web-analyzer/logs"
	"github.com/AkshithaLiyanage/go-web-analyzer/model/analysisresult"
	"github.com/PuerkitoBio/goquery"
)

// Analyze recieved URL
func AnalyzeURL(rawURL string) (*analysisresult.AnalysisResult, int, error) {
	logs.Init()
	logs.Log.Info("Starting Analysing URL : ", rawURL)
	// Validate URL
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		logs.Log.Error("Invalid URL")
		return nil, http.StatusBadRequest, errors.New("invalid url")
	}
	// set time out value to 5 sec
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(parsedURL.String())
	if err != nil {
		logs.Log.Error(err.Error())
		return nil, resp.StatusCode, err
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
	//check login form is available or not
	isLoginForm := doc.Find("input[type='password']").Length() > 0

	// detect html version
	var htmlVersion string
	html, err := doc.Html()
	if err != nil {
		htmlVersion = "Unknown"
	} else {
		htmlVersion = detectHTMLVersion(html)
	}
	//build result
	result := &analysisresult.AnalysisResult{
		Title:             strings.TrimSpace(title),
		Headings:          make(map[string]int),
		HTMLVersion:       htmlVersion,
		InternalLinks:     []string{},
		ExternalLinks:     []string{},
		InaccessibleLinks: []string{},
		LoginForm:         isLoginForm,
	}

	calcNumOfHeadings(result, doc)

	allLinks := findAllLinks(parsedURL, doc, result)

	result.InaccessibleLinks = validateUrls(allLinks, client)

	return result, http.StatusCreated, nil
}

// calculate number of headings
func calcNumOfHeadings(result *analysisresult.AnalysisResult, doc *goquery.Document) {
	logs.Log.Info("Calculating number of headings")
	const numberOfHeadings = 6
	for i := 1; i <= numberOfHeadings; i++ {
		tag := fmt.Sprintf("h%d", i)
		result.Headings[tag] = doc.Find(tag).Length()
	}
}

// find all links
func findAllLinks(parsedURL *url.URL, doc *goquery.Document, result *analysisresult.AnalysisResult) []string {
	logs.Log.Info("Finding all links")
	baseDomain := parsedURL.Hostname()

	var allLinks []string

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
			result.InternalLinks = append(result.InternalLinks, resolvedURL.String())
		} else {
			result.ExternalLinks = append(result.ExternalLinks, resolvedURL.String())
		}
		allLinks = append(allLinks, resolvedURL.String())
	})
	return allLinks
}

// validate accessibility by checking all the links concurrently
func validateUrls(allLinks []string, client *http.Client) []string {
	logs.Log.Info("Validating  accessibility")
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
func detectHTMLVersion(html string) string {
	logs.Log.Info("Detecting  HTML version")
	lowered := strings.ToLower(html)

	switch {
	case strings.Contains(lowered, "<!doctype html>"):
		return "HTML5"
	case strings.Contains(lowered, "html 4.01"):
		return "HTML 4.01"
	case strings.Contains(lowered, "html 4.0"):
		return "HTML 4.0"
	case strings.Contains(lowered, "html 3.2"):
		return "HTML 3.2"
	case strings.Contains(lowered, "html 2.0"):
		return "HTML 2.0"
	case strings.Contains(lowered, "xhtml 1.1"):
		return "XHTML 1.1"
	case strings.Contains(lowered, "xhtml 1.0"):
		return "XHTML 1.0"
	default:
		return "Unknown or pre-HTML 4"
	}
}
