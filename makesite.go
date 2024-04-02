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
	// Define a flag for the input file name
	fileName := flag.String("file", "first-post.txt", "Name of the input .txt file")
	// Parse the flags - required to access the flag values
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
	// Handle any errors that occur during file reading
	if err != nil {
		// Print the error message and exit the program
		panic(err)
	}

	// Assign the content of the text file to the Page instance
	page.Content = string(fileContents)

	// Create a new template in memory named "template.tmpl"
	tmpl := template.Must(template.New("template.tmpl").ParseFiles("template.tmpl"))

	// Create a new HTML file with the appropriate name
	htmlFile, err := os.Create(page.HTMLPagePath)
	// Handle any errors that occur during file creation
	if err != nil {
		panic(err)
	}
	// Defer the closing of the HTML file until the function completes
	defer htmlFile.Close()

	// Execute the template with the Page instance's data and write to the HTML file
	err = tmpl.Execute(htmlFile, page)
	if err != nil {
		panic(err)
	}

	// Print a message to the console indicating the successful creation of the HTML file
	fmt.Printf("HTML template written to %s\n", page.HTMLPagePath)
}
