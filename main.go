package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/skip2/go-qrcode"
)

func main() {
	var srcfile string
	var baseurl string
	flag.StringVar(&srcfile, "source", "", "File to read URL suffixes from")
	flag.StringVar(&baseurl, "base", "https://example.com", "Base url for generated QR Codes")
	flag.Parse()

	if srcfile == "" {
		_, _ = fmt.Println("No source file given")
		return
	}

	readFile, err := os.Open(srcfile)
	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	for _, name := range fileLines {
		desturl := fmt.Sprintf("%s/%s", baseurl, name)
		outname := fmt.Sprintf("%s.png", name)
		qrcode.WriteFile(desturl, qrcode.Medium, 256, outname)
	}

}
