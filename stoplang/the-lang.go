package stoplang

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// TODO: might want to add scanner here and deal with state pollution
type Lang struct {
	hadErrors bool
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

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	l.run(string(bytes))

	if l.hadErrors {
		os.Exit(65)
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
	}
}

// TODO: wont work until add method for tokens isnt implemented
func (l *Lang) run(source string) {
	lex := NewScanner(source)

	tokens := lex.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func (l *Lang) Error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lang) report(line int, where string, message string) {
	fmt.Fprint(os.Stdin, "[line %d] errors%s: %s\n", line, where, message)
	l.hadErrors = true
}
