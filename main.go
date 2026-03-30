package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"
	tsize "github.com/kopoli/go-terminal-size"
	graphics "DocumentationWriter/graphics"
)

type UserInputs []map[string]string
type Category []map[string]UserInputs

type Template map[string][]Content
type Content struct {
	// name : string
	// contents : [any]
	Name     string
	Contents []interface{}
}

func Print(text string) {
	fmt.Print(text)
}

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Throw(err any) {
	if err != nil {
		panic(err)
	}
}

func Get() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	newInput := input[:len(input)-2] // remove \r\n

	Throw(err)
	return newInput
}

func Contains(str string, substr string) bool {
	if strings.Contains(str, substr) {
		return true
	} else {
		return false
	}
}

func Bigger(var1 string, var2 string) bool {
	if len(var1) >= len(var2) {
		return true
	} else {
		return false
	}
}

func Smaller(var1 string, var2 string) bool {
	if len(var1) <= len(var2) {
		return true
	} else {
		return false
	}
}

func Convention(input string, Convention string, conventionType string, oldConvention string) (string, string) {
	if Bigger(input, Convention) && Contains(input, Convention) {
		return input[2:], conventionType
	}

	return input, oldConvention
}

func ContentInputs(Document *docx.RootDoc, JSONdata Template, Width int, Offset int, selectedTemplate string) []Content {
	/*
		- inputsList = list of all inputs before "next category"
		- categoryList = the name of the category + all the inputs of the category (inputsList)
	*/

	var input string
	var conType string

	contents, ok := JSONdata[selectedTemplate] // or "template1", whatever your key is
	if !ok {
		fmt.Println("Template not found")
	}

	// contentsList
	for parsed := range contents {
		var UserInputs UserInputs
		Document.AddHeading(contents[parsed].Name, 1)
		for {
			// Document.AddHeading(contents[parsed].Name, 1)
			inputText := ""
			if len(input) >= 20 {
				inputText = ("Contents: lots of text Added!")
			} else if len(input) <= 1 {
				inputText = ("Nothing added yet.")
			} else {
				inputText = ("Contents: " + input)
			}

			graphics.GenerateBox(Width, Offset, "Current Category: "+contents[parsed].Name, func() {
				graphics.PrintTextLine(Width, Offset, 4, "Type 'help' for conventions.")
				graphics.PrintLine(Width, Offset, " ", "│", "│")
				graphics.PrintTextLine(Width, Offset, 2, "Type: "+conType)
				graphics.PrintTextLine(Width, Offset, 2, inputText)
			})
			conType = "" // reset previous Convention

			input = Get() // user input

			input, conType = Convention(input, "h=", "header", conType)
			input, conType = Convention(input, "p=", "paragraph", conType)
			input, conType = Convention(input, "b=", "bulletpoint", conType)

			// List :
			if conType != "" {
				// inserting conventions into docx file.
				if conType == "header" {
					Document.AddHeading(input, 2)
				}
				if conType == "paragraph" {
					Document.AddParagraph(input)
				}
				if conType == "bulletpoint" {
					Document.AddParagraph("- " + input).Style("List Bullet")
				}

				UserInputs = append(UserInputs,
					map[string]string{"type": conType, "contents": input})
			}
			if input == "n" {
				contents[parsed].Contents = append(contents[parsed].Contents, UserInputs)
				break
			}

			if input == "help" || input == "Help" {
				Clear()
				Print(
					"│\n" +
						"│ 'enter'-> Return.\n" +
						"│ 'h=' -> Header Element.\n" +
						"│ 'p=' -> Paragraph Element\n" +
						"│ 'b=' -> Bullet Element.\n",
				)
				input = Get()
			}
			Clear()
		}
		Clear()
	}
	fmt.Println("Generating .docx file!")

	return contents
}

func SelectTemplate(data Template) string {
	var intInput int
	var err error
	for {
		Clear()
		// Show all templates available.
		keys := make([]string, 0, len(data)) // creates a list with allocated memory of len(data)
		for k := range data {
			keys = append(keys, k)
		}

		sort.Strings(keys) // sort da keys

		for i, k := range keys {
			fmt.Println("(", i+1, ").", k)
		}

		// Input handling
		Print("\nTemplate: ")
		input := Get()
		intInput, err = strconv.Atoi(input) // Convert string > integer
		if err == nil {                     // If there's no error select the template
			// check if input is equal to existing template.
			for i, k := range keys {
				if i+1 == intInput {
					Clear()
					return k
				}
			}
		}
	}
}

func GetDocumentInformation(Width int, Offset int, Document *docx.RootDoc) string {
	graphics.GenerateBox(Width, Offset, "Enter the Filename.", func() {
		graphics.PrintTextLine(Width, Offset, 2, "Filename: None")
		graphics.PrintTextLine(Width, Offset, 2, "Title:    None")

		graphics.PrintLine(Width, Offset, " ", "│", "│")
		graphics.PrintTextLine(Width, Offset, 4, "This will be the filename + '.docx'")
	})

	input := Get()

	DocName := input
	Clear()
	graphics.GenerateBox(Width, Offset, "Enter the Document Title", func() {
		graphics.PrintTextLine(Width, Offset, 2, "Filename: "+DocName)
		graphics.PrintTextLine(Width, Offset, 2, "Title:    None")
		graphics.PrintLine(Width, Offset, " ", "│", "│")
		graphics.PrintTextLine(Width, Offset, 4, "Title will be inside the file!")
	})

	input = Get()
	Clear()
	Document.AddHeading(input, 0)

	graphics.GenerateBox(Width, Offset, "Who's it made by?", func() {
		graphics.PrintTextLine(Width, Offset, 2, "Filename: "+DocName)
		graphics.PrintTextLine(Width, Offset, 2, "Title:    "+input)
		graphics.PrintLine(Width, Offset, " ", "│", "│")
		graphics.PrintTextLine(Width, Offset, 4, "This will be inside the file!")
	})
	input = Get()
	Clear()
	Document.AddParagraph("Gemaakt door: " + input)
	return DocName
}


func main() {
	Clear()

	// TODO: Width can't be Smaller than 100

	var s tsize.Size
	s, err := tsize.GetSize()
	Throw(err)

	var offset int = 10

	// STEP 0 | Get TITLE, MADE BY & FILE TITLE
	document, err := godocx.NewDocument()
	Throw(err)

	docName := GetDocumentInformation(s.Width, offset, document)

	// STEP 1 | FILL INPUTS.JSON WITH TEMPLATE.

	// templates file --
	templateFile, err := os.ReadFile("templates.json")
	Throw(err)

	var data Template

	err = json.Unmarshal(templateFile, &data) // write json contents to (data)
	Throw(err)

	selectedTemplate := SelectTemplate(data)

	// create file
	file, err := os.Create("inputs.json")
	Throw(err)

	defer file.Close()

	var FinalInputs []Content = ContentInputs(document, data, s.Width, offset, selectedTemplate)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(FinalInputs)
	Throw(err)

	// STEP 2 | WRITING WORD DOCUMENT
	err = document.SaveTo(docName + ".docx")
	Throw(err)

	// STEP 3 | OPEN THE WORD FILE.
	cmd := exec.Command("cmd", "/c", "start", docName+".docx")
	err = cmd.Run()
	Throw(err)
}
