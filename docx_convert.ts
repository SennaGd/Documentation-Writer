import * as fs from "fs";
import { Document, Packer, Paragraph, TextRun, AlignmentType  } from "docx";

type DataItem = {
  [key: string]: string[][];
};

const data: DataItem[] = require("./file.json"); 

// ---- Code ---- //

export function GenerateWordDocument(title: string) {
    let contentsList: Paragraph[] = []


    data.forEach(category => {
        var indentation = 0;
        let findKey = Object.keys(category)[0];

        if (!findKey) { 
            return; 
        } else { 
            var FirstPlaced: boolean = false;
            var headerSpacing = 0;

            contentsList.push(new Paragraph(
                { 
                    spacing : {before: 400}, 
                    children : [ new TextRun({
                        text: findKey, 
                        size: 40, 
                        bold: true, 
                        characterSpacing: 20 
                    })] 
                }));
        }


        const findValues = category[findKey];
        if(!findValues?.[0])
        { return; }

        
        //* Traverses all contents from a (category | key)
        //* Data inside category
        for(let i: number = 0; i < (findValues.length); i++) { 
            if (findValues[i][0] == "Header") 
            {   
                var indentation = 360;
                if (FirstPlaced) {
                    headerSpacing = 100;
                }
                
                var content = new Paragraph(
                { 
                    indent:{ left: indentation },
                    spacing : {before: headerSpacing},    
                    children: [ 
                        new TextRun({font: "Calibri", text: findValues[i][1], size: 30, characterSpacing: 10}) 
                    ]
                })

                var indentation = 360;
                FirstPlaced = true; 
            }
            
            if (findValues[i][0] == "Paragraph") 
            {           
                var indentation = 180;
                var content = new Paragraph(
                { 
                    indent:{ left: indentation },
                    spacing : {before: 50},    
                    children: [ 
                        new TextRun({font: "Calibri", text: findValues[i][1], size: 24}) 
                    ]
                })
                
                
            }; 

            if (findValues[i][0] == "BulletPoint") {
                var content = new Paragraph(
                { 
                    spacing : {before: 0},    
                    children: [ 
                        new TextRun({font: "Calibri", text: findValues[i][1], size: 24, italics: true}) 
                    ]
                })
            }

            contentsList.push(content); // append content on content.
            console.log(findValues[i]); 
        }

    });

    // Document!

    const doc = new Document({
        sections: [ { children: contentsList } ],
           
        styles: {
            default: {
                document: {
                run: {
                    font: "Calibri",
                    size: 22
                }
            }
        }
    }
})


    Packer.toBuffer(doc).then((buffer) => {
        fs.writeFileSync(title + ".docx", buffer);
    });
}


GenerateWordDocument("")










