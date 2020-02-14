package main

import (
	"errors"
	"fmt"
	"io/ioutil"
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
		matched, _ := regexp.Match(`^[:alpha:]+\w*`, []byte(arg))

		if matched {
			if (n >= 1 && args[n-1] != "-o") || n <= 0 {
				fileName = arg
				break
			}
		}
	}

	return fileName
}

func readInputFile(args []string) ([]string, error) {
	fileContent, err := ioutil.ReadFile(parseFileName(args))

	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(fileContent), "\n"), nil
}

func  parseCommands(args []string) (map[string]int, error) {
	commands := make(map[string]int)

	for pos, arg := range args {
		matchedFlag, err := regexp.Match(`^-?\w$|[:alpha:]+|[:digit:]+`, []byte(arg))

		//TODO доделать определение пути файла

		if err != nil || !matchedFlag {
			return map[string]int{}, errors.New("Invalid commands")
		} else {
			commands[arg] = pos
		}
	}

	return commands, nil
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

func sortByColumn(words []wordState, commands map[string]int, pos int) ([]wordState, error) {
	var column int
	for key, value := range commands {
		if value == pos+1 {
			var errParse error
			column, errParse = strconv.Atoi(key)

			if errParse != nil {
				return []wordState{}, errParse
			} else {
				break
			}
		}
	}

	for pos := range words {
		columns := strings.Split(words[pos].original, " ")

		if column >= len(columns) {
			words[pos].current = ""
		} else {
			words[pos].current = columns[column]
		}

	}
	return words, nil
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

func sortNumb(words []wordState) ([]wordState, error) {
	var errParse error
	sort.Slice(words, func(i, j int) bool {
		firstNum, errFirst := strconv.ParseFloat(words[i].current, 32)
		secondNum, errSecond := strconv.ParseFloat(words[j].current, 32)

		if errFirst != nil {
			errParse = errFirst
		} else if errSecond != nil {
			errFirst = errSecond
		}
		return  firstNum < secondNum
	})

	if errParse != nil {
		return []wordState{}, errParse
	} else {
		return words, nil
	}
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

func writeInFile(words []wordState, commands map[string]int, pos int) error {
	outPath := findOutPath(commands, pos)

	file, errFile := os.Create(outPath)
	if errFile != nil {
		return errFile
	}

	defer file.Close()

	for pos := range words {
		_, errWrite := file.WriteString(words[pos].original)
		if errWrite != nil {
			return errWrite
		}
	}

	return nil
}

func stringifyResult(words []wordState) string {
	var result string
	for pos := range words {
		result += words[pos].original + "\n"
	}

	return result
}

func sortInputSlice(words []string, commands map[string]int) string {
	if len(words) == 0 {
		return ""
	}

	wordsState := lockOrder(words)

	if pos, commandExist := commands["-k"]; commandExist {
		sortByColumn(wordsState, commands, pos)
	}

	if _, commandExist := commands["-f"]; commandExist {
		wordsState = ignoreCase(wordsState)
	}

	if _, commandExist := commands["-n"]; commandExist {
		var errNumber error
		wordsState, errNumber = sortNumb(wordsState)
		if errNumber != nil {
			return errNumber.Error()
		}
	} else {
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
	}

	return stringifyResult(wordsState)
}

func Execute(args []string) string {
	inputSlice, errRead := readInputFile(args)
	if errRead != nil {
		return errRead.Error()
	}

	commands, errParse := parseCommands(args)
	if errParse != nil {
		return errParse.Error()
	}

	return sortInputSlice(inputSlice, commands)
}

func main() {
	fmt.Println(Execute(os.Args[1:]))
}