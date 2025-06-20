# Go Web Analyzer
Go Web Analyzer is a web-based application that analyzes given URL. It detects HTML versions, check whether it contains login form, extracts links and identifies whether links are internal or external and accesible.  This project is integrated with Prometheus metrics and pprof to support performance diagnostics and can be deploy in Docker.

## Project Overview

This project is a Go application which contains,
- Back end : which consist of RESTful API endpoint `/analyze` for analyzing a given URL
- Frond end : which consist of UI that accepts URLs , analyse and provide response with multiple analysis results
- This project can be build and run in a Docker container
- Use Prometheus metrics and pprof to support performance diagnostics

## Prerequisites
- Go (1.24 or higher)
- Git
- Docker (Latest)
- Browser (any latest briowser)

## Setup Instructions

#### 1. Clone the repository
``` git clone https://github.com/AkshithaLiyanage/go-web-analyzer.git```

``` cd go-web-analyzer ```

#### 2. Get dependencies
``` go mod tidy ```

### Run locally
``` go run cmd/server/main.go ```

### Run in Docker

#### 1. Install Docker:
- For Linux:
``` curl -fsSL https://get.docker.com -o get-docker.sh ``` 

``` sudo sh get-docker.sh ```

- For Windows:
``` Download Docker Desktop for Windows from https://www.docker.com/products/docker-desktop ``` 

- For MacOS:
``` Download Docker Desktop for MacOS from https://www.docker.com/products/docker-desktop ``` 

After installation, verify with: ``` docker --version ``` 

#### 2. Build Docker image
``` docker build -t go-web-analyzer . ```

#### 3. Run Docker container
``` docker run -p 8080:8080 -p 6060:6060 go-web-analyzer ``` 

### Run Tests
#### 1. Run All Unit Tests

``` go test ./... ```

#### 2. Displays detailed output for each test.
``` go test -v ./... ```

#### 3. Run Tests with Code Coverage
``` go test -cover ./... ```

#### 4. HTML Coverage Report
``` go test -coverprofile=coverage.out ./... ```
``` go tool cover -html=coverage.out ```


## Usage
- Open [http://localhost:8080](http://localhost:8080) to access the web UI.
- Enter a URL in the input field and submit for analysis.
- Web UI will display the result.
- If the entered URL is Invalid or unreachable, a descriptive error messages will display in web UI.

 ## Main Functionalities
- Web UI :	Serves frontend for URL submission and showing results/errors.
- Analyze API:	`/analyze` returns submitted URL's HTML version, link stats, and login form presence
- Logging : Custom logging using logrus
- Error Handling : return custom errors, with user friendly HTTP response and description
- Metrics	: Prometheus metrics for API request count, latency, etc.
- Profiling :	pprof support for performance diagnostics

## Technologies Used

### Backend (BE) 
- Golang (1.24)
- Prometheus
- pprof

### Frontend (FE)
- HTML
- CSS
- JavaScript

### DevOps
- Docker

## URLs
- Website : [http://localhost:8080](http://localhost:8080)
- Analyze Endpoint : [http://localhost:8080/analyze](http://localhost:8080/analyze)
- Prometheus Metrics : [http://localhost:8080/metrics](http://localhost:8080/metrics)
- pprof Dashboard  : [http://localhost:6060/debug/pprof](http://localhost:6060/debug/pprof)  


 ## API specs

### Analyze URL

- API Endpoint : [http://localhost:8080/analyze](http://localhost:8080/analyze)
- Method : POST
- Request : `url:https://google.com`
- Request Content-Type: ``` multipart/form-data ```
- Response : ```
{
  "title": "Google",
  "headings": {
    "h1": 0,
    "h2": 0,
    "h3": 0,
    "h4": 0,
    "h5": 0,
    "h6": 0
  },
  "html_version": "HTML5",
  "internal_links": [
    "https://google.com/preferences?hl=en"
  ],
  "external_links": [
    "https://maps.google.lk/maps?hl=en&tab=wl"
  ],
  "inaccessible_links": null,
  "is_login_form": false
} ```

- Response content type : application/json

##  Challenges and the approaches took to overcome

1. - Challenge : Errors such as DNS resolution failures, timeouts, and TLS handshake errors returned vague messages or caused crashes.
   - Solution : Added custom error handling by inspecting error messages and returning user-friendly HTTP responses such as Bad Gateway, Gateway Timeout, etc.
2. - Challenge : Many websites do not clearly declare their HTML version, making it difficult to determine reliably.
   - Solution : Analyzed the presence and content of doctype declarations and certain meta tags using simple string checks to identify common versions (e.g., HTML5, HTML 4.01, XHTML).
3. - Challenge : Understanding Go’s syntax, coding styles, libraies as a Java developer
   - Solution : Followed tutorials and other online metirials availbale to learn

## Possible improvements to the project

- Better UI/UX on the frontend.
- Deploy to Kubernetes with proper monitoring stack (Prometheus + Grafana)
- Integrate CI/CD pipeline
- Add support for concurrent URL analyzing
- Implement caching for repeated analysis of the same URL
- Add authentication

## WEB UI
  ![image](https://github.com/user-attachments/assets/fdaf516d-1e92-45d0-bb63-f0e380ff3bab)
