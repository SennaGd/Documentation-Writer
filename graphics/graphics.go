package graphics

import ("fmt"
 "os/exec"
	"os")

func Print(text string) {
	fmt.Print(text)
}

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func PrintLine(width int, offset int, centerLine string, leftEdge string, rightEdge string) {
	// Rounding to be an equal num
	width /= 2            // Divide w by 2
	width -= 2            // Edges
	width -= (offset * 2) // offset | left & right

	// Offset 1
	for range offset {
		Print(" ")
	}

	Print(leftEdge)

	for range width * 2 {
		Print(centerLine)
	}
	Print(rightEdge + "\n")
}

func PrintIndexTitle(width int, offset int, text string) {
	// Rounding to be an equal num
	if width%2 != 0 {
		width -= 1
	}

	diff := 0

	width /= 2            // Divide w by 2
	width -= 2            // Edges
	width -= (offset * 2) // offset | left & right

	textLength := len(text)

	lineWidth := GetTotalLineWidth(width, offset)
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
		Print(" ")
	}

	if textLength%2 != 0 {
		textLength += 1
	}

	Print("│")

	// width =  len(text) / 2
	for range width - (textLength / 2) {
		Print(" ")
	}
	Print(text)

	for range width - (textLength / 2) + diff {
		Print(" ")
	}
	Print("│\n")
}

func PrintTextLine(width int, offset int, textOffset int, text string) {
	// Rounding to be an equal num
	if width%2 != 0 {
		width -= 1
	}

	diff := 0
	width = width/2 - (2 + (offset * 2))

	textLength := len(text)

	lineWidth := GetTotalLineWidth(width, offset)
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
		Print(" ")
	}

	if textLength%2 != 0 {
		textLength += 1
	}

	Print("├") // i = 11

	// width =  len(text) / 2
	for range textOffset {
		Print("─")
	}

	Print(" ")
	Print(text)

	lineWidth = GetTotalLineWidth(width, offset)

	for range remaining {
		Print(" ")
	}

	Print("│\n")
}

func GetTotalLineWidth(width int, offset int) int {
	return (offset + 1 + (width * 2) + 1 + offset)
}

func GenerateBox(width int, offset int, title string, funcs ...func()) {
	PrintLine(width, offset, "─", "┌", "┐")

	PrintLine(width, offset, " ", "│", "│")
	PrintIndexTitle(width, offset, title)
	if len(funcs) >= 1 {
		PrintLine(width, offset, " ", "│", "│")
	}

	for _, f := range funcs {
		f()
	}

	PrintLine(width, offset, " ", "│", "│")
	PrintLine(width, offset, "─", "├", "┘")
	for range offset {
		Print(" ")
	}

	Print("└───> ")
}
