package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"unicode/utf8"
)

type previousSymbolType int

const (
	nothing previousSymbolType = 1 + iota
	operation
	value
)

type stackF []float64

func (s *stackF) Push(v float64) {
	*s = append(*s, v)
}

func (s *stackF) Pop() float64 {
	res:=(*s)[len(*s)-1]
	*s=(*s)[:len(*s)-1]
	return res
}

type stackS []string

func (s *stackS) Push(v string) {
	*s = append(*s, v)
}

func (s *stackS) Pop() string {
	res:=(*s)[len(*s)-1]
	*s=(*s)[:len(*s)-1]
	return res
}

func getOpPriority(operator string) int {
	switch operator {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}

	return -1
}

func doOperation(values stackF, operators stackS) (stackF, stackS) {
	second := values.Pop()
	first := values.Pop()

	switch operators.Pop() {
	case "+": values.Push(first+second)
	case "-": values.Push(first-second)
	case "*": values.Push(first*second)
	case "/": values.Push(first/second)
	}

	return values, operators
}

func operationMatched(operation string) bool {
	return operation == "+" || operation == "-" || operation == "*" || operation == "/"
}

func valueMatching(expression string) (float64, string, bool) {
	regularVal := regexp.MustCompile("(^([0-9]*[.])?[0-9]+)|(^)")
	matchedVal := regularVal.FindString(expression)

	parsedVal, err := strconv.ParseFloat(matchedVal, 64)

	return parsedVal, matchedVal, err == nil
}

func calculate(expression string) (float64, error) {
	if expression == "" {
		return 0, nil
	}

	var funcs stackS
	var values stackF
	previousSymbol := nothing

	for len(expression) != 0 {
		if expression[:1] == "(" {
			if previousSymbol == value {
				return 0, errors.New("Пропущен оператор")
			}

			funcs.Push("(")
			expression = expression[1:]
			previousSymbol = nothing
		} else if expression[:1] == ")" {
			for len(funcs) != 0 && funcs[len(funcs)-1] != "(" {
				values, funcs = doOperation(values, funcs)
			}
			if len(funcs) == 0 {
				return 0, errors.New("Скобка не открыта")
			}

			funcs.Pop()
			expression = expression[1:]
			previousSymbol = value
		} else if operationMatched(expression[:1]) {
			if previousSymbol == operation  {
				return 0, errors.New("Пропущен операнд")
			}

			if previousSymbol == nothing {
				values.Push(0)
			}

			for len(funcs) != 0 && getOpPriority(expression[:1]) <= getOpPriority(funcs[len(funcs)-1]) {
				values, funcs = doOperation(values, funcs)
			}

			funcs.Push(expression[:1])
			expression = expression[1:]
			previousSymbol = operation
		} else if parsedVal, matchedVal, valueMatched := valueMatching(expression); valueMatched {
			if previousSymbol == value  {
				return 0, errors.New("Пропущен оператор")
			}

			values.Push(parsedVal)
			expression = expression[utf8.RuneCountInString(matchedVal):]
			previousSymbol = value
		} else {
			return 0, errors.New("Нераспознанный символ")
		}

	}

	for len(funcs) != 0 && len(values) > 1 {
		if funcs[len(funcs)-1] == "(" {
			return 0, errors.New("Скобка не закрыта")
		}
		values, funcs = doOperation(values, funcs)
	}

	return values[0], nil
}

func Execute(expression string) string {
	result, errCalc := calculate(expression)
	if errCalc != nil {
		return errCalc.Error()
	}

	return fmt.Sprint(result)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	fmt.Println(Execute(text[:len(text)-1]))
}