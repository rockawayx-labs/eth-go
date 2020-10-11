package eth

import "strings"

// SanitizeHex removes the prefix `0x` if it exists
// and ensures there is an even number of characters in the string,
// padding on the left of the string is it's not the case.
func SanitizeHex(input string) string {
	if Has0xPrefix(input) {
		input = input[2:]
	}

	if len(input)%2 != 0 {
		input = "0" + input
	}

	return strings.ToLower(input)
}

func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}
