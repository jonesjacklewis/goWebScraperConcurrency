package main

import (
	"fmt"
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

func main() {
	fmt.Println("Hello World")
}
