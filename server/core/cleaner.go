// TODO: Replace use of this file with "encoding/hex" in the standard libraries? It will make things take a bit more space and be less readable, but drastically reduces the amount of code that I need to maintain.

package core

import (
	"regexp"
)

var cleaner = regexp.MustCompile("[^a-zA-Z0-9/]")
var uncleaner = regexp.MustCompile("%[a-fA-F0-9][a-fA-F0-9]")

const hexDigits = "0123456789ABCDEF"

// converts a byte to two hex characters representing the first and second nibbles of the byte
func escapeByte(b byte) (byte, byte) {
	first := b >> 4
	last := b & 0x0F
	return hexDigits[first], hexDigits[last]
}

// converts two hex characters representing the first and second nibbles of the byte to a byte
func unescapeByte(first, last byte) byte {
	decode := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		default:
			panic("invalid hex digit")
		}
		panic("should not be able to reach this")
	}
	return decode(first)<<4 + decode(last)
}

func cleanChars(unclean []byte) []byte {
	clean := make([]byte, len(unclean)*3)
	for i, v := range unclean {
		clean[3*i] = '%'
		clean[3*i+1], clean[3*i+2] = escapeByte(v)
	}
	return clean
}
func uncleanChars(clean []byte) []byte {
	if len(clean)%3 != 0 {
		panic("charUncleaner given non-multiple of 3 bytes")
	}
	unclean := make([]byte, len(clean)/3)
	for i, _ := range unclean {
		unclean[i] = unescapeByte(clean[3*i+1], clean[3*i+2])
	}
	return unclean
}

// makes data names safe for file systems and URLs
func Sanitize(uncleanText string) string {
	return string(cleaner.ReplaceAllFunc([]byte(uncleanText), cleanChars))
}

// reverses sanitize
func Unsanitize(cleanText string) string {
	return string(uncleaner.ReplaceAllFunc([]byte(cleanText), uncleanChars))
}
