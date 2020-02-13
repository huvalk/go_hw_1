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

type pair struct {
	value string
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

func usualSort(words []string) []string{
	sort.Strings(words)

	return words
}

func fixateOrder(words []string) []pair {
	var fixedOrderWords []pair
	for pos := range words {
		fixedOrderWords = append(fixedOrderWords, pair {
			value: words[pos],
			pos: pos,
		})
	}

	return fixedOrderWords
}

func freeOrder(fixedOrderWords []pair) []string {
	var words []string
	for pos := range fixedOrderWords {
		words = append(words, fixedOrderWords[pos].value)
	}

	return words
}

func sortIgnoreCase(words []string) []string{
	fixedOrderWords := fixateOrder(words)

	sort.Slice(fixedOrderWords, func(i, j int) bool {
		firstNum := strings.ToUpper(fixedOrderWords[i].value)
		secondNum := strings.ToUpper(fixedOrderWords[j].value)
		if firstNum == secondNum {
			return fixedOrderWords[i].pos < fixedOrderWords[j].pos
		}

		return  firstNum < secondNum
	})

	return freeOrder(fixedOrderWords)
}

func sortNumb(words []string) []string{
	sort.Slice(words, func(i, j int) bool {
		firstNum, _ := strconv.ParseFloat(words[i], 32)
		secondNum, _ := strconv.ParseFloat(words[j], 32)
		return  firstNum < secondNum
	})

	return words
}

func reverseSlice(words []string) []string{
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}

	return words
}

func uniqueSlice(words []string) []string{
	var uniqueWords []string
	uniqueWords = append(uniqueWords, words[0])
	for i, j := 0, 0; i < len(words); i++ {
		if words[i] != uniqueWords[j] {
			uniqueWords = append(uniqueWords, words[i])
			j++
		}
	}

	return uniqueWords
}

func uniqueSliceIgnoreCase(words []string) []string{
	var uniqueWords []string
	uniqueWords = append(uniqueWords, words[0])
	for i, j := 1, 0; i < len(words); i++ {
		if strings.ToUpper(words[i]) == strings.ToUpper(uniqueWords[j]) {
			continue
		}
		uniqueWords = append(uniqueWords, words[i])
		j++
	}

	return uniqueWords
}

func findOutPath(commands map[string]int, pos int) string {
	for key, value := range commands {
		if value == pos+1 {
			return key
		}
	}

	return ""
}

func writeInFile(words []string, commands map[string]int, pos int) {
	outPath := findOutPath(commands, pos)

	file, err := os.Create(outPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	for pos := range words {
		_, err := file.WriteString(words[pos])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func sortInputSlice(words []string, commands map[string]int) {
	if len(words) == 0 {
		return
	}

	var ignoreCase bool
	if _, commandExist := commands["-n"]; commandExist {
		words = sortNumb(words)
	} else if _, commandExist := commands["-f"]; commandExist {
		words = sortIgnoreCase(words)
		ignoreCase = true
	} else {
		words = usualSort(words)
	}

	if _, commandExist := commands["-r"]; commandExist {
		words = reverseSlice(words)
	}

	if _, commandExist := commands["-u"]; commandExist {
		if ignoreCase {
			words = uniqueSliceIgnoreCase(words)
		} else {
			words = uniqueSlice(words)
		}
	}

	if pos, commandExist := commands["-o"]; commandExist {
		writeInFile(words, commands, pos)
	} else {
		fmt.Println(words)
	}
}

func main() {
	args := os.Args[1:]

	inputSlice := readInputFile(args)
	commands := parseCommands(args)

	sortInputSlice(inputSlice, commands)
}
