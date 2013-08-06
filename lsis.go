package main

import (
	"bufio"
	"fmt"
	"github.com/mischief/lsystem"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
/*
	variables = []rune{'A', 'B'}
	constants = []rune{'-', '+'}

	rules = map[rune]string{
		'A': "B-A-B",
		'B': "A+B+A",
	}

	tg = map[rune]string{
		'A': "drawForward 10",
		'B': "drawForward 10",
		'+': "turnLeft 60",
		'-': "turnRight 60",
	}

	tgfuncs = map[string]lsystem.TGFunc{
		"drawForward": DrawForward,
		"turnLeft":    TurnLeft,
		"turnRight":   TurnRight,
	}
*/
)

func FirstRune(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

func main() {
	var vars lsystem.Variables
	var cons lsystem.Constants

	infile := os.Stdin

	if len(os.Args) == 2 {
		var err error
		if infile, err = os.Open(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s", os.Args[1], err)
			os.Exit(1)
		}
	}

	rules := lsystem.NewRules()

	tgrules := lsystem.NewTurtleGraphicsRules()

	ls := lsystem.NewLSystem("", &vars, &cons, rules)
	tg := lsystem.NewTurtleGraphics(640, 480, tgrules)

	input := bufio.NewScanner(infile)

	for input.Scan() {
		line := input.Text()

		// should always have at least one field, because of bufio.Scanner
		fields := strings.Fields(line)

		if FirstRune(line) == '#' || len(fields) < 1 {
			continue
		}

		switch fields[0] {
		case "start":
			if len(fields) != 2 {
				fmt.Fprintln(os.Stderr, "invalid call to %s", fields[0])
				os.Exit(1)
			}
			ls.SetState(fields[1])
		case "addvar":
			if len(fields) != 2 {
				fmt.Fprintln(os.Stderr, "invalid call to %s", fields[0])
				os.Exit(1)
			}
			vars.Add(FirstRune(fields[1]))
		case "addconst":
			if len(fields) != 2 {
				fmt.Fprintln(os.Stderr, "invalid call to %s", fields[0])
				os.Exit(1)
			}
			cons.Add(FirstRune(fields[1]))
		case "addrule":
			if len(fields) != 3 {
				fmt.Fprintln(os.Stderr, "invalid call to %s", fields[0])
				os.Exit(1)
			}

			rules.Add(FirstRune(fields[1]), fields[2])

		case "step":
			if len(fields) != 2 {
				fmt.Fprintln(os.Stderr, "invalid call to %s", fields[0])
				os.Exit(1)
			}
			count, err := strconv.Atoi(fields[1])
			if err != nil {
				fmt.Fprintln(os.Stderr, "invalid step count")
				os.Exit(1)
			}
			ls.Run(count)

		case "tgaddrule":
			if len(fields) != 4 {
				fmt.Fprintln(os.Stderr, "invalid call to %s", fields[0])
				os.Exit(1)
			}

			fun := fields[2]
			var fn lsystem.TGFunc
			switch fun {
			case "drawfwd":
				fn = lsystem.DrawFwd
			case "turn":
				fn = lsystem.Turn
			default:
				fmt.Fprintln(os.Stderr, "invalid function ", fields[2])
			}

			num, err := strconv.Atoi(fields[3])
			if err != nil {
				fmt.Fprintln(os.Stderr, "invalid number ", fields[3])
			}

			tgrules.Add(FirstRune(fields[1]), fn, num)

		case "tgsave":
			if len(fields) != 2 {
				fmt.Fprintln(os.Stderr, "invalid call to %s", fields[0])
				os.Exit(1)
			}
			tg.Draw(ls)
			if err := tg.SavePNG(fields[1]); err != nil {
				fmt.Fprintln(os.Stderr, "%s failed: %s", fields[0], err)
			} else {
				fmt.Printf("saved to %s\n", fields[1])
			}

		default:
			fmt.Fprintln(os.Stderr, "invalid command: ", fields[0])

		}
	}

}
