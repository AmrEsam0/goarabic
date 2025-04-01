// Package goarabic contains utility functions for working with Arabic strings.
package goarabic

import "strings"

// Reverse returns its argument string reversed rune-wise left to right.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// SmartLength returns the length of the given string
// without considering the Arabic Vowels (Tashkeel).
func SmartLength(s *string) int {
	// len() use int as return value, so we'd better follow for compatibility
	length := 0

	for _, value := range *s {
		if tashkeel[value] {
			continue
		}
		length++
	}

	return length
}

// RemoveTashkeel returns its argument as rune-wise string without Arabic vowels (Tashkeel).
func RemoveTashkeel(s string) string {
	// var r []rune
	// the capcity of the slice wont be greater than the length of the string itself
	// hence, cap = len(s)
	r := make([]rune, 0, len(s))

	for _, value := range s {
		if tashkeel[value] {
			continue
		}
		r = append(r, value)
	}

	return string(r)
}

// RemoveTatweel returns its argument as rune-wise string without Arabic Tatweel character.
func RemoveTatweel(s string) string {
	r := make([]rune, 0, len(s))

	for _, value := range s {
		if TATWEEL.equals(value) {
			continue
		}
		r = append(r, value)
	}

	return string(r)
}

func getCharGlyph(previousChar, currentChar, nextChar rune) rune {
	glyph := currentChar
	previousIn := false // in the Arabic Alphabet or not
	nextIn := false     // in the Arabic Alphabet or not

	for _, s := range alphabet {
		if s.equals(previousChar) { // previousChar in the Arabic Alphabet ?
			previousIn = true
		}

		if s.equals(nextChar) { // nextChar in the Arabic Alphabet ?
			nextIn = true
		}
	}

	for _, s := range alphabet {

		if !s.equals(currentChar) { // currentChar in the Arabic Alphabet ?
			continue
		}

		if previousIn && nextIn { // between two Arabic Alphabet, return the medium glyph
			for s, _ := range beggining_after {
				if s.equals(previousChar) {
					return getHarf(currentChar).Beggining
				}
			}

			return getHarf(currentChar).Medium
		}

		if nextIn { // beginning (because the previous is not in the Arabic Alphabet)
			return getHarf(currentChar).Beggining
		}

		if previousIn { // final (because the next is not in the Arabic Alphabet)
			for s, _ := range beggining_after {
				if s.equals(previousChar) {
					return getHarf(currentChar).Isolated
				}
			}
			return getHarf(currentChar).Final
		}

		if !previousIn && !nextIn {
			return getHarf(currentChar).Isolated
		}

	}
	return glyph
}

// equals() return if true if the given Arabic char is alphabetically equal to
// the current Harf regardless its shape (Glyph)
func (c *Harf) equals(char rune) bool {
	switch char {
	case c.Unicode:
		return true
	case c.Beggining:
		return true
	case c.Isolated:
		return true
	case c.Medium:
		return true
	case c.Final:
		return true
	default:
		return false
	}
}

// getHarf gets the correspondent Harf for the given Arabic char
func getHarf(char rune) Harf {
	for _, s := range alphabet {
		if s.equals(char) {
			return s
		}
	}

	return Harf{Unicode: char, Isolated: char, Medium: char, Final: char}
}

// RemoveAllNonAlphabetChars deletes all characters which are not included in Arabic Alphabet
func RemoveAllNonArabicChars(text string) string {
	runes := []rune(text)
	newText := []rune{}
	for _, current := range runes {
		inAlphabet := false
		for _, s := range alphabet {
			if s.equals(current) {
				inAlphabet = true
			}
		}
		if inAlphabet {
			newText = append(newText, current)
		}
	}
	return string(newText)
}

// ToGlyph returns the glyph representation of the given text
func ToGlyph(text string) string {
	var prev, next rune

	runes := []rune(text)
	length := len(runes)
	newText := make([]rune, 0, length)

	for i, current := range runes {
		// get the previous char
		if (i - 1) < 0 {
			prev = 0
		} else {
			prev = runes[i-1]
		}

		// get the next char
		if (i + 1) <= length-1 {
			next = runes[i+1]
		} else {
			next = 0
		}

		// get the current char representation or return the same if unnecessary
		glyph := getCharGlyph(prev, current, next)

		// append the new char representation to the newText
		newText = append(newText, glyph)
	}

	return string(newText)
}

// RemoveTashkeel returns its argument as rune-wise string without Arabic vowels (Tashkeel).
/*
func RemoveTashkeelExtended(s string) string {
	r := []rune(s)

	m := map[string]bool{"\u064e": true, "\u064b": true, "\u064f": true,
		"\u064c": true, "\u0650": true, "\u064d": true,
		"\u0651": true, "\u0652": true}

	for key, value := range s {
		if m[value] {
			continue
		}
		r[key] = value
	}

	return string(r)
}
*/

var isArabic map[rune]bool

func fillIsArabicMap() {
	if isArabic != nil {
		return
	}
	isArabic = make(map[rune]bool)
	for r := rune(0x0600); r <= rune(0x06FF); r++ {
		isArabic[r] = true
	}
	for r := rune(0x0750); r <= rune(0x077F); r++ {
		isArabic[r] = true
	}
	for r := rune(0x08A0); r <= rune(0x08FF); r++ {
		isArabic[r] = true
	}
	for r := rune(0xFB50); r <= rune(0xFDFF); r++ {
		isArabic[r] = true
	}
	for r := rune(0xFE70); r <= rune(0xFEFF); r++ {
		isArabic[r] = true
	}
	for r := rune(0x10E60); r <= rune(0x10E7F); r++ {
		isArabic[r] = true
	}
}

// FixBidiText fixes the bidirectional text for Arabic, English and numbers in a string.
func FixBidiText(text string) string {
	if len(text) == 0 {
		return text
	}

	fillIsArabicMap()
	words := strings.Fields(text)
	var result []string

	// Process each word
	for _, word := range words {
		runes := []rune(word)
		isArabicWord := false
		for _, r := range runes {
			if isArabic[r] {
				isArabicWord = true
				break
			}
		}

		if isArabicWord {
			// Apply Arabic transformations
			result = append(result, Reverse(ToGlyph(word)))
		} else {
			// Keep English words as is (we'll handle their order later)
			result = append(result, word)
		}
	}

	// Reverse the entire sentence for RTL (Arabic) display
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	// Now, group and reverse consecutive English words
	var finalResult []string
	var englishGroup []string

	for _, word := range result {
		runes := []rune(word)
		isEnglishWord := true
		for _, r := range runes {
			if isArabic[r] {
				isEnglishWord = false
				break
			}
		}

		if isEnglishWord {
			englishGroup = append(englishGroup, word)
		} else {
			// If we have an English group, reverse its word order before adding Arabic
			if len(englishGroup) > 0 {
				// Reverse the English group's word order
				for i, j := 0, len(englishGroup)-1; i < j; i, j = i+1, j-1 {
					englishGroup[i], englishGroup[j] = englishGroup[j], englishGroup[i]
				}
				finalResult = append(finalResult, englishGroup...)
				englishGroup = nil
			}
			finalResult = append(finalResult, word)
		}
	}

	// Add any remaining English words
	if len(englishGroup) > 0 {
		for i, j := 0, len(englishGroup)-1; i < j; i, j = i+1, j-1 {
			englishGroup[i], englishGroup[j] = englishGroup[j], englishGroup[i]
		}
		finalResult = append(finalResult, englishGroup...)
	}

	return strings.Join(finalResult, " ")
}
