import * as fs from "fs";
import { Document, Packer, Paragraph, TextRun } from "docx";
import path from 'path';

type DataItem = {
  [key: string]: string[][];
};

// const ConventionTypes: any[] = [
//     {"None" : ""},
// ] 

const data: DataItem[] = require("./file.json"); 

let index: number = 0;
data.forEach(category => {
    const firstItem = data[index];
    if (!firstItem) 
    { return; }

    const findKey = Object.keys(firstItem)[0];
    if (!findKey)
    { return; }

    const findValues = firstItem[findKey];
    if(!findValues?.[0])
    { return; }

    // Traverses all contents from a (category | key)
    for(let i: number = 0; i < (findValues.length); i++) 
    { console.log(findValues[i]); }

    index++;
});
// const doc = new Document({
//     sections: [
//         {
//             children: [
//                 new Paragraph({
                    
//                     children: [
//                         new TextRun("Hello World"),
//                         new TextRun({
//                             text: " - Bold text",
//                             bold: true,
//                             size: 50
//                         }),
//                         new TextRun({
//                             text: "header",
//                             size: 20,
//                         })
//                     ],
//                 }),
//             ],
//         },
//     ],
// });

// Packer.toBuffer(doc).then((buffer) => {
//     fs.writeFileSync("My Document.docx", buffer);
// });