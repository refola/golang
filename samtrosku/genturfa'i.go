// genturfa'i.go
// x1 parses according to formal grammer x2 text x3
// x1 = gemturfa'i.go
// x2 = lojban's formal grammar
// x3 = input

// This is based off of server/wiki/formatter.go.
// TODO: implement more of teh lojban grammarz
/*
TODO: convert to syntax tree parsing, like this:
bridi[
	sumti[
		cmavo[
			=convert-selbri-to-sumti
		]
		selbri
		cmavo[
			=terminator
		]
	]
	selbri[
		=tanru[
			brivla
			brivla[
				=lujvo[
					rafsi
					rafsi
					rafsi
				]
			]
		]
	]
	sumti
]
*/
// NOTE list for future upgrader:
// * cmavo context should find fu'ivla, or "OTHER" catch-all will catch them all.
// * lujvo -> rafsi+ might be separate substage miniparser.
// * It's probably wayyy easier to (keep the parse-to-word-by-morphology parser, discard spaces, and use grammar to convert to trees) than it is to do it all in one place. It can still be one-pass by making the morphology-to-word-type parser emit tokens on a channel instead of generating a list. Streams and trees, dude. Streams and trees.

package samtrosku

import (
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"
)

// How to match text to token types? Use patterns of some sort, like simplified regex
type tokenType struct {
	pattern *regexp.Regexp
	string
}

// regex notes:
// 	Anything in [CV'] gets turned into custom character class before regex compile. These are used in Lojban grammar docs and tutorials.
// 	\b = "word boundary" and \B = "not \b" should surround Lojban word types, with \b next to (C|V) and \B next to ('|.). Some day the preparsing of the parsing patterns may use regex to avoid this need, becoming more directly and automatically meta.
// 	Make sure to double "\" into "\\" or use raw string literals.
var (
	SPACE       = &tokenType{parseTokenTypePattern(" "), "SPACE"}
	ATTITUDINAL = &tokenType{parseTokenTypePattern(`\B.V'?V.\B`), "ATTITUDINAL"}
	SELBRI      = &tokenType{parseTokenTypePattern(`\bCCVCV|CVCCV\b`), "SELBRI"}
	OTHER       = &tokenType{parseTokenTypePattern("[^ ]+"), "OTHER"}
)

// Convert a pattern string to a regex automaton.
// Pattern syntax is regex, with the following characters -> character class conversions:
// 	C = Lojban consonant = [bcdfgjklmnprstvxzBCDFGJKLMNPRSTVXZ]
// 	V = Lojban vowel = [aeiouAEIOU]
// 	' = Lojban vowel separator = ['h]
func parseTokenTypePattern(s string) *regexp.Regexp {
	type replacement struct {
		old string
		new string
	}
	replacements := []replacement{
		{"C", "[bcdfgjklmnprstvxzBCDFGJKLMNPRSTVXZ]"},
		{"V", "[aeiouAEIOU]"},
		{"'", "['h]"}}
	for _, rpl := range replacements {
		s = strings.Replace(s, rpl.old, rpl.new, -1)
	}
	return regexp.MustCompile(s)
}

// attempt to find an instance of t at the beginning of s
// returns length of successful match, or unspecified nonpositive number if unsuccessful
func (t *tokenType) attemptMatch(s string) int {
	if loc := t.pattern.FindStringIndex(s); loc != nil {
		return loc[1]
	}
	return -1
}

// prioritization of which token type to check for first
// TODO: sort by Lojban morphology rules IF applicable (Lojban doesn't have any ambiguity about phonology -> word type, right?), otherwise sort by frequency of usage in Lojban text sample
var tokenTypeParseOrder = []*tokenType{
	SPACE, // should be first for performance reasons; is right half the time
	ATTITUDINAL,
	SELBRI,
	OTHER} // should be last as it consumes everything

type token struct {
	*tokenType
	string
}

func tokenize(s string) []token {
	charsPer := 3 // Lojban uses a lot of short words and we're tokenizing spaces, so this assumes charsPer*2-1 characters per non-space tokens. Too high and the runtime will reallocate the token slice. Too low and excess memory will be consumed.
	tokens := make([]token, 0, len(s)/charsPer)
	at := 0

outer:
	for at != len(s) {
		for _, v := range tokenTypeParseOrder {
			if matchLen := v.attemptMatch(s[at:]); matchLen > 0 {
				tokens = append(tokens, token{v, s[at : at+matchLen]}) // TODO: check for off-by-1
				at += matchLen
				continue outer
			}
		}
		panic("unreachable line") // catch-all OTHER token should catch all
	}

	return tokens
}

// generates HTML tagging for use with syntax-highlighting CSS
// TODO: Write parse functions that generate, e.g., English or computer code.
func parseToHTML(t []token, w io.Writer) {
	add := func(s string) {
		w.Write([]byte(s))
	}
	esc := func(s string) {
		add(html.EscapeString(s))
	}
	defer func() {
		if err := recover(); err != nil {
			add("<pre>Parsing error! Details follow.\n\n")
			esc(fmt.Sprintln(err))
			add("\n\nTokens are as follows:\n\n")
			esc(fmt.Sprintln(t))
			add("\n\nPlease report this bug at github.com/refola with \"samtrosku\" in the title.\n</pre>\n")
		}
	}()

	for _, v := range t {
		switch v.tokenType {
		case SPACE:
			add(" ")
		case OTHER:
			esc(v.string)
		default:
			add("<div class=\"" + v.tokenType.string + ">" + v.string + "</div>")
		}
	}
}
