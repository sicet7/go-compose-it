package term

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"io"
	"strings"
)

func StringPrompt(label string, in io.Reader, out io.Writer) string {
	var s string
	r := bufio.NewReader(in)
	for {
		fmt.Fprintln(out, label+"")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func PasswordPrompt(label string, inFd int, out io.Writer) string {
	var s string
	for {
		_, err := fmt.Fprint(out, label+" ")
		if err != nil {
			return ""
		}
		b, _ := term.ReadPassword(inFd)
		s = string(b)
		if s != "" {
			break
		}
	}
	fmt.Println()
	return s
}
