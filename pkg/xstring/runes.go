package xstring

import "slices"

func IsNumber(c rune) bool {
	return c >= '0' && c <= '9'
}

func IsLowerCaseLetter(c rune) bool {
	return c >= 'a' && c <= 'z'
}

func IsUpperCaseLetter(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func IsLetter(c rune) bool {
	return IsLowerCaseLetter(c) || IsUpperCaseLetter(c)
}

func IsWhiteSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func IsSpace(c rune) bool {
	return c == ' '
}

func IsUnderscore(c rune) bool {
	return c == '_'
}

func IsSpecialCharacter(c rune) bool {
	return slices.Contains([]rune{
		' ',
		'~', '`', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '-', '+', '=',
		'[', ']', '{', '}', '\\', '|',
		':', ';', '"', '\'',
		'<', '>', ',', '.', '?', '/',
	}, c)
}
