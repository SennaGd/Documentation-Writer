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

func clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func print(text string) {
	fmt.Print(text)
}

func throw(err any) {
	if err != nil {
		panic(err)
	}
}

func get() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	newInput := input[:len(input)-2] // remove \r\n

	throw(err)
	return newInput
}
func contains(str string, substr string) bool {
	if strings.Contains(str, substr) {
		return true
	} else {
		return false
	}
}

func bigger(var1 string, var2 string) bool {
	if len(var1) >= len(var2) {
		return true
	} else {
		return false
	}
}

func smaller(var1 string, var2 string) bool {
	if len(var1) <= len(var2) {
		return true
	} else {
		return false
	}
}

func convention(input string, convention string, conventionType string, oldConvention string) (string, string) {
	if bigger(input, convention) && contains(input, convention) {
		return input[2:], conventionType
	}

	return input, oldConvention
}

func contentInputs(Document *docx.RootDoc, JSONdata Template, Width int, Offset int, selectedTemplate string) []Content {
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

			generateBox(Width, Offset, "Current Category: "+contents[parsed].Name, func() {
				printTextLine(Width, Offset, 4, "Type 'help' for conventions.")
				printLine(Width, Offset, " ", "│", "│")
				printTextLine(Width, Offset, 2, "Type: "+conType)
				printTextLine(Width, Offset, 2, inputText)

			})
			conType = "" // reset previous convention

			input = get() // user input

			input, conType = convention(input, "h=", "header", conType)
			input, conType = convention(input, "p=", "paragraph", conType)
			input, conType = convention(input, "b=", "bulletpoint", conType)

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
				clear()
				print(
					"│\n" +
						"│ 'enter'-> Return.\n" +
						"│ 'h=' -> Header Element.\n" +
						"│ 'p=' -> Paragraph Element\n" +
						"│ 'b=' -> Bullet Element.\n",
				)
				input = get()
			}
			clear()
		}
		clear()
	}
	fmt.Println("Generating .docx file!")

	return contents
}

func selectTemplate(data Template) string {
	var intInput int
	var err error
	for {
		clear()
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
		print("\nTemplate: ")
		input := get()
		intInput, err = strconv.Atoi(input) // Convert string > integer
		if err == nil {                     // If there's no error select the template
			// check if input is equal to existing template.
			for i, k := range keys {
				if i+1 == intInput {
					clear()
					return k
				}
			}
		}
	}
}

func getDocumentInformation(Width int, Offset int, Document *docx.RootDoc) string {

	generateBox(Width, Offset, "Enter the Filename.", func() {
		printTextLine(Width, Offset, 2, "Filename: None")
		printTextLine(Width, Offset, 2, "Title:    None")

		printLine(Width, Offset, " ", "│", "│")
		printTextLine(Width, Offset, 4, "This will be the filename + '.docx'")
	})

	input := get()

	DocName := input
	clear()
	generateBox(Width, Offset, "Enter the Document Title", func() {
		printTextLine(Width, Offset, 2, "Filename: "+DocName)
		printTextLine(Width, Offset, 2, "Title:    None")
		printLine(Width, Offset, " ", "│", "│")
		printTextLine(Width, Offset, 4, "Title will be inside the file!")
	})

	input = get()
	clear()
	Document.AddHeading(input, 0)

	generateBox(Width, Offset, "Who's it made by?", func() {
		printTextLine(Width, Offset, 2, "Filename: "+DocName)
		printTextLine(Width, Offset, 2, "Title:    "+input)
		printLine(Width, Offset, " ", "│", "│")
		printTextLine(Width, Offset, 4, "This will be inside the file!")
	})
	input = get()
	clear()
	Document.AddParagraph("Gemaakt door: " + input)
	return DocName
}

func printLine(width int, offset int, centerLine string, leftEdge string, rightEdge string) {
	// Rounding to be an equal num
	width /= 2            // Divide w by 2
	width -= 2            // Edges
	width -= (offset * 2) // offset | left & right

	// Offset 1
	for range offset {
		print(" ")
	}

	print(leftEdge)

	for range width * 2 {
		print(centerLine)
	}
	print(rightEdge + "\n")
}

func printIndexTitle(width int, offset int, text string) {
	// Rounding to be an equal num
	if width%2 != 0 {
		width -= 1
	}

	diff := 0

	width /= 2            // Divide w by 2
	width -= 2            // Edges
	width -= (offset * 2) // offset | left & right

	textLength := len(text)

	lineWidth := getTotalLineWidth(width, offset)
	totalWidth := offset + 1 + (width - (textLength / 2)) + len(text) + (width - (textLength / 2)) + 1 + offset
	if lineWidth > totalWidth {
		// 168 - 169 = -1
		diff = lineWidth - totalWidth
	} else {
		// 169 - 168 = 1
		diff = totalWidth - lineWidth
	}
	totalWidth -= diff

	for range offset {
		print(" ")
	}

	if textLength%2 != 0 {
		textLength += 1
	}

	print("│")

	// width =  len(text) / 2
	for range width - (textLength / 2) {
		print(" ")
	}
	print(text)

	for range width - (textLength / 2) + diff {
		print(" ")
	}
	print("│\n")
}

func printTextLine(width int, offset int, textOffset int, text string) {
	// Rounding to be an equal num
	if width%2 != 0 {
		width -= 1
	}

	diff := 0
	width = width/2 - (2 + (offset * 2))

	textLength := len(text)

	lineWidth := getTotalLineWidth(width, offset)
	totalWidth := offset + 1 + (width - (textLength / 2)) + len(text) + (width - (textLength / 2)) + 1 + offset

	if lineWidth > totalWidth {
		// 168 - 169 = -1
		diff = lineWidth - totalWidth
	} else {
		// 169 - 168 = 1
		diff = totalWidth - lineWidth
	}
	totalWidth -= diff
	remaining := totalWidth - offset - offset - textOffset - textLength - 2 - 1

	// i = 10
	for range offset {
		print(" ")
	}

	if textLength%2 != 0 {
		textLength += 1
	}

	print("├") // i = 11

	// width =  len(text) / 2
	for range textOffset {
		print("─")
	}

	print(" ")
	print(text)

	lineWidth = getTotalLineWidth(width, offset)

	for range remaining {
		print(" ")
	}

	print("│\n")
}

func getTotalLineWidth(width int, offset int) int {
	return (offset + 1 + (width * 2) + 1 + offset)
}

func generateBox(width int, offset int, title string, funcs ...func()) {
	printLine(width, offset, "─", "┌", "┐")

	printLine(width, offset, " ", "│", "│")
	printIndexTitle(width, offset, title)
	if len(funcs) >= 1 {
		printLine(width, offset, " ", "│", "│")
	}

	for _, f := range funcs {
		f()
	}

	printLine(width, offset, " ", "│", "│")
	printLine(width, offset, "─", "├", "┘")
	for range offset {
		print(" ")
	}

	print("└───> ")
}

func main() {
	clear()

	// TODO: Width can't be smaller than 100

	var s tsize.Size
	s, err := tsize.GetSize()
	throw(err)

	var offset int = 10

	// STEP 0 | GET TITLE, MADE BY & FILE TITLE
	document, err := godocx.NewDocument()
	throw(err)

	docName := getDocumentInformation(s.Width, offset, document)

	// STEP 1 | FILL INPUTS.JSON WITH TEMPLATE.

	// templates file --
	templateFile, err := os.ReadFile("templates.json")
	throw(err)

	var data Template

	err = json.Unmarshal(templateFile, &data) // write json contents to (data)
	throw(err)

	selectedTemplate := selectTemplate(data)

	// create file
	file, err := os.Create("inputs.json")
	throw(err)

	defer file.Close()

	var FinalInputs []Content = contentInputs(document, data, s.Width, offset, selectedTemplate)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(FinalInputs)
	throw(err)

	// STEP 2 | WRITING WORD DOCUMENT
	err = document.SaveTo(docName + ".docx")
	throw(err)

	// STEP 3 | OPEN THE WORD FILE.
	cmd := exec.Command("cmd", "/c", "start", docName+".docx")
	err = cmd.Run()
	throw(err)
}
