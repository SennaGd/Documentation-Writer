import * as readline from 'readline';
import { writeFileSync } from "fs";
import { GenerateWordDocument } from "./docx_convert"

console.clear(); // Initial clear cli screen.

const rl = readline.createInterface ({
    input: process.stdin,
    output: process.stdout
});



function getNamingConvention(input:string, convention: string, type: string): any  {
    convention = convention.toLowerCase();

    let search:number =  input.toLowerCase().search(convention);
    let conventionLength = convention.length;
    
    if (search !== -1) {
        return [type, input.substring(conventionLength)];
    } else {
        return false;
    }
}


const categories: string[] = [
    "Samenvatting", 
    "Analyse huidige situatie",
    "Analyse gewenste situatie",
    "Consequenties",
    "Kosten",
    "Planning"
];

const testcategories: string[] = [
    "Readability", "Context"
];

var categoryAnswers: any[] = []; 

function appendConventionToArray(NamingConvention: any, categoryObject: any): any {
    if (NamingConvention != false) {        
        let keys = Object.keys(categoryObject);
        category[keys[0]].push(NamingConvention);   // adds "Senna" to the array

        console.log(category[_category]);
        console.clear();
        console.log("Appended: ", NamingConvention);
        // console.log("categoryAnswers: ", categoryAnswers);
    } 
}

function parseCategoryAnswers(answersArray: string[]){
    console.log(answersArray)
    let json = JSON.stringify(answersArray);
    writeFileSync("file.json", json, {
        flag: "w"
    });
}

let _category: string = categories[0];
let  category: { [key: string] : string[] } = { [_category] : [] };

function assingCategoryObject(categoryName:string){
    let category: { [key: string] : string[] } = { [categoryName] : [] };
    return category;
}

// User Input Loop | Basically creates the categoryAnswers array 
function askNext(i: number, categories: string[]): void {    
    rl.question(`${categories[i]} → `, (input: string) => {
        // Next Category
        if(input == "N" || input == "n") 
        {     
            if (i + 1 >= categories.length) {
                console.clear();
                console.log("Writing to output file");
                parseCategoryAnswers(categoryAnswers);

                rl.close();
                return;
            } else {
                console.clear();
                
                categoryAnswers.push(category);
                console.log(category);
                console.log(categoryAnswers);

                category = assingCategoryObject(categories[i + 1]);    

                askNext(i + 1, categories);
            }
        }

        // Naming Conventions 
        let Header = getNamingConvention(input, "H=", "Header"); // header : big text
        let Paragraph = getNamingConvention(input, "P=", "Paragraph"); // paragraph : normal text
        let BulletPoint = getNamingConvention(input, "B=", "BulletPoint"); // bullet : - text
        let Nonetype = getNamingConvention(input, "N=", "None"); // Nonetype for names, dates, etc.
        
        let conventionsList: any[] = [Header, Paragraph, BulletPoint, Nonetype];

        conventionsList.forEach(Convention => {
            appendConventionToArray(Convention, category);
        });

        askNext(i, categories);
    });
}

askNext(0, categories);

GenerateWordDocument("")
