package main

import (
	"regexp"
	"testing"
)

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