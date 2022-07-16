package main

import (
	"log"
	"os"

	"github.com/stopdiiacity/stopdiiacity-app-go/templates"
	"github.com/stopdiiacity/stopdiiacity-app-go/verify"
)

func main() {
	f, err := os.OpenFile("./public/index.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	templates.WriteGenerate(f, verify.Prefixes())

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
