import * as readline from 'readline';
import { writeFileSync } from "fs";
import { GenerateWordDocument } from "./docx_convert"

console.Clear(); // Initial Clear cli screen.

const rl = readline.createInterface ({
    input: process.stdin,
    output: process.stdout
});



function getNamingConvention(input:string, Convention: string, type: string): any  {
    Convention = Convention.toLowerCase();

    let search:number =  input.toLowerCase().search(Convention);
    let conventionLength = Convention.length;
    
    if (search !== -1) {
        return [type, input.substring(conventionLength)];
    } else {
        return false;
    }
}


const test: string[] = [
    "test2", 
    "Analyse huidige situatie",
];
const FunctioneelOntwerp: string[] = [
    "Samenvatting", 
    "Analyse huidige situatie",
    "Analyse gewenste situatie",
    "Consequenties",
    "Kosten",
    "Planning"
];

const templates: any[] = [
    {"Functioneel Ontwerp" : FunctioneelOntwerp},
    {"Project Plan" : test}, 
]



var categoryAnswers: any[] = []; 
var jsonGenerated: boolean;
var hasChosenTemplate: boolean;





function appendConventionToArray(NamingConvention: any, categoryObject: any): any {
    if (NamingConvention != false) {        
        let keys = Object.keys(categoryObject);
        category[keys[0]].push(NamingConvention);   // adds "Senna" to the array

        console.log(category[_category]);
        console.Clear();
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

let _category: string = FunctioneelOntwerp[0];
let  category: { [key: string] : string[] } = { [_category] : [] };

function assingCategoryObject(categoryName:string){
    let category: { [key: string] : string[] } = { [categoryName] : [] };
    return category;
}

let hasTemplate: boolean;
let selectedCategory: string[] = [];

function getTemplate() {
    let i = 1;
    templates.forEach(template => {
        let templateString = Object.keys(template)[0]
        console.log("("+ i +"). " + templateString);
        i++;
    });
}
// User Input Loop | Basically creates the categoryAnswers array 
function askNext(i: number, categories: string[]): void { 
    console.Clear();
    getTemplate();

    var b = 1;
    rl.question(`Select a template: `, (input: string) => {
        let selectedCategory;
        let categoryName;
        let IntInput = parseInt(input, 10)
        templates.forEach(template => {
            if (IntInput == b) {
                
                selectedCategory = Object.values(template)[0]
                categoryName = Object.keys(template)[0]
            } 
            b++;
        });

        if (selectedCategory == null) {
            askNext(i, categories);
        } 
        else {
            console.Clear();
        }

    
    
    rl.question(`${selectedCategory[i]} → `, (input2: string) => {
        console.log(categoryName)
        // let intInput: number 
        // if ((intInput =+ input) != null) {
        //     intInput =+ input;
            
        //     console.log(Object.keys(templates[intInput-1])[0]);

        //     if(intInput) 
        //     {
        //         let selectedTemplate = (Object.values(templates)[intInput-1]) 
        //         hasChosenTemplate = true; 
        //     }

        //     }
 
        // Next Category
        if(input2 == "N" || input2 == "n") 
        {     
            if (i + 1 >= categories.length) {
                console.Clear();
                jsonGenerated = true;
                console.log("Writing to output file");
                parseCategoryAnswers(categoryAnswers);

                rl.close();
                return;

            } else {
                console.Clear();
                
                categoryAnswers.push(category);
                console.log(category);
                console.log(categoryAnswers);

                category = assingCategoryObject(categories[i + 1]);    

                askNext(i + 1, categories);
            }
        }

        // Naming Conventions 
        let Header = getNamingConvention(input2, "H=", "Header"); // header : big text
        let Paragraph = getNamingConvention(input2, "P=", "Paragraph"); // paragraph : normal text
        let BulletPoint = getNamingConvention(input2, "B=", "BulletPoint"); // bullet : - text
        let Nonetype = getNamingConvention(input2, "N=", "None"); // Nonetype for names, dates, etc.
        
        let conventionsList: any[] = [Header, Paragraph, BulletPoint, Nonetype];

        conventionsList.forEach(Convention => {
            appendConventionToArray(Convention, category);
        });

        askNext(i, categories);
    });
    })
}


askNext(0, FunctioneelOntwerp)
if(jsonGenerated){
    GenerateWordDocument("");
}

