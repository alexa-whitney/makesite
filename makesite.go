package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"strings"
)

// Page holds all the information needed to generate a new HTML page
// from a text file on the filesystem.
type Page struct {
	TextFilePath string
	TextFileName string
	HTMLPagePath string
	Content      string
}

func main() {
	// Define a flag named "file" to specify the name of the input text file
	fileName := flag.String("file", "first-post.txt", "Name of the input .txt file")
	flag.Parse()

	// Trim the file extension from the file name
	textFileName := strings.TrimSuffix(*fileName, ".txt")

	// Create a Page instance with relevant information
	page := Page{
		TextFilePath: *fileName,
		TextFileName: textFileName,
		HTMLPagePath: textFileName + ".html",
		Content:      "",
	}

	// Read the contents of the input text file
	fileContents, err := os.ReadFile(page.TextFilePath)
	if err != nil {
		panic(err)
	}

	// Assign the content of the text file to the Page instance
	page.Content = string(fileContents)

	// Create a new template in memory named "template.tmpl"
	tmpl := template.Must(template.New("template.tmpl").ParseFiles("template.tmpl"))

	// Create a new HTML file with the appropriate name
	htmlFile, err := os.Create(page.HTMLPagePath)
	if err != nil {
		panic(err)
	}
	defer htmlFile.Close()

	// Execute the template with the Page instance's data and write to the HTML file
	err = tmpl.Execute(htmlFile, page)
	if err != nil {
		panic(err)
	}

	fmt.Printf("HTML template written to %s\n", page.HTMLPagePath)
}
