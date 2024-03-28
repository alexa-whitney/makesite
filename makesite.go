package main

import (
	"html/template"
	"os"
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
	// Read the contents of first-post.txt
	content, err := os.ReadFile("first-post.txt")
	if err != nil {
		panic(err)
	}

	// Create a Page instance with relevant information
	page := Page{
		TextFilePath: "first-post.txt",
		TextFileName: "first-post",
		HTMLPagePath: "first-post.html",
		Content:      string(content),
	}

	// Create a new template in memory named "template.tmpl"
	t := template.Must(template.New("template.tmpl").ParseFiles("template.tmpl"))

	// Create a new, blank HTML file
	newFile, err := os.Create("first-post.html")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	// Executing the template injects the Page instance's data,
	// allowing us to render the content of our text file.
	// Furthermore, upon execution, the rendered template will be
	// saved inside the new file we created earlier.
	err = t.Execute(newFile, page)
	if err != nil {
		panic(err)
	}

	println("HTML template written to first-post.html")
}
