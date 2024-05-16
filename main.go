package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/muesli/reflow/wordwrap"
	"github.com/skip2/go-qrcode"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const DPI = 96
const CARDHEIGHT = 2
const CARDWIDTH = 3.5

func main() {
	var srcfile string
	var baseurl string
	var baseText string
	flag.StringVar(&srcfile, "source", "", "File to read URL suffixes from")
	flag.StringVar(&baseurl, "base", "https://example.com", "Base url for generated QR Codes")
	flag.StringVar(&baseText, "text", "", "Text to include on card with title")
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
		titleText := strings.Replace(name, "-", " ", -1)
		title := csr.String(titleText)
		imagefile := fmt.Sprintf("%s.png", name)
		outfile := fmt.Sprintf("%s.svg", name)

		outwriter, err := os.Create(outfile)
		if err != nil {
			log.Fatalf("Creating SVG file: %v", err)
		}
		canvas := svg.New(outwriter)

		cardText := title
		if baseText != "" {
			cardText += " "
			cardText += baseText
		}

		cardLines := strings.Split(wordwrap.String(cardText, 18), "\n")

		qrWidth := int(1.5 * DPI)
		qrX := int(CARDWIDTH*DPI) - qrWidth

		canvas.Start(int(3.5*DPI), 2*DPI, "style=\"font-family: Arial\"")
		canvas.Title(fmt.Sprintf("QR Code for %s", title))
		canvas.Image(qrX, 0, qrWidth, int(qrWidth), imagefile)
		canvas.Textlines(24, 32, cardLines, 20, 24, "black", "start")
		canvas.Text(qrX+int(qrWidth/2.0), qrWidth+18, name, "text-anchor: middle; color: black; font-family: Arial; font-size: 10pt")
		canvas.End()
	}

}
