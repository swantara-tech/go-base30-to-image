package decoder

import (
	"strings"
)

// jSignature Base30 character set
// Full charset: "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWX"
// base = 30 (half of 60)
const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWX"
const base = 30 // len(charset) / 2

// qMap: lower → upper (for tail-char remapping during encode)
// mMap: upper → lower (for decode: upper chars are "tail" remapped to lower)
var qMap map[rune]rune
var mMap map[rune]rune

func init() {
	qMap = make(map[rune]rune)
	mMap = make(map[rune]rune)

	runes := []rune(charset)
	for r := base - 1; r >= 0; r-- {
		qMap[runes[r]]        = runes[r+base]
		mMap[runes[r+base]]   = runes[r]
	}
}

// UncompressLeg decodes one stroke leg (X or Y array) from jSignature base30 string.
// Rules (ported from jSignature source):
//   - 'Z' → flip sign to -1, reset buffer
//   - 'Y' → flip sign to +1, reset buffer
//   - char in lowerSet → part of current number (direct digit in base-30)
//   - char in mMap (uppercase/tail) → remap to lower, part of current number
//   - flush buffer when encountering Z / Y / start-of-new-number char (in lowerSet)
func UncompressLeg(s string) []int {
	chars := []rune(s)
	result := []int{}
	sign := 1
	buf := []rune{}
	last := 0

	flush := func() {
		if len(buf) == 0 {
			return
		}
		// parse buf as base-30 number
		val := 0
		for _, ch := range buf {
			val = val*base + charVal(ch)
		}
		val = val*sign + last
		result = append(result, val)
		last = val
		buf = []rune{}
	}

	for _, ch := range chars {
		if ch == 'Z' {
			flush()
			sign = -1
		} else if ch == 'Y' {
			flush()
			sign = 1
		} else if _, isLower := qMap[ch]; isLower {
			// lower char → flush previous number, start new token
			flush()
			buf = []rune{ch}
		} else if lower, isTail := mMap[ch]; isTail {
			// upper/tail char → remap to lower, continue current token
			buf = append(buf, lower)
		}
		// ignore any other chars
	}
	flush()

	return result
}

// charVal returns the base-30 value of a lower-set character
func charVal(ch rune) int {
	for i, r := range charset[:base] {
		if r == ch {
			return i
		}
	}
	return 0
}

// UncompressStrokes parses full base30 data string into strokes.
// Each stroke is represented as {X: []int, Y: []int}.
// Strokes are separated by '_', with xLeg and yLeg alternating.
func UncompressStrokes(data string) []Stroke {
	// Strip prefix
	lower := strings.ToLower(data)
	if strings.HasPrefix(lower, "image/jsignature;base30,") {
		data = data[len("image/jsignature;base30,"):]
	}

	parts := strings.Split(data, "_")
	strokes := []Stroke{}

	// Parts come in pairs: xLeg, yLeg
	for i := 0; i+1 < len(parts); i += 2 {
		xs := UncompressLeg(parts[i])
		ys := UncompressLeg(parts[i+1])
		if len(xs) > 0 && len(ys) > 0 {
			strokes = append(strokes, Stroke{X: xs, Y: ys})
		}
	}

	return strokes
}

// Stroke holds decoded X and Y coordinate arrays for a single pen stroke
type Stroke struct {
	X []int
	Y []int
}
