package stoplang

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Lang struct {
	hadErrors       bool
	hadRuntimeError bool
	interpreter     Interpreter
}

func (l *Lang) main(args []string) {
	if len(args) > 1 {
		fmt.Println("usage: lang [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		l.runFile(args[0])
	} else {
		l.runPromt()
	}
}

func (l *Lang) runFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	l.run(string(bytes))

	if l.hadErrors {
		os.Exit(65)
	}

	if l.hadRuntimeError {
		os.Exit(70)
	}
}

func (l *Lang) runPromt() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()

		l.run(line)
		l.hadErrors = false
		l.hadRuntimeError = false
	}
}

func (l *Lang) run(source string) {
	lex := NewScanner(source)
	tokens := lex.ScanTokens()

	parser := &Parser{
		tokens: tokens,
		lang:   l,
	}

	statements := parser.Parse()

	if l.hadErrors {
		return
	}

	l.interpreter.Interpret(statements, l)
}

func (l *Lang) Error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lang) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	l.hadErrors = true
}

func (l *Lang) ErrorForParser(token Token, message string) {
	if token.tokenType == EOF {
		l.report(token.line, " at the end", message)
	} else {
		l.report(token.line, " at '"+token.lexeme+"'", message)
	}
}
