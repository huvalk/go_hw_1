package main

import (
	"github.com/stretchr/testify/require"
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

	require.Equal(t, TestUsualExpOut, result, "TestUsualExp faild")
}

func TestBracketsExp(t *testing.T) {
	result := Execute("12-(-2+2)*3")

	require.Equal(t, TestBracketsExpOut, result, "TestBracketsExp faild")
}

func TestUnaryExp(t *testing.T) {
	result := Execute("-2+2")

	require.Equal(t, TestUnaryExpOut, result, "TestUnaryExp faild")
}

func TestNotFinishedExp(t *testing.T) {
	result := Execute("2*2-7+")

	require.Equal(t, TestNotFinishedExpOut, result, "TestNotFinishedExp faild")
}

func TestOperatorsInARowExp(t *testing.T) {
	result := Execute("-2++2")

	require.Equal(t, TestOperatorsInARowExpOut, result, "TestOperatorsInARowExp faild")
}

func TestOperandsInARowExp(t *testing.T) {
	result := Execute("-2.2.2")

	require.Equal(t, TestOperandsInARowExpOut, result, "TestOperandsInARowExp faild")
}

func TestUnknownSymbolExp(t *testing.T) {
	result := Execute("-22;2")

	require.Equal(t, TestUnknownSymbolExpOut, result, "TestUnknownSymbolExp faild")
}

func TestMissingOpeningBracketExp(t *testing.T) {
	result := Execute("3+2)")

	require.Equal(t, TestMissingClosingBracketExpOut, result, "TestMissingOpeningBracketExp faild")
}

func TestMissingClosingBracketExp(t *testing.T) {
	result := Execute("2+(3+2")

	require.Equal(t, TestMissingOpeningBracketExpOut, result, "TestMissingnClosingBracketExp faild")
}