package main

import (
	"os"
	"regexp"
	"testing"
)

var testColumnSortOut = `ffw
ga awg
dad daw
g1 ga
`

var testDefaultSortOut = `Apple
BOOK
Book
Go
Hauptbahnhof
January
January
Napkin
`

var testReverseSortOut = `Napkin
January
January
Hauptbahnhof
Go
Book
BOOK
Apple
`

var testUniqueSortOut = `Apple
BOOK
Book
Go
Hauptbahnhof
January
Napkin
`

var testUniqueCaseSortOut = `Apple
BOOK
Go
Hauptbahnhof
January
Napkin
`
var testNumberSortOut = `.12
15.6
32
1e2
`

var testInvalidNumberSortOut = "Invalid commands"

func TestColumnSort(t *testing.T) {
	result := Execute([]string{
		"-k",
		"1",
		"data_columns.txt",
	})

	if testColumnSortOut != result {
		t.Errorf("TestColumnSort faild")
	}
}

func TestDefaultSort(t *testing.T) {
	result := Execute([]string{
		"data.txt",
	})

	if testDefaultSortOut != result {
		t.Errorf("TestDefaultSort faild")
	}
}

func TestInvalidCommand(t *testing.T) {
	result := Execute([]string{
		"?:faa",
		"data.txt",
	})

	if "Invalid commands" != result {
		t.Errorf("TestInvalidCommand faild")
	}
}

func TestReverseSort(t *testing.T) {
	result := Execute([]string{
		"-r",
		"data.txt",
	})

	if result != testReverseSortOut {
		t.Errorf("TestReverseSort faild")
	}
}

func TestUniqueSort(t *testing.T) {
	result := Execute([]string{
		"-u",
		"data.txt",
	})

	if result != testUniqueSortOut {
		t.Errorf("TestUniqueSort faild")
	}
}

func TestUniqueCaseSort (t *testing.T) {
	result := Execute([]string{
		"-u",
		"-f",
		"data.txt",
	})

	if result != testUniqueCaseSortOut {
		t.Errorf("TestUniqueCaseSort faild")
	}
}

func TestNumberSort (t *testing.T) {
	result := Execute([]string{
		"-n",
		"data_number.txt",
	})

	if result != testNumberSortOut {
		t.Errorf("TestNumberSort faild")
	}
}

func TestInvalidNumberSort (t *testing.T) {
	result := Execute([]string{
		"-n",
		"data_columns.txt",
	})

	matched, _ := regexp.Match(`invalid syntax$`, []byte(result))

	if !matched {
		t.Errorf("TestInvalidNumberSort faild")
	}
}

func TestWriteInFile (t *testing.T) {
	Execute([]string{
		"-o",
		"result.txt",
		"data.txt",
	})

	exists := true
	info, err := os.Stat("result.txt")
	if os.IsNotExist(err) {
		exists = false
	}
	exists = !info.IsDir()

	if !exists {
		t.Errorf("TestWriteInFile faild")
	}
}