package env

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/roeldev/go-errs"
)

const (
	runeQuot    = 39 // '
	runeDblQuot = 34 // "
	runeHash    = 35 // #
)

// Map represents a map of env key value pairs.
type Map map[string]string

func (m Map) parse(env string) bool {
	key, val := ParsePair(env)
	if key != "" && val != "" {
		m[key] = val
		return true
	}

	return false
}

// Merge any map of strings with this Map.
func (m Map) Merge(e map[string]string) {
	for k, v := range e {
		m[k] = v
	}
}

// Environ returns a `Map` with the os' current environment variables.
func Environ() (Map, int) {
	s := os.Environ()
	m := make(Map, len(s))
	n := ParseSlice(s, m)
	return m, n
}

// Read from an `io.Reader`, parse its results and parse them to the provided Map. Each line is
// cleaned before being parsed with `ParsePair`.
// It returns the number of parsed lines and any error that occurs while scanning for lines.
func Read(r io.Reader, dest Map) (int, error) {
	n := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == runeHash {
			continue // skip empty lines and comments
		}
		if dest.parse(cleanLine(line)) {
			n++
		}
	}

	return n, errs.Wrap(scanner.Err())
}

// ParseSlice parses a slice of strings to a map with key value pairs. The slice should be clean,
// entries are not checked on starting/trailing whitespace or comment tags.
// It returns the number of parsed lines.
func ParseSlice(env []string, dest Map) (n int) {
	for _, e := range env {
		if dest.parse(e) {
			n++
		}
	}

	return n
}

// ParseFlagArgs parses the
func ParseFlagArgs(flag string, args []string, dest Map) (n int) {
	if flag == "" || len(args) == 0 {
		return 0
	}

	sd, dd := "-"+flag, "--"+flag
	sdl, ddl := len(sd), len(dd)

	nextIsPair := false
	for _, arg := range args {
		if nextIsPair && arg[0] != '-' {
			if dest.parse(arg) {
				n++
			}
			nextIsPair = false
			continue
		}

		if strings.Index(arg, sd) == 0 {
			if len(arg) == sdl {
				nextIsPair = true
			} else if dest.parse(arg[sdl+1:]) {
				n++
			}
		} else if strings.Index(arg, dd) == 0 {
			if len(arg) == ddl {
				nextIsPair = true
			} else if dest.parse(arg[ddl+1:]) {
				n++
			}
		}
	}

	return n
}

// ParsePair parses an environment variable's paired key and value and returns the separated values.
func ParsePair(pair string) (key string, val string) {
	// environment variables can begin with = so start splitting after the
	// first character
	split := strings.SplitAfterN(pair[1:], "=", 2)
	if len(split) != 2 {
		return "", ""
	}

	key = string(pair[0]) + split[0]
	// after splitting, the last character of key is '='
	// obviously we do not want to keep it
	key = key[0 : len(key)-1]

	val = split[1]
	if val == "" {
		return
	}

	last := len(val) - 1

	// remove optional double quotes from string
	if (val[0] == runeDblQuot || val[0] == runeQuot) && val[0] == val[last] {
		val = val[1:last]
	}

	return
}

func cleanLine(line string) string {
	return line
}
