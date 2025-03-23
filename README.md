🚀 SigmaText - My Own Vim-Like Text Editor in Go! 🚀  

I am in process of learning Go and what better way to learn other than working on a project? So I decided to build my own Vim like terminal based text editor from scratch! 🖥️  

### 🔹 Features So Far:  
✅ Create, edit, and view files  
✅ Edit & View modes for better control  
✅ Copy/Paste functionality  
✅ Undo/Redo for seamless editing  
✅ Syntax highlighting to enhance readability  

### 💡 Challenges Faced:  
🔸 Maintaining the cursor state dynamically while navigating  
🔸 Implementing an efficient undo/redo system for endless iterations  

This project has been a deep dive into **state management, data structures, and terminal-based UI design**. There's still a lot to improve and optimize, but it's been an exciting journey so far!

### Installation
🔹 Clone the repo on your machine:
```
git clone https://github.com/Amandeep2230/SigmaText.git
```
🔹 Set directory to the project and run the executable file

🔹 To open an existing file: 
```
./sigmatext <existing_file_name>
```
🔹 To create a new file:
```
./sigmatext <new_file_name>
```

### Commands
🔹 'e' => Edit Mode
🔹 'Esc' => Toggle back to view mode
🔹 'q' => close the editor
🔹 'w' => write/save changes to a file
🔹 'c' => copy
🔹 'v' => paste
🔹 'd' => cut
🔹 's' => undo (will maintain state of changes being made)
🔹 'l' => redo (will rollback to previous pre-change state)
🔹 'h' => toggle text highlight


This is a code along project along with some modifications from my end, special thanks to: https://youtube.com/playlist?list=PLLfIBXQeu3aa0NI4RT5OuRQsLo6gtLwGN&si=X5b6wCgdzx9NJlAE

-Cheers!
