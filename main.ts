import * as readline from 'readline';

console.clear() // Initial clear cli screen.

const rl = readline.createInterface({
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
    "Made by",
    "Doelgroep",
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
var firstPushed: boolean = false;


function pushFirstCategory(answersArray: string[], categoryArray: string[]): void {
    answersArray.push( categoryArray[0] );
}

function appendConventionToArray(NamingConvention: any): any {
            if (NamingConvention != false) 
            {   
                categoryAnswers.push(NamingConvention);
                console.clear();
                console.log("categoryAnswers: ", categoryAnswers);
            } 
}


// User Input Loop | Basically creates the categoryAnswers array 
function askNext(i: number, categories: string[]): void {    
    if (!firstPushed){
        pushFirstCategory(categoryAnswers, categories);
        firstPushed = true;
    } 

    rl.question(`${categories[i]}: `, (input: string) => {
        // Next Category
        if(input == "N" || input == "n") 
        {     
            if (i + 1 >= categories.length) {
                console.clear()
                console.log("bye, bye!")
                rl.close();
                return;
            } else {
                console.clear();
                categoryAnswers.push(categories[i+1]);
                askNext(i + 1, categories);
            }
        }

        // Naming Conventions 
        let Header = getNamingConvention(input, "H=", "Header"); // header : big text
        let Paragraph = getNamingConvention(input, "P=", "Paragraph"); // paragraph : normal text
        let BulletPoint = getNamingConvention(input, "B=", "BulletPoint"); // bullet : - text
        let Nonetype = getNamingConvention(input, "../", "None"); // Nonetype for names, dates, etc.
        
        let conventionsList: any[] = [Header, Paragraph, BulletPoint, Nonetype] 

        conventionsList.forEach(Convention => {
            appendConventionToArray(Convention);
        });

        askNext(i, categories);
    });
}

askNext(0, categories);




