package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
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
	// Define the file and dir flags
	fileName := flag.String("file", "", "Name of the input .txt file")
	dirName := flag.String("dir", ".", "The directory where the text files are located")

	// Parse the flags
	flag.Parse()

	// If the dir flag is provided, list all .txt files in the given directory
	if *fileName == "" {
		fmt.Printf("List of .txt files in directory '%s':\n", *dirName)
		err := filepath.Walk(*dirName, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
				// Generate HTML page for each .txt file found
				generateHTMLPage(filepath.Join(*dirName, info.Name()))
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// If both file and dir flags are provided, exit with an error
	if *fileName != "" && *dirName != "" {
		fmt.Println("Error: Please provide either the 'file' flag or the 'dir' flag, not both.")
		os.Exit(1)
	}

	// Construct the full path of the text file
	filePath := fmt.Sprintf("%s/%s", *dirName, *fileName)

	// Check if the provided file exists
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' does not exist in directory '%s'\n", *fileName, *dirName)
			os.Exit(1)
		}
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Generate HTML page for the specified .txt file
	generateHTMLPage(filePath)
}

// generateHTMLPage generates an HTML page for the given text file
func generateHTMLPage(filePath string) {
	// Read the contents of the input text file
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Trim the file extension from the file name
	textFileName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	// Create a Page instance with relevant information
	page := Page{
		TextFilePath: filePath,
		TextFileName: textFileName,
		HTMLPagePath: textFileName + ".html",
		Content:      string(fileContents),
	}

	// Create a new template in memory named "template.tmpl"
	tmpl := template.Must(template.New("template.tmpl").ParseFiles("template.tmpl"))

	// Create a new HTML file with the appropriate name
	htmlFile, err := os.Create(page.HTMLPagePath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer htmlFile.Close()

	// Execute the template with the Page instance's data and write to the HTML file
	err = tmpl.Execute(htmlFile, page)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Print a message to the console indicating the successful creation of the HTML file
	fmt.Printf("HTML template written to %s\n", page.HTMLPagePath)
}
