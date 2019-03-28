# typeme

Simple typing game where you need to type words as they appear on the screen before
leaving it

[![asciicast](https://asciinema.org/a/RwX3DQEoZOzsy634Fld8sgPkn.svg)](https://asciinema.org/a/RwX3DQEoZOzsy634Fld8sgPkn)

## About

This was made over the course of two days so the code is a little bit ugly (but hey! it works).
I mainly wanted any idea to use the (https://github.com/gdamore/tcell)[https://github.com/gdamore/tcell] library

Basically words will appear between the first 10 rows of your terminal and they will slowly start to go towards the end of the screen.
You need to type each word before they arrive to the end (this ia a game over), also the speed of the rows increase along the game length

## Words

All the words come from `words.json`, basically a big `JSON` file of english words.
You can change this file with anything (that is an array) and it should work fine (if you want to use another language for example)

# License

Do whatever you want. This 'game' is built mainly using (https://github.com/gdamore/tcell)[https://github.com/gdamore/tcell]

