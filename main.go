package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

const filename = "links.csv"
const numberOfWorkers = 5

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
	fmt.Printf("Success: %t\n", result.Success)

	fmt.Println(" ========= ")
}

func ensureFileExists() error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)

	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}

	defer file.Close()

	file.WriteString("suggestedName,Link\n")
	file.WriteString("Jack Jones Portfolio,https://www.jackljones.com/\n")
	file.WriteString("Books,https://books.toscrape.com/\n")
	file.WriteString("\"This is a test of getting, JSON\",https://jsonplaceholder.typicode.com/todos/1")

	return nil

}

func worker(index int, jobs <-chan LinkInfo, results chan<- Result) {
	for j := range jobs {
		fmt.Printf("Worked %d is targetting job %s\n", index, j.TargetUrl)

		output, err := GetResult(j)

		if err != nil {
			log.Default().Printf("Worked %d had an error %v getting %s", index, err, j.TargetUrl)

			results <- Result{
				Success:          false,
				OriginalLinkInfo: j,
			}

			continue
		}

		results <- output

	}
}

func main() {

	err := ensureFileExists()

	if err != nil {
		log.Fatalf("%v", err)
	}

	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("%v", err)
	}

	defer file.Close()

	jobs := make(chan LinkInfo, 100)
	results := make(chan Result, 100)

	for a := 1; a <= numberOfWorkers; a++ {
		go worker(a, jobs, results)
	}

	reader := csv.NewReader(file)
	urlCount := 0

	// skip first line

	_, err = reader.Read()

	if err != nil {
		log.Fatalf("%v", err)
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Default().Printf("Skipping due to %v\n", err)
			continue
		}

		linkInfo := LinkInfo{
			SuggestedName: strings.TrimSpace(record[0]),
			TargetUrl:     strings.TrimSpace(record[1]),
		}

		urlCount++

		jobs <- linkInfo

	}

	close(jobs)

	for i := 1; i <= urlCount; i++ {
		output := <-results
		formatOutput(output)
	}

}
