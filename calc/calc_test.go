package main

import (
	"testing"
)

var TestUsualExpOut = "60"

var TestBracketsExpOut = "12"

var TestUnaryExpOut = "0"

var TestNotFinishedExpOut = "-3"

var TestOperatorsInARowExpOut = "Пропущен операнд"

var TestOperandsInARowExpOut = "Пропущен оператор"

var TestUnknownSymbolExpOut = "Нераспознанный символ"

var TestMissingClosingBracketExpOut = "Скобка не открыта"

var TestMissingOpeningBracketExpOut = "Скобка не закрыта"


func TestUsualExp(t *testing.T) {
	result := Execute("48+42*18/63")

	if TestUsualExpOut != result {
		t.Errorf("TestUsualExp faild")
	}
}

func TestBracketsExp(t *testing.T) {
	result := Execute("12-(-2+2)*3")

	if TestBracketsExpOut != result {
		t.Errorf("TestBracketsExp faild")
	}
}

func TestUnaryExp(t *testing.T) {
	result := Execute("-2+2")

	if TestUnaryExpOut != result {
		t.Errorf("TestUnaryExp faild")
	}
}

func TestNotFinishedExp(t *testing.T) {
	result := Execute("2*2-7+")

	if TestNotFinishedExpOut != result {
		t.Errorf("TestNotFinishedExp faild")
	}
}

func TestOperatorsInARowExp(t *testing.T) {
	result := Execute("-2++2")

	if TestOperatorsInARowExpOut != result {
		t.Errorf("TestOperatorsInARowExp faild")
	}
}

func TestOperandsInARowExp(t *testing.T) {
	result := Execute("-2.2.2")

	if TestOperandsInARowExpOut != result {
		t.Errorf("TestOperandsInARowExp faild")
	}
}

func TestUnknownSymbolExp(t *testing.T) {
	result := Execute("-22;2")

	if TestUnknownSymbolExpOut != result {
		t.Errorf("TestUnknownSymbolExp faild")
	}
}

func TestMissingOpeningBracketExp(t *testing.T) {
	result := Execute("3+2)")

	if TestMissingClosingBracketExpOut != result {
		t.Errorf("TestMissingOpeningBracketExp faild")
	}
}

func TestMissingClosingBracketExp(t *testing.T) {
	result := Execute("2+(3+2")

	if TestMissingOpeningBracketExpOut != result {
		t.Errorf("TestMissingnClosingBracketExp faild")
	}
}