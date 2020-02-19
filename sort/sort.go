package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type wordState struct {
	original string
	current string
}

type Flags struct {
	input string
	ignoreFont bool
	unique bool
	reverse bool
	numbers bool
	output string
	column int
}

func readInputFile(args string) ([]string, error) {
	fileName := args
	if fileName == "" {
		return []string{}, errors.New("Входной файл не задан")
	}

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(fileContent), "\n"), nil
}

func  parseCommands() (parsedFlags Flags) {
	parsedFlags = Flags{
		"",
		false,
		false,
		false,
		false,
		"",
		-1,
	}
	flag.BoolVar(&parsedFlags.ignoreFont, "f", false, "Ignore font")
	flag.BoolVar(&parsedFlags.unique, "u", false, "Unique")
	flag.BoolVar(&parsedFlags.reverse, "r", false, "Reverse sort")
	flag.BoolVar(&parsedFlags.numbers, "n", false, "Sort numbers")
	flag.StringVar(&parsedFlags.output, "o", "", "Output file")
	flag.IntVar(&parsedFlags.column, "k", -1, "Output file")
	flag.Parse()
	fmt.Println(flag.Args()[0])
	parsedFlags.input = flag.Args()[0]

	return parsedFlags
}

func lockOrder(words []string) []wordState {
	fixedOrderWords := make([]wordState, 0, len(words))
	for pos := range words {
		fixedOrderWords = append(fixedOrderWords, wordState{
			original: words[pos],
			current: words[pos],
		})
	}

	return fixedOrderWords
}

func sortByColumn(words []wordState, column int) ([]wordState, error) {
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

func usualSort(words []wordState, reverse bool) []wordState{
	sort.SliceStable(words, func(i, j int) bool {
		if reverse {
			return words[i].current > words[j].current
		}
		return  words[i].current < words[j].current
	})

	return words
}

func sortNumb(words []wordState, reverse bool) ([]wordState, error) {
	var errParse error
	sort.SliceStable(words, func(i, j int) bool {
		firstNum, errFirst := strconv.ParseFloat(words[i].current, 32)
		secondNum, errSecond := strconv.ParseFloat(words[j].current, 32)

		if errFirst != nil {
			errParse = errFirst
		} else if errSecond != nil {
			errFirst = errSecond
		}

		if reverse {
			return  firstNum > secondNum
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

func uniqueSlice(words []wordState) (uniqueWords []wordState) {
	keys := make(map[string]bool)
	for _, entry := range words {
		if _, value := keys[entry.current]; !value {
			keys[entry.current] = true
			uniqueWords = append(uniqueWords, entry)
		}
	}

	return uniqueWords
}

func writeInFile(words string, outPath string) error {
	file, errFile := os.Create(outPath)
	if errFile != nil {
		return errFile
	}

	defer file.Close()

	_, errWrite := file.WriteString(words)
	if errWrite != nil {
		return errWrite
	}

	return nil
}

func stringifyResult(words []wordState) (result string) {
	for pos := range words {
		result += words[pos].original
		if pos != len(words)-1 {
			result += "\n"
		}
	}

	return result
}

func sortInputSlice(words []string, commands Flags) (string, error) {
	if len(words) == 0 {
		return "", nil
	}

	wordsState := lockOrder(words)

	if commands.column > -1 {
		sortByColumn(wordsState, commands.column)
	}

	if commands.ignoreFont {
		wordsState = ignoreCase(wordsState)
	}

	if commands.unique {
		wordsState = uniqueSlice(wordsState)
	}

	if commands.numbers {
		var errNumber error
		wordsState, errNumber = sortNumb(wordsState, commands.reverse)
		if errNumber != nil {
			return "", errNumber
		}
	} else {
		wordsState = usualSort(wordsState, commands.reverse)
	}

	return stringifyResult(wordsState), nil
}

func Execute(parsedFlags Flags) (string, error) {
	inputSlice, errRead := readInputFile(parsedFlags.input)
	if errRead != nil {
		return "", errRead
	}

	resultString, err := sortInputSlice(inputSlice, parsedFlags)
	if err != nil {
		return "", err
	}

	if parsedFlags.output != "" {
		writeInFile(resultString, parsedFlags.output)
	}
	return sortInputSlice(inputSlice, parsedFlags)
}

func main() {
	parsedFlags := parseCommands()
	resultString, err := Execute(parsedFlags)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resultString)
}