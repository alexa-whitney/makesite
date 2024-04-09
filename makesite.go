package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v4"
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

	// Timer starts
	startTime := time.Now()

	// Counters for pages generated and total size
	pagesGenerated := 0
	totalSize := int64(0)

	// If the dir flag is provided, list all .txt files in the given directory
	if *fileName == "" {
		// Print the list of .txt files in the directory
		fmt.Println(aurora.Bold(aurora.BgBrightRed(fmt.Sprintf("List of .txt files in directory '%s':", *dirName))))
		// Walk the directory and list all .txt files
		err := filepath.Walk(*dirName, func(path string, info os.FileInfo, err error) error {
			// Check if the file is a .txt file
			if err != nil {
				return err
			}
			// If the file is a .txt file, generate an HTML page
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
				size, err := generateHTMLPage(filepath.Join(*dirName, info.Name()))
				if err != nil {
					return err
				}
				pagesGenerated++                                                                                                          // Increment the counter
				totalSize += size                                                                                                         // Add the size of the generated HTML page to the total size
				fmt.Printf("%s %s\n", aurora.Bold(aurora.Green("HTML template written to")), aurora.Italic(aurora.BgBrightMagenta(path))) // Print the path of the generated HTML page
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	} else { // Else branch to handle single file conversion
		filePath := fmt.Sprintf("%s/%s", *dirName, *fileName)
		_, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Error: File '%s' does not exist in directory '%s'\n", *fileName, *dirName)
				os.Exit(1)
			}
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		size, err := generateHTMLPage(filePath)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		pagesGenerated++
		totalSize += size
	}

	// Calculate the total size in kilobytes with one significant digit after the decimal
	totalSizeKb := float64(totalSize) / 1024 // Convert bytes to kilobytes

	// Duration of execution
	duration := time.Since(startTime)

	// Success message parts with color and style
	successPart := aurora.Bold(aurora.Green("Success!")).String() // The word "Success!" is bold and green
	countPart := aurora.Bold(pagesGenerated).String()             // The count is bold
	pagesPart := "pages"
	sizePart := fmt.Sprintf("(%.1fKB total)", totalSizeKb)
	timePart := fmt.Sprintf("in %.2f seconds.", duration.Seconds())

	// Combine parts with different styles
	successMessage := fmt.Sprintf("%s Generated %s %s %s %s",
		successPart,
		countPart,
		aurora.Bold(aurora.Cyan(pagesPart)).String(),
		aurora.Bold(aurora.BrightYellow(sizePart)).String(),
		aurora.Bold(aurora.BrightBlue(timePart)).String(),
	)

	// Print the formatted success message
	fmt.Println(successMessage)
}

// generateHTMLPage reads the content of a text file and generates an HTML page
func generateHTMLPage(filePath string) (int64, error) {
	// Read the content of the text file
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	// Create a Page struct with the file content
	textFileName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	page := Page{
		TextFilePath: filePath,               // The full path to the text file
		TextFileName: textFileName,           // The name of the text file without the extension
		HTMLPagePath: textFileName + ".html", // The name of the HTML file to be generated
		Content:      string(fileContents),   // The content of the text file
	}

	// Create an HTML file from the template
	tmpl := template.Must(template.New("template.tmpl").ParseFiles("template.tmpl"))
	// Create the HTML file
	htmlFile, err := os.Create(page.HTMLPagePath)
	if err != nil {
		return 0, err
	}
	// Close the file when the function returns
	defer htmlFile.Close()

	err = tmpl.Execute(htmlFile, page)
	if err != nil {
		return 0, err
	}

	// Get the file info to calculate the size
	fileInfo, err := htmlFile.Stat()
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}
