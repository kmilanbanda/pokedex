package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	cleanText := strings.ToLower(text)
	cleanedWords := strings.Fields(cleanText)	
	return cleanedWords
}

func main() {
    fmt.Println("Hello, World!")
}
