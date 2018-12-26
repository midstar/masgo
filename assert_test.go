// Assertion helpers of golang unit tests
package main

import (
	"fmt"
	"runtime/debug"
	"testing"
)

func assertTrue(t *testing.T, message string, check bool) {
	if !check {
		debug.PrintStack()
		t.Fatal(message)
	}
}

func assertEqualsInt(t *testing.T, message string, expected int, actual int) {
	assertTrue(t, fmt.Sprintf("%s\nExpected: %d, Actual: %d", message, expected, actual), expected == actual)
}

func assertEqualsStr(t *testing.T, message string, expected string, actual string) {
	assertTrue(t, fmt.Sprintf("%s\nExpected: %s, Actual: %s", message, expected, actual), expected == actual)
}

func assertEqualsSlice(t *testing.T, message string, expected []uint32, actual []uint32) {
	assertEqualsInt(t, fmt.Sprintf("%s\nSize missmatch", message), len(expected), len(actual))
	for index, expvalue := range expected {
		actvalue := actual[index]
		assertTrue(t, fmt.Sprintf("%s\nIndex %d - Expected: %d, Actual: %d", message, index, expvalue,
			actvalue), expvalue == actvalue)
	}
}
