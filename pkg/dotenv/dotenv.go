package dotenv

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type SyntaxError struct{ Message string }

func (e SyntaxError) Error() string {
	return "syntax error: " + e.Message
}

func isCharAllowed(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_'
}

func ParseLine(line string) (string, string, bool, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", "", false, nil
	}
	if line[0] == '#' {
		return "", "", false, nil
	}

	line = strings.TrimPrefix(line, "export ")
	line = strings.TrimSpace(line)

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false, SyntaxError{"line does not contain '=' char"}
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return "", "", false, SyntaxError{"key is empty"}
	}
	if value == "" {
		return key, value, true, nil
	}

	var quoted bool

	if len(value) >= 2 && value[0] == '"' {
		last := value[len(value)-1]
		if last != '"' {
			return "", "", false, SyntaxError{"invalid value: unterminated quote"}
		}
		value = value[1 : len(value)-1]
		quoted = true
	}

	if key[0] >= '0' && key[0] <= '9' {
		return "", "", false, SyntaxError{"invalid key: first character must not be a digit"}
	}

	for i, c := range key {
		if !isCharAllowed(c) {
			return "", "", false, SyntaxError{
				fmt.Sprintf("invalid key: non-alphanumeric char %d on index %d", c, i),
			}
		}
	}

	if !quoted {
		return key, value, true, nil
	}

	var b strings.Builder
	var esc bool

	for _, c := range value {
		if c == '\\' && !esc {
			esc = true
			continue
		}

		if esc {
			switch c {
			case 'n':
				b.WriteRune('\n')
			case 't':
				b.WriteRune('\t')
			case 'r':
				b.WriteRune('\r')
			default:
				b.WriteRune(c)
			}

			esc = false
			continue
		}

		b.WriteRune(c)
	}

	return key, b.String(), true, nil
}

func Parse(r io.Reader) (map[string]string, error) {
	sc := bufio.NewScanner(r)
	m := make(map[string]string)

	for sc.Scan() {
		line := sc.Text()
		key, value, found, err := ParseLine(line)

		if err != nil {
			return nil, fmt.Errorf("parse line: %w", err)
		}
		if !found {
			continue
		}

		if _, ok := m[key]; ok {
			return nil, fmt.Errorf("key %s found twice", key)
		}
		m[key] = value
	}

	if err := sc.Err(); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	return m, nil
}

func Load(r io.Reader) error {
	m, err := Parse(r)
	if err != nil {
		return err
	}
	for k, v := range m {
		os.Setenv(k, v)
	}
	return nil
}

func LoadFiles(files []string) error {
	for _, file := range files {
		fd, err := os.OpenFile(file, os.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("open file %s: %w", file, err)
		}
		err = Load(fd)
		fd.Close()

		if err != nil {
			return fmt.Errorf("parse file %s: %w", file, err)
		}
	}
	return nil
}
