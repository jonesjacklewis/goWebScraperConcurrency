package main

import (
	"fmt"
	"log"
	"time"
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
}

func GetResult(linkInfo LinkInfo) (Result, error) {
	return Result{}, fmt.Errorf("error happened")
}

func main() {
	testSuggestedName := "Jack Jones Portfolio"
	testTargetUrl := "https://www.jackljones.com/"

	linkInfo := LinkInfo{
		SuggestedName: testSuggestedName,
		TargetUrl:     testTargetUrl,
	}

	_, err := GetResult(linkInfo)

	if err != nil {
		log.Fatalf("%v", err)
	}
}
