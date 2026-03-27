package main
import (
	"fmt"
	"os/exec"
	"os"
)

func clear() {
	cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func print(text string) {
	fmt.Printf(text)
}

func input(variable any) {
	fmt.Scan(variable)
}


func main() {
	functioneelOntwerp := [] string {
		"Samenvatting", 
		"Analyse huidige situatie",
		"Analyse gewenste situatie",
		"Consequenties",
		"Kosten",
		"Planning",
	}

	
	
	print("(1). Functioneel Ontwerp\n")
	print("Choose Category: ")

	var selectedCategoryIndex int;
	input(&selectedCategoryIndex)

	
	for _i:=0; _i < len(functioneelOntwerp); _i++ {
		print(functioneelOntwerp[_i] + "\n")
		
	}

var input int
for ok := true; ok; ok = (input != 2) {
    n, err := fmt.Scanln(&input)
    if n < 1 || err != nil {
        fmt.Println("invalid input")
        break
    }

    switch input {
    case 1:
        fmt.Println("hi")
    case 2:
        // Do nothing (we want to exit the loop)
        // In a real program this could be cleanup
    default:
        fmt.Println("def")
    }
}


	print(functioneelOntwerp[0])




}

