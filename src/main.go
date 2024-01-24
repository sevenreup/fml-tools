package main

import (
	"fhirLSP/src/formatter"
)

func main() {
	formatter := formatter.NewFormatter()
	formatter.Format()
}