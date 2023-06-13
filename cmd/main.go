package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/unidoc/unipdf/v3/common"
	pdfcontent "github.com/unidoc/unipdf/v3/contentstream"
	pdfcore "github.com/unidoc/unipdf/v3/core"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func main() {
	// Define default values
	defaultPDFPath := "config/example.pdf"
	defaultChapterRegex := `(\d+\.\d+)`
	defaultRequirementRegex := `(\d+\.\d+\.\d+)`
	defaultHeader := []string{"Chapter", "Requirement", "Description"}
	defaultOutputPath := "output.csv"

	// Define command-line flags
	pdfPathFlag := flag.String("pdfPath", defaultPDFPath, "Path to the PDF file")
	chapterRegexFlag := flag.String("chapterRegex", defaultChapterRegex, "Regex pattern for chapters")
	requirementRegexFlag := flag.String("requirementRegex", defaultRequirementRegex, "Regex pattern for requirements")
	headerFlag := flag.String("header", strings.Join(defaultHeader, ","), "CSV header")
	outputPathFlag := flag.String("outputPath", defaultOutputPath, "Path to the output CSV file")

	// Parse command-line flags
	flag.Parse()

	// Load the PDF file.
	pdfPath := *pdfPathFlag
	reader, err := pdf.NewPdfReaderFromFile(pdfPath)
	if err != nil {
		fmt.Printf("Failed to load PDF file: %v\n", err)
		return
	}

	// Define the regex patterns for chapters and requirements.
	chapterRegex := regexp.MustCompile(*chapterRegexFlag)
	requirementRegex := regexp.MustCompile(*requirementRegexFlag)

	// Create a new CSV file.
	outputPath := *outputPathFlag
	csvFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Failed to create CSV file: %v\n", err)
		return
	}
	defer csvFile.Close()

	// Create a CSV writer.
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write CSV header.
	header := strings.Split(*headerFlag, ",")
	err = csvWriter.Write(header)
	if err != nil {
		fmt.Printf("Failed to write CSV header: %v\n", err)
		return
	}

	// Iterate through each page.
	totalPages := reader.GetNumPages()
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page, err := reader.GetPage(pageNum)
		if err != nil {
			fmt.Printf("Failed to extract content from page %d: %v\n", pageNum, err)
			continue
		}

		extractor := pdfcontent.NewContentExtractor(page)

		// Extract text content from the page.
		textContent, err := extractor.ExtractText()
		if err != nil {
			fmt.Printf("Failed to extract text from page %d: %v\n", pageNum, err)
			continue
		}

		// Search for chapters and requirements using regex patterns.
		chapters := chapterRegex.FindAllStringSubmatch(textContent, -1)
		for _, chapterMatch := range chapters {
			chapter := chapterMatch[0]
			requirements := requirementRegex.FindAllStringSubmatch(textContent, -1)
			for _, requirementMatch := range requirements {
				requirement := requirementMatch[0]
				description := getDescriptionForRequirement(textContent, chapter, requirement)

				// Write to CSV.
				record := []string{chapter, requirement, description}
				err = csvWriter.Write(record)
				if err != nil {
					fmt.Printf("Failed to write record to CSV: %v\n", err)
					continue
				}
			}
		}
	}

	fmt.Println("CSV file created successfully.")
}

// Helper function to extract the description for a requirement.
func getDescriptionForRequirement(textContent, chapter, requirement string) string {
	chapterIndex := strings.Index(textContent, chapter)
	requirementIndex := strings.Index(textContent, requirement)

	if chapterIndex == -1 || requirementIndex == -1 {
		return ""
	}

	startIndex := requirementIndex + len(requirement)
	endIndex := strings.Index(textContent[requirementIndex:], chapter)

	if endIndex == -1 {
		endIndex = len(textContent)
	} else {
		endIndex += requirementIndex
	}

	return strings.TrimSpace(textContent[startIndex:endIndex])
}


