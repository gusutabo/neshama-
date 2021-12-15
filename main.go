package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Token int

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT
	SEMI
	OPERATOR
	SEPARATOR
)

var tokens = []string{
	EOF:       "EOF",
	ILLEGAL:   "ILLEGAL",
	IDENT:     "IDENT",
	INT:       "INT",
	OPERATOR:  "OPERATOR",
	SEPARATOR: "SEPARATOR",
}

func (t Token) String() string {
	return tokens[t]
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			panic(err)
		}

		l.pos.column++

		switch r {
		case '\n':
			l.resetPos()
		case ';':
			return l.pos, SEPARATOR, ";"
		case '+':
			return l.pos, OPERATOR, "+"
		case '-':
			return l.pos, OPERATOR, "-"
		case '*':
			return l.pos, OPERATOR, "*"
		case '/':
			return l.pos, OPERATOR, "/"
		case '=':
			return l.pos, OPERATOR, "="
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			} else if unicode.IsLetter(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()
				return startPos, IDENT, lit
			} else {
				return l.pos, ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) resetPos() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func main() {
	filename := flag.String("f", "/dev/stdin", "File to parse")
	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lexer := NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}

		fmt.Printf("%d:%d -> %s %s\n", pos.line, pos.column, tok, lit)
	}
}
