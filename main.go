package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type LinkInfo struct {
	SuggestedName string
	TargetUrl     string
}

type Result struct {
	Title            string
	StatusCode       int
	RequestDuration  time.Duration
	OriginalLinkInfo LinkInfo
	Success          bool
}

func GetResult(linkInfo LinkInfo) (Result, error) {

	if strings.TrimSpace(linkInfo.SuggestedName) == "" {
		return Result{
			OriginalLinkInfo: linkInfo,
			Success:          false,
		}, fmt.Errorf("provided LinkInfo had an empty SuggestedName")
	}

	if strings.TrimSpace(linkInfo.TargetUrl) == "" {
		return Result{
			OriginalLinkInfo: linkInfo,
			Success:          false,
		}, fmt.Errorf("provided TargetUrl had an empty SuggestedName")
	}

	// prevent time out of resource takes long time,
	client := &http.Client{Timeout: 10 * time.Second}

	start := time.Now()

	response, err := client.Get(linkInfo.TargetUrl)

	if err != nil {
		return Result{
			OriginalLinkInfo: linkInfo,
			Success:          false,
		}, err
	}

	defer response.Body.Close()

	timeTaken := time.Since(start)

	if response.StatusCode != http.StatusOK {
		return Result{
			Title:            linkInfo.SuggestedName,
			StatusCode:       response.StatusCode,
			RequestDuration:  timeTaken,
			OriginalLinkInfo: linkInfo,
			Success:          false,
		}, nil
	}

	contentType := response.Header.Get("Content-Type")
	contentType = strings.ToLower(contentType)

	if !strings.Contains(contentType, "html") {
		// this is still a valid outcome, just not HTML so can't extract further info
		return Result{
			Title:            linkInfo.SuggestedName,
			StatusCode:       response.StatusCode,
			RequestDuration:  timeTaken,
			OriginalLinkInfo: linkInfo,
			Success:          true,
		}, nil
	}

	reader, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		// could still be valid so not passing error downstream, but should log
		log.Default().Printf("%v", err)
		return Result{
			Title:            linkInfo.SuggestedName,
			StatusCode:       response.StatusCode,
			RequestDuration:  timeTaken,
			OriginalLinkInfo: linkInfo,
			Success:          true,
		}, nil
	}

	title := strings.TrimSpace(reader.Find("title").First().Text())

	if len(title) == 0 {
		title = strings.TrimSpace(reader.Find("h1").First().Text())
	}

	if len(title) == 0 {
		title = linkInfo.SuggestedName
	}

	return Result{
		Title:            title,
		StatusCode:       response.StatusCode,
		RequestDuration:  timeTaken,
		OriginalLinkInfo: linkInfo,
		Success:          true,
	}, nil
}

func formatOutput(result Result) {
	fmt.Println(" ========= ")

	fmt.Printf("Title: %s\n", result.Title)
	fmt.Printf("Status Code: %d\n", result.StatusCode)
	fmt.Printf("Request Duration: %d Milliseconds\n", result.RequestDuration.Milliseconds())
	fmt.Printf("Original Suggested Name: %s\n", result.OriginalLinkInfo.SuggestedName)
	fmt.Printf("Target URL: %s\n", result.OriginalLinkInfo.TargetUrl)
	fmt.Printf("Title: %t\n", result.Success)

	fmt.Println(" ========= ")
}

func main() {
	testSuggestedName := "Jack Jones Portfolio"
	testTargetUrl := "https://www.jackljones.com/"

	linkInfo := LinkInfo{
		SuggestedName: testSuggestedName,
		TargetUrl:     testTargetUrl,
	}

	result, err := GetResult(linkInfo)

	if err != nil {
		log.Fatalf("%v", err)
	}

	formatOutput(result)
}
