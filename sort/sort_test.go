package main

import (
	"os"
	"testing"
)

var testColumnSortOut = `ffw
ga awg
dad daw
g1 ga`

var testDefaultSortOut = `Apple
BOOK
Book
Go
Hauptbahnhof
January
January
Napkin`

var testReverseSortOut = `Napkin
January
January
Hauptbahnhof
Go
Book
BOOK
Apple`

var testUniqueSortOut = `Apple
BOOK
Book
Go
Hauptbahnhof
January
Napkin`

var testUniqueCaseSortOut = `Apple
BOOK
Go
Hauptbahnhof
January
Napkin`

var testNumberSortOut = `.12
15.6
32
1e2`

var testInvalidNumberSortOut = "Invalid commands"

func TestColumnSort(t *testing.T) {
	result, _ := Execute(Flags {
		"data_columns.txt",
		false,
		false,
		false,
		false,
		"",
		1,

	})

	if testColumnSortOut != result {
		t.Errorf("TestColumnSort faild")
	}
}

func TestDefaultSort(t *testing.T) {
	result, _ := Execute(Flags {
		"data.txt",
		false,
		false,
		false,
		false,
		"",
		-1,

	})

	if testDefaultSortOut != result {
		t.Errorf("TestDefaultSort faild")
	}
}

func TestFileDoesntExist(t *testing.T) {
	_, err := Execute(Flags {
		"nodata.txt",
		false,
		false,
		false,
		false,
		"",
		-1,

	})

	if err == nil {
		t.Errorf("TestInvalidCommand faild")
	}
}

func TestFileDoesntSet(t *testing.T) {
	_, err := Execute(Flags {
		"",
		false,
		false,
		false,
		false,
		"",
		-1,

	})

	if err == nil {
		t.Errorf("TestInvalidCommand faild")
	}
}

func TestReverseSort(t *testing.T) {
	result, _ := Execute(Flags {
		"data.txt",
		false,
		false,
		true,
		false,
		"",
		-1,

	})

	if result != testReverseSortOut {
		t.Errorf("TestReverseSort faild")
	}
}

func TestUniqueSort(t *testing.T) {
	result, _ := Execute(Flags {
		"data.txt",
		false,
		true,
		false,
		false,
		"",
		-1,

	})

	if result != testUniqueSortOut {
		t.Errorf("TestUniqueSort faild")
	}
}

func TestUniqueCaseSort (t *testing.T) {
	result, _ := Execute(Flags {
		"data.txt",
		true,
		true,
		false,
		false,
		"",
		-1,

	})

	if result != testUniqueCaseSortOut {
		t.Errorf("TestUniqueCaseSort faild")
	}
}

func TestNumberSort (t *testing.T) {
	result, _ := Execute(Flags {
		"data_number.txt",
		false,
		false,
		false,
		true,
		"",
		-1,

	})

	if result != testNumberSortOut {
		t.Errorf("TestNumberSort faild")
	}
}

func TestInvalidNumberSort (t *testing.T) {
	_, err := Execute(Flags {
		"data_columns.txt",
		false,
		false,
		false,
		true,
		"",
		-1,

	})

	if err == nil {
		t.Errorf("TestInvalidNumberSort faild")
	}
}

func TestWriteInFile (t *testing.T) {
	Execute(Flags {
		"data.txt",
		false,
		false,
		false,
		false,
		"result.txt",
		-1,

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