package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var mode int
var ROWS, COLS int
var offset_col, offset_row int
var currentCol, currentRow int
var source_file string
var text_buffer = [][]rune{}
var undo_buffer = [][]rune{}
var copy_buffer = []rune{}
var modified bool

func read_file(filename string) {
	file, err := os.Open(filename)

	if err != nil {
		source_file = filename
		text_buffer = append(text_buffer, []rune{})
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		line := scanner.Text()
		text_buffer = append(text_buffer, []rune{})

		for i := 0; i < len(line); i++ {
			text_buffer[lineNumber] = append(text_buffer[lineNumber], rune(line[i]))
		}
		lineNumber++
	}
	if lineNumber == 0 {
		text_buffer = append(text_buffer, []rune{})
	}
}

func write_file(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for row, line := range text_buffer {
		new_line := "\n"
		if row == len(text_buffer)-1 {
			new_line = ""
		}
		write_line := string(line) + new_line
		_, err = writer.WriteString(write_line)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		writer.Flush()
		modified = false
	}
}

func insert_rune(event termbox.Event) {
	insert_rune := make([]rune, len(text_buffer[currentRow])+1)
	copy(insert_rune[:currentCol], text_buffer[currentRow][:currentCol])
	if event.Key == termbox.KeySpace {
		insert_rune[currentCol] = rune(' ')
	} else if event.Key == termbox.KeyTab {
		insert_rune[currentCol] = rune(' ')
	} else {
		insert_rune[currentCol] = rune(event.Ch)
	}
	copy(insert_rune[currentCol+1:], text_buffer[currentRow][currentCol:])
	text_buffer[currentRow] = insert_rune
	currentCol++
}

func delete_rune() {
	if currentCol > 0 {
		currentCol--
		delete_line := make([]rune, len(text_buffer[currentRow])-1)
		copy(delete_line[:currentCol], text_buffer[currentRow][:currentCol])
		copy(delete_line[currentCol:], text_buffer[currentRow][currentCol+1:])
		text_buffer[currentRow] = delete_line
	} else if currentRow > 0 {
		append_line := make([]rune, len(text_buffer[currentRow]))
		copy(append_line, text_buffer[currentRow][currentCol:])
		new_text_buffer := make([][]rune, len(text_buffer)-1)
		copy(new_text_buffer[:currentRow], text_buffer[:currentRow])
		copy(new_text_buffer[currentRow:], text_buffer[currentRow+1:])
		text_buffer = new_text_buffer
		currentRow--
		currentCol = len(text_buffer[currentRow])
		insert_line := make([]rune, len(text_buffer[currentRow])+len(append_line))
		copy(insert_line[:len(text_buffer[currentRow])], text_buffer[currentRow])
		copy(insert_line[len(text_buffer[currentRow]):], append_line)
		text_buffer[currentRow] = insert_line
	}
}

func insert_line() {
	right_line := make([]rune, len(text_buffer[currentRow][currentCol:]))
	copy(right_line, text_buffer[currentRow][currentCol:])
	left_line := make([]rune, len(text_buffer[currentRow][:currentCol]))
	copy(left_line, text_buffer[currentRow][:currentCol])
	text_buffer[currentRow] = left_line
	currentRow++
	currentCol = 0
	new_text_buffer := make([][]rune, len(text_buffer)+1)
	copy(new_text_buffer, text_buffer[:currentRow])
	new_text_buffer[currentRow] = right_line
	copy(new_text_buffer[currentRow+1:], text_buffer[currentRow:])
	text_buffer = new_text_buffer
}

func scroll_text_buffer() {
	if currentRow < offset_row {
		offset_row = currentRow
	}
	if currentCol < offset_col {
		offset_col = currentCol
	}
	if currentRow >= offset_row+ROWS {
		offset_row = currentRow - ROWS + 1
	}
	if currentCol >= offset_col+COLS {
		offset_col = currentCol - COLS + 1
	}
}

func display_text_buffer() {
	var row, col int
	for row = 0; row < ROWS; row++ {
		text_buffer_row := row + offset_row
		for col = 0; col < COLS; col++ {
			textBufferCol := col + offset_col
			if text_buffer_row >= 0 && text_buffer_row < len(text_buffer) && textBufferCol < len(text_buffer[text_buffer_row]) {
				if text_buffer[text_buffer_row][textBufferCol] != '\t' {
					termbox.SetChar(col, row, text_buffer[text_buffer_row][textBufferCol])
				} else {
					termbox.SetCell(col, row, rune(' '), termbox.ColorDefault, termbox.ColorGreen)
				}
			} else if row+offset_row > len(text_buffer)-1 {
				termbox.SetCell(0, row, rune('*'), termbox.ColorBlue, termbox.ColorDefault)
			}
		}
		termbox.SetChar(col, row, rune('\n'))
	}
}

func display_status_bar() {
	var mode_status string
	var file_status string
	var copy_status string
	var undo_status string
	var cursor_status string
	if mode > 0 {
		mode_status = " EDIT: "
	} else {
		mode_status = " VIEW: "
	}
	filename_length := len(source_file)
	if filename_length > 8 {
		filename_length = 8
	}
	file_status = source_file[:filename_length] + " - " + strconv.Itoa(len(text_buffer)) + " lines"
	if modified {
		file_status += " modified"
	} else {
		file_status += " saved"
	}
	cursor_status = " Row " + strconv.Itoa(currentRow+1) + ", Col " + strconv.Itoa(currentCol+1) + " "
	if len(copy_buffer) > 0 {
		copy_status = " [Copy] "
	}
	if len(undo_buffer) > 0 {
		undo_status = " [Undo] "
	}
	used_space := len(mode_status) + len(file_status) + len(cursor_status) + len(copy_status) + len(undo_status)
	spaces := strings.Repeat(" ", COLS-used_space)
	message := mode_status + file_status + copy_status + undo_status + spaces + cursor_status
	print_messages(0, ROWS, termbox.ColorBlack, termbox.ColorWhite, message)
}

func print_messages(col, row int, fg, bg termbox.Attribute, message string) {
	for _, ch := range message {
		termbox.SetCell(col, row, ch, fg, bg)
		col += runewidth.RuneWidth(ch)
	}
}

func get_key() termbox.Event {
	var key_event termbox.Event
	switch event := termbox.PollEvent(); event.Type {
	case termbox.EventKey:
		key_event = event
	case termbox.EventError:
		panic(event.Err)
	}
	return key_event
}

func process_keypress() {
	key_event := get_key()
	if key_event.Key == termbox.KeyEsc {
		mode = 0
	} else if key_event.Ch != 0 {
		if mode == 1 {
			insert_rune(key_event)
			modified = true
		} else {
			switch key_event.Ch {
			case 'q':
				termbox.Close()
				os.Exit(0)
			case 'e':
				mode = 1
			case 'w':
				write_file(source_file)
			}
		}
	} else {
		switch key_event.Key {
		case termbox.KeyEnter:
			if mode == 1 {
				insert_line()
				modified = true
			}
		case termbox.KeyBackspace:
			delete_rune()
			modified = true
		case termbox.KeyBackspace2:
			delete_rune()
			modified = true
		case termbox.KeyTab:
			if mode == 1 {
				for i := 0; i < 4; i++ {
					insert_rune(key_event)
				}
				modified = true
			}
		case termbox.KeySpace:
			if mode == 1 {
				insert_rune(key_event)
				modified = true
			}
		case termbox.KeyHome:
			currentCol = 0
		case termbox.KeyEnd:
			currentCol = len(text_buffer[currentRow])
		case termbox.KeyPgup:
			if currentRow-int(ROWS/4) > 0 {
				currentRow -= int(ROWS / 4)
			}
		case termbox.KeyPgdn:
			if currentRow+int(ROWS/4) < len(text_buffer)-1 {
				currentRow += int(ROWS / 4)
			}
		case termbox.KeyArrowUp:
			if currentRow != 0 {
				currentRow--
			}
		case termbox.KeyArrowDown:
			if currentRow < len(text_buffer)-1 {
				currentRow++
			}
		case termbox.KeyArrowLeft:
			if currentCol != 0 {
				currentCol--
			} else if currentRow > 0 {
				currentRow--
				currentCol = len(text_buffer[currentRow])
			}
		case termbox.KeyArrowRight:
			if currentCol < len(text_buffer[currentRow]) {
				currentCol++
			} else if currentRow < len(text_buffer)-1 {
				currentRow++
				currentCol = 0
			}
		}
		if currentCol > len(text_buffer[currentRow]) {
			currentCol = len(text_buffer[currentRow])
		}
	}
}

func run_editor() {
	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		source_file = os.Args[1]
		read_file((source_file))
	} else {
		source_file = "out.txt"
		text_buffer = append(text_buffer, []rune{})
	}

	for {
		COLS, ROWS = termbox.Size()
		ROWS--
		if COLS < 80 {
			COLS = 80
		}
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		scroll_text_buffer()
		display_text_buffer()
		display_status_bar()
		termbox.SetCursor(currentCol-offset_col, currentRow-offset_row)
		termbox.Flush()
		process_keypress()
	}
}

func main() {
	run_editor()
}
