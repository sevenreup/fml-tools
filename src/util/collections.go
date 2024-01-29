package util

import "fhirLSP/src/lexer"

// check if item is in slice
func Contains(slice []lexer.Token, item lexer.Token) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}