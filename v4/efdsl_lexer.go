package sameriver

import (
	"regexp"
	"strings"
	"text/scanner"
	"unicode"
)

var DSLIdentRuneRe = regexp.MustCompile("[A-Za-z0-9.<>\\[\\]{}]")

func DSLIdentRune(r rune) bool {
	return DSLIdentRuneRe.MatchString(string(r))
}

type EFDSLToken int

const (
	EOF EFDSLToken = iota
	Not
	And
	Or
	Function
	Identifier
	OpenParen
	CloseParen
	Comma
	Semicolon
)

func (t EFDSLToken) String() string {
	switch t {
	case EOF:
		return "EOF"
	case Not:
		return "Not"
	case And:
		return "And"
	case Or:
		return "Or"
	case Function:
		return "Function"
	case Identifier:
		return "Identifier"
	case OpenParen:
		return "OpenParen"
	case CloseParen:
		return "CloseParen"
	case Comma:
		return "Comma"
	case Semicolon:
		return "Semicolon"
	default:
		return "Unknown"
	}
}

type EFDSLLexer struct {
	scanner.Scanner
	token       EFDSLToken
	stringValue string
}

func (l *EFDSLLexer) IsEOF() bool {
	return l.Peek() == scanner.EOF
}

func (l *EFDSLLexer) TokenText() string {
	return l.stringValue
}

func (l *EFDSLLexer) Lex() EFDSLToken {
	l.stringValue = ""
	l.token = EOF

	for !l.IsEOF() {
		r := l.Peek()

		if unicode.IsSpace(r) {
			l.Next()
			continue
		}

		switch {
		case r == '!':
			l.Next()
			l.token = Not
		case r == '&':
			l.Next()
			if l.Peek() == '&' {
				l.Next()
				l.token = And
			} else {
				l.token = EOF
			}
		case r == '|':
			l.Next()
			if l.Peek() == '|' {
				l.Next()
				l.token = Or
			} else {
				l.token = EOF
			}
		case r == '(':
			l.Next()
			l.token = OpenParen
		case r == ')':
			l.Next()
			l.token = CloseParen
		case r == ',':
			l.Next()
			l.token = Comma
		case r == ';':
			l.Next()
			l.token = Semicolon
		case unicode.IsUpper(r):
			str := l.scanString(func(r rune) bool {
				return unicode.IsLetter(r)
			})
			if str != "" {
				l.stringValue = str
				l.token = Function
			} else {
				l.token = EOF
			}
		case DSLIdentRune(r):
			str := l.scanString(DSLIdentRune)
			if str != "" {
				l.stringValue = str
				l.token = Identifier
			} else {
				l.token = EOF
			}
		default:
			l.token = EOF
			l.Next()
		}

		if l.token != EOF {
			break
		}
	}

	return l.token
}

func (l *EFDSLLexer) scanString(isValid func(rune) bool) string {
	var buf strings.Builder
	for !l.IsEOF() && isValid(l.Peek()) {
		buf.WriteRune(l.Next())
	}
	return buf.String()
}
