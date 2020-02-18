package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

func parseFileName(args map[string]int) (fileName string, err error) {
	regFileName := regexp.MustCompile(`^\/?([A-z0-9-_+]+\/)*([A-z0-9]+\.(txt))`)

	for key, value := range args {
		matched := regFileName.MatchString(key)
		if matched && value == len(args)-1 {
			fileName = key
			break
		}
	}

	if fileName == "" {
		return "", errors.New("Невалидное имя файла")
	}
	return fileName, nil
}

func readInputFile(args map[string]int) ([]string, error) {
	fileName, errFileName := parseFileName(args)
	if errFileName != nil {
		return []string{}, errFileName
	}

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(fileContent), "\n"), nil
}

func  parseCommands(args []string) (map[string]int, error) {
	commands := make(map[string]int)

	regFlag := regexp.MustCompile(`^-?[A-Za-z]+$|^[0-9]+`)
	regFile := regexp.MustCompile(`^\/?([A-z0-9-_+]+\/)*([A-z0-9]+\.(txt|zip))`)

	for pos, arg := range args {
		matchedFlag := regFlag.MatchString(arg)
		matchedFile := regFile.MatchString(arg)

		if matchedFlag == false && matchedFile == false {
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

	return words, errParse
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
		_, errWrite := file.WriteString(words[pos].original + "\n")
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

func sortInputSlice(words []string, commands map[string]int) (string, error) {
	if len(words) == 0 {
		return "", nil
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
			return "", errNumber
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

	return stringifyResult(wordsState), nil
}

func Execute(args []string) (string, error) {
	commands, errParse := parseCommands(args)
	if errParse != nil {
		return "", errParse
	}

	inputSlice, errRead := readInputFile(commands)
	if errRead != nil {
		return "", errRead
	}

	return sortInputSlice(inputSlice, commands)
}

func main() {
	resultString, err := Execute(os.Args[1:])
	if err != nil {
		log.Fatal(err)
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resultString)
}