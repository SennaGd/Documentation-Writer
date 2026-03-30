const readline = require('node:readline');
const { stdin: input, stdout: output } = require('node:process');

const rl = readline.createInterface({ input, output });
rl.question('What do you think of Node.js? ', (answer1) => {
  console.log(`Feedback: ${answer1}`);

  rl.question('Why? ', (answer2) => {
    console.log(`Reason: ${answer2}`);

    rl.close();
  });
});