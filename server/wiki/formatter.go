package wiki

import (
	"fmt"
	"html"
	"io"
	"strings"
)

// Formats wiki syntax into HTML.
func wikiFormatter(w io.Writer, s string, data ...interface{}) {
	b, ok := data[0].([]byte)
	if !ok {
		panic("Invalid data passed to wikiFormatter!")
	}
	parse(tokenize(string(b)), w)
}

// Formats wiki syntax into HTML.
func tokenFormatter(w io.Writer, s string, data ...interface{}) {
	b, ok := data[0].([]byte)
	if !ok {
		panic("Invalid data passed to tokenFormatter!")
	}
	w.Write([]byte(html.EscapeString(fmt.Sprintln(tokenize(string(b))))))
}

type tokenCode int

const (
	RAW tokenCode = iota
	NOWIKI_START
	NOWIKI_END
	LINE_FEED
	BULLET
	LINK_START
	LINK_END
)

func (tc tokenCode) String() string {
	ret := "Unknown token type."
	switch tc {
	case RAW:
		ret = "RAW"
	case NOWIKI_START:
		ret = "NOWIKI_START"
	case NOWIKI_END:
		ret = "NOWIKI_END"
	case LINE_FEED:
		ret = "LINE_FEED"
	case BULLET:
		ret = "BULLET"
	case LINK_START:
		ret = "LINK_START"
	case LINK_END:
		ret = "LINK_END"
	}
	return ret
}

type token struct {
	code tokenCode
	string
}

var key = []token{token{NOWIKI_START, "<nowiki>"},
	token{NOWIKI_END, "</nowiki>"},
	token{LINE_FEED, "\r\n"}, // This should treat Windows' new lines the same as
	token{LINE_FEED, "\n"},   // the standard ones.
	token{BULLET, "*"},       // Update parseList if changing bullets.
	token{BULLET, "#"},       // Update parseList if changing bullets.
	token{LINK_START, "["},
	token{LINK_END, "]"}}

func tokenize(s string) []token {
	charsPer := 25 // should be an acceptable guess; only affects performance
	tokens := make([]token, 0, len(s)/charsPer)
	used, at := 0, 0

outer:
	for used != len(s) {
		for _, v := range key {
			if strings.HasPrefix(s[at:], v.string) {
				if at != used {
					tokens = append(tokens, token{RAW, s[used:at]})
					used = at
				}
				tokens = append(tokens, v)
				used += len(v.string)
				at = used
				continue outer
			}
		}
		at++
		if at == len(s) {
			tokens = append(tokens, token{RAW, s[used:]})
			break
		}
	}
	return tokens
}

func parse(t []token, w io.Writer) {
	var bullets, lineEnd string
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
			add("\n\nPlease report this bug to Mark.\n</pre>\n")
		}
	}()
outer:
	for i := 0; i < len(t); i++ {
		v := t[i]
		switch v.code {
		case RAW:
			add(v.string)
		case NOWIKI_START: // treat anything in <nowiki> tags as raw text
			for j := i + 1; j < len(t); j++ {
				if t[j].code == NOWIKI_END {
					i = j
					break
				} else {
					add(t[j].string)
				}
			}
		case LINE_FEED: // show each consecutive line feed after the first
			add(lineEnd + "\n") // so other cases (currently only lists) can have stuff appear at the end of their line
			lineEnd = ""
			if bullets != "" && i+1 != len(t) && t[i+1].string[0] != bullets[0] {
				add(parseList(&bullets, ""))
			}
			for j := i + 1; j < len(t); j++ {
				if t[j].code == LINE_FEED {
					add("<br/>\n")
				} else {
					i = j - 1
					continue outer
				}
			}
			break outer // used up all tokens if "continue outer" not reached
		case BULLET:
			if !(i == 0 || t[i-1].code == LINE_FEED) { // not start of line
				add(v.string)
				break
			}
			current := v.string
			for j := i + 1; j < len(t); j++ {
				if t[j].code == BULLET {
					current += t[j].string
				} else {
					i = j - 1
					break
				}
			}
			add(parseList(&bullets, current) + "<li>")
			lineEnd += "</li>"
		case LINK_START:
			// TODO: Rewrite parser to handle multiple tokens inside a link or rewrite tokenizer to prevent multiple tokens inside the same link.
			if i+2 >= len(t) || t[i+1].code != RAW || t[i+2].code != LINK_END {
				add(v.string) // Display extraneous link starts as raw.
				continue outer
			}
			add(parseLink(t[i+1].string))
			i += 2 // RAW and LINK_END tokens
		default:
			add(v.string) // display all extraneous formatting as raw
		}
	}
	add(lineEnd)
}

func parseList(prev *string, current string) string {
	var from int
	var short, long string
	if len(current) > len(*prev) {
		long, short = current, *prev
	} else {
		short, long = current, *prev
	}
	for j := 0; j <= len(short); j++ {
		if j == len(short) || long[j] != short[j] {
			from = j
			break
		}
	}
	bText := ""
	for j := len(*prev) - 1; j >= from; j-- {
		bText += map[byte]string{'*': "</ul>", '#': "</ol>"}[(*prev)[j]]
	}
	for j := from; j < len(current); j++ {
		bText += map[byte]string{'*': "<ul>", '#': "<ol>"}[current[j]]
	}
	*prev = current
	return bText
}

func parseLink(s string) string {
	at := strings.Index(s, ":")
	var place, display, protocol, separator string
	if at == -1 {
		protocol, separator = "", "|"
	} else {
		protocol, separator = s[:at], " "
	}
	if i := strings.Index(s, separator); i == -1 {
		place, display = s, s
	} else {
		place, display = s[:i], s[i+1:]
	}
	place = html.EscapeString(place)
	display = html.EscapeString(display)
	var ret string
	if protocol != "" || place[0] == '/' {
		ret = "<a href=\"" + place + "\">" + display + "</a>"
	} else {
		wikiName := "wiki" // TODO: refactor to get this from Wiki.Prefix
		ret = "<a href=\"/" + wikiName + "/view/" + place + "\">" + display + "</a>"
	}
	return ret
}
