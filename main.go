package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/skip2/go-qrcode"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

	csr := cases.Title(language.AmericanEnglish)
	for _, name := range fileLines {
		title := csr.String(name)
		imagefile := fmt.Sprintf("%s.png", name)
		outfile := fmt.Sprintf("%s.svg", name)

		outwriter, err := os.Create(outfile)
		if err != nil {
			log.Fatalf("Creating SVG file: %v", err)
		}
		canvas := svg.New(outwriter)

		canvas.Start(256, 300)
		canvas.Title(fmt.Sprintf("QR Code for %s", title))
		canvas.Image(0, 0, 256, 256, imagefile)
		canvas.Text(128, 257, title, "text-anchor: middle; color: black; font-family: Arial; font-size: 10pt")
		canvas.End()
	}

}
