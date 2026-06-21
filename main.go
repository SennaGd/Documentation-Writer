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

	// "github.com/acarl005/stripansi"
	"github.com/fatih/color"
	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"

	graphics "DocumentationWriter/graphics"

	tsize "github.com/kopoli/go-terminal-size"
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

	newInput := input[:len(input)-1] // remove \r\n

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

func HelpPage(Width int, Offset int) {
	graphics.GenerateBox(Width, Offset, color.BlueString("Help Page"), func() {
		graphics.PrintIndexTitle(Width, Offset, "Where you felt lost, you found it.")
		graphics.PrintLine(Width, Offset, " ", "│", "│")
		graphics.PrintLine(Width, Offset, "─", "├", "┤")

		graphics.PrintIndexTitle(Width, Offset, color.YellowString("Naming Conventions."))
		graphics.PrintLine(Width, Offset, " ", "│", "│")
		graphics.PrintTextLine(Width, Offset, 4, "Write Convention ("+color.RedString("'P='")+") : followed by text")
		graphics.PrintLine(Width, Offset, " ", "│", "│")

		graphics.PrintTextLine(Width, Offset, 2, color.GreenString("'B='")+" This is a "+color.GreenString("Bullet")+" Element")
		graphics.PrintTextLine(Width, Offset, 2, color.GreenString("'H='")+" This is a "+color.GreenString("Header")+" Element")
		graphics.PrintTextLine(Width, Offset, 2, color.GreenString("'P='")+" This is a "+color.GreenString("Paragraph")+" Element")
		graphics.PrintLine(Width, Offset, "─", "├", "┤")
		graphics.PrintIndexTitle(Width, Offset, color.YellowString("Go To Next Category"))
		graphics.PrintLine(Width, Offset, " ", "│", "│")
		graphics.PrintTextLine(Width, Offset, 2, color.GreenString("'N'")+" To traverse to the next "+color.GreenString("Category"))
		graphics.PrintLine(Width, Offset, " ", "│", "│")
	})
}

func ContentInputs(Document *docx.RootDoc, JSONdata Template, Width int, Offset int, selectedTemplate string) []Content {
	/*
		- inputsList = list of all inputs before "next category"
		- categoryList = the name of the category + all the inputs of the category (inputsList)
	*/

	var input string
	var conType string = "Nothing Written yet."

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
				inputText = ("Contents: " + color.RedString("lots of text Added!"))
			} else if len(input) <= 1 {
				inputText = ("Contents: " + color.YellowString("Nothing added yet."))
			} else {
				inputText = ("Contents: " + color.GreenString(input))
			}

			graphics.GenerateBox(Width, Offset, (color.BlueString("Current Category: ") + color.CyanString(contents[parsed].Name)), func() {
				graphics.PrintLine(Width, Offset, " ", "│", "│")
				graphics.PrintTextLine(Width, Offset, 2, "Type: "+color.GreenString(conType))
				graphics.PrintTextLine(Width, Offset, 2, inputText)
				graphics.PrintLine(Width, Offset, " ", "│", "│")
				graphics.PrintIndexTitle(Width, Offset, "Type "+color.YellowString("'help'")+" for conventions.")
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

			if input == "help" || input == "Help" || input == "HELP" {
				Clear()
				HelpPage(Width, Offset)
				input = Get()
			}
			Clear()
		}
		Clear()
	}
	fmt.Println("Generating .docx file!")

	return contents
}

func SelectTemplate(data Template, Width int, Offset int) string {
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

		graphics.GenerateBox(Width, Offset, color.BlueString("Select a category."), func() {
			for i, k := range keys {
				text := (color.WhiteString("("+strconv.Itoa(i+1)+"). ") + color.GreenString(k))
				graphics.PrintTextLine(Width, Offset, 2, text)

			}
			graphics.PrintLine(Width, Offset, " ", "│", "│")
			graphics.PrintIndexTitle(Width, Offset, "Type "+color.YellowString("'help'")+" for conventions.")

		})

		// Input handling
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

		if input == "help" || input == "Help" || input == "HELP" {
			Clear()
			HelpPage(Width, Offset)
			Get()
		}
	}
}

func GetDocumentInformation(Width int, Offset int, Document *docx.RootDoc) string {
	fileName := color.YellowString("None")
	title := color.CyanString("None")

	graphics.GenerateBox(Width, Offset, color.BlueString("Enter the Filename."), func() {
		graphics.PrintTextLine(Width, Offset, 2, "Filename: "+fileName)
		graphics.PrintTextLine(Width, Offset, 2, "Title:    "+title)

		graphics.PrintLine(Width, Offset, " ", "│", "│")
		graphics.PrintIndexTitle(Width, Offset, "This will be the filename + '.docx'")
	})

	input := Get()

	DocName := input
	fileName = color.GreenString(input)
	Clear()

	title = color.YellowString("None")
	graphics.GenerateBox(Width, Offset, color.BlueString("Enter the Document Title"), func() {
		graphics.PrintTextLine(Width, Offset, 2, "Filename: "+fileName+color.GreenString(".docx"))
		graphics.PrintTextLine(Width, Offset, 2, "Title:    "+title)
		graphics.PrintLine(Width, Offset, " ", "│", "│")
		graphics.PrintIndexTitle(Width, Offset, "Title will be inside the file!")
	})
	input = Get()
	title = color.GreenString(input)
	Clear()

	Document.AddHeading(input, 0)
	graphics.GenerateBox(Width, Offset, color.BlueString("Who's it made by?"), func() {
		graphics.PrintTextLine(Width, Offset, 2, "Filename: "+fileName+color.GreenString(".docx"))
		graphics.PrintTextLine(Width, Offset, 2, "Title:    "+title)
		graphics.PrintLine(Width, Offset, " ", "│", "│")
		graphics.PrintIndexTitle(Width, Offset, "This will be inside the file!")
	})
	input = Get()
	Clear()
	Document.AddParagraph("Gemaakt door: " + input)
	return DocName
}

func GetSize(Offset int) tsize.Size {
	for {
		var s tsize.Size
		var err error
		s, err = tsize.GetSize()
		Throw(err)

		if s.Width < 100 {
			Clear()
			graphics.GenerateBox(s.Width, Offset, color.BlueString("Increase Window Size..."), func() {
				stringSize := strconv.Itoa(s.Width)
				size := color.YellowString(stringSize)

				graphics.PrintTextLine(s.Width, Offset, 2, "Size: "+size)
				graphics.PrintTextLine(s.Width, Offset, 2, "Minimum: "+color.GreenString("100"))
				graphics.PrintLine(s.Width, Offset, " ", "│", "│")
				graphics.PrintIndexTitle(s.Width, Offset, "Press Enter To Refresh")
			})
			Get()
		} else {
			Clear()
			return s
		}
	}
}

func main() {
	Clear()

	// TODO: Width can't be Smaller than 100

	var s tsize.Size
	var offset int = 10
	s = GetSize(offset)

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
	
	fmt.Println(&data)
	selectedTemplate := SelectTemplate(data, s.Width, offset)
	// create file
	file, err := os.Create("inputs.json")
	Throw(err)

	fmt.Println(&selectedTemplate)
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
