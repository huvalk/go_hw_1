package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type wordState struct {
	original string
	current string
	pos int
}

func parseFileName(args []string) string {
	var fileName string
	for n, arg := range args {
		matched, _ := regexp.Match(`^\w+`, []byte(arg))

		if matched {
			if (n >= 1 && args[n-1] != "-o") || n <= 0 {
				fileName = arg
				break
			}
		}
	}

	return fileName
}

func readInputFile(args []string) []string {
	file, err := os.Open(parseFileName(args))

	defer file.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	text := make([]byte, 64)
	for{
		_, err := file.Read(text)
		if err == io.EOF {
			break
		}
	}

	return strings.Split(string(text), "\n")
}

func  parseCommands(args []string) map[string]int {
	commands := make(map[string]int)

	for pos, arg := range args {
		matched, _ := regexp.Match(`^-?\w+`, []byte(arg))

		if matched {
			commands[arg] = pos
		} else {
			fmt.Println("Invalid command")
			os.Exit(1)
		}
	}

	return commands
}

func lockOrder(words []string) []wordState {
	var fixedOrderWords []wordState
	for pos := range words {
		fixedOrderWords = append(fixedOrderWords, wordState{
			original: words[pos],
			current: words[pos],
			pos: pos,
		})
	}

	return fixedOrderWords
}

func usualSort(words []wordState) []wordState{
	sort.Slice(words, func(i, j int) bool {
		if words[i].current == words[j].current {
			return words[i].pos < words[j].pos
		} else {
			return  words[i].current < words[j].current
		}
	})

	return words
}

func sortNumb(words []wordState) []wordState{
	sort.Slice(words, func(i, j int) bool {
		firstNum, _ := strconv.ParseFloat(words[i].current, 32)
		secondNum, _ := strconv.ParseFloat(words[j].current, 32)
		return  firstNum < secondNum
	})

	return words
}

func ignoreCase(words []wordState) []wordState{
	for pos := range words {
		words[pos].current = strings.ToUpper(words[pos].current)
	}

	return words
}



func uniqueSlice(words []wordState) []wordState{
	var uniqueWords []wordState
	uniqueWords = append(uniqueWords, words[0])
	for i, j := 0, 0; i < len(words); i++ {
		if words[i].current != uniqueWords[j].current {
			uniqueWords = append(uniqueWords, words[i])
			j++
		}
	}

	return uniqueWords
}

func reverseSlice(words []wordState) []wordState{
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}

	return words
}

func findOutPath(commands map[string]int, pos int) string {
	for key, value := range commands {
		if value == pos+1 {
			return key
		}
	}

	return ""
}

func writeInFile(words []wordState, commands map[string]int, pos int) {
	outPath := findOutPath(commands, pos)

	file, err := os.Create(outPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	for pos := range words {
		_, err := file.WriteString(words[pos].original)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func printResult(words []wordState) {
	for pos := range words {
		fmt.Println(words[pos].original)
	}
}

func sortInputSlice(words []string, commands map[string]int) {
	if len(words) == 0 {
		return
	}

	wordsState := lockOrder(words)

	if _, commandExist := commands["-f"]; commandExist {
		wordsState = ignoreCase(wordsState)
	}

	needUsualSort := true
	if _, commandExist := commands["-n"]; commandExist {
		wordsState = sortNumb(wordsState)
		needUsualSort = false
	}
	if needUsualSort {
		wordsState = usualSort(wordsState)
	}

	if _, commandExist := commands["-u"]; commandExist {
		wordsState = uniqueSlice(wordsState)
	}

	if _, commandExist := commands["-r"]; commandExist {
		wordsState = reverseSlice(wordsState)
	}

	if pos, commandExist := commands["-o"]; commandExist {
		writeInFile(wordsState, commands, pos)
	} else {
		printResult(wordsState)
	}
}

func main() {
	args := os.Args[1:]

	inputSlice := readInputFile(args)
	commands := parseCommands(args)

	sortInputSlice(inputSlice, commands)
}
