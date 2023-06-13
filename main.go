package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func main() {
	// Define command-line flags.
	pdfPath := flag.String("pdfPath", "sample.pdf", "Specifies the path to the PDF file.")
	header := flag.String("header", "Chapter,Requirement,Description", "Specifies the CSV header.")
	outputPath := flag.String("outputPath", "output.csv", "Specifies the path to the output CSV file.")
	flag.Parse()

	// Read the input PDF file.
	file, err := os.Open(*pdfPath)
	if err != nil {
		log.Fatalf("Error opening PDF file: %v", err)
	}
	defer file.Close()

	// Read the content of the PDF file.
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading PDF file: %v", err)
	}

	// Load the PDF content using pdfcpu.
	pdf, err := api.ReadPDF(content, nil)
	if err != nil {
		log.Fatalf("Error loading PDF: %v", err)
	}

	// Create a new CSV file for writing.
	csvFile, err := os.Create(*outputPath)
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	defer csvFile.Close()

	// Create a new CSV writer.
	csvWriter := csv.NewWriter(csvFile)

	// Write the CSV header.
	err = csvWriter.Write(strings.Split(*header, ","))
	if err != nil {
		log.Fatalf("Error writing CSV header: %v", err)
	}

	// Define the regex pattern for matching integers in sequential order (e.g., 1.1.1).
	integerPattern := regexp.MustCompile(`\b(\d+\.\d+\.\d+)\b`)

	// Iterate through the pages and search for matching integers.
	for _, page := range pdf.Pages {
		// Extract the text content from the page.
		content, err := api.ExtractText(content, []int{page.Number}, nil)
		if err != nil {
			log.Printf("Error extracting text from page %d: %v", page.Number, err)
			continue
		}

		// Find all matching integers in the content.
		matches := integerPattern.FindAllStringSubmatch(content, -1)

		// Write the matching integers to the CSV file.
		for _, match := range matches {
			record := []string{match[1]}
			err = csvWriter.Write(record)
			if err != nil {
				log.Printf("Error writing to CSV file: %v", err)
				break
			}
		}
	}

	// Flush any buffered data and check for errors.
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Fatalf("Error flushing CSV writer: %v", err)
	}

	fmt.Println("Conversion completed successfully.")
}
