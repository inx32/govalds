package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

func CallerFormatter(noColor bool, module string) zerolog.Formatter {
	return func(a any) string {
		if path, ok := a.(string); ok {
			if cwd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(cwd, path); err == nil {
					path = rel
				}
			}

			path = strings.TrimPrefix(path, module+"/")

			if noColor {
				path = "[" + path + "]"
			} else {
				path = "\033[2m[" + path + "]\033[0m"
			}

			return path
		}
		if a == nil {
			return ""
		}
		return fmt.Sprint(a)
	}
}

func LevelFormatter(noColor bool) zerolog.Formatter {
	return func(a any) string {
		if level, ok := a.(string); ok {
			parsedLevel, err := zerolog.ParseLevel(level)
			if err != nil {
				return "<level invalid>"
			}
			levelString := strings.ToUpper(parsedLevel.String())
			if !noColor {
				color := zerolog.LevelColors[parsedLevel]
				return fmt.Sprintf("\033[%dm%s\033[0m", color, levelString)
			}
			return levelString
		}
		if a == nil {
			return ""
		}
		return fmt.Sprint(a)
	}
}

func ValueFormatter(noColor bool) zerolog.FormatterByFieldName {
	return func(a any, s string) string {
		if s == "from" {
			if a == nil {
				return ""
			}

			if noColor {
				return fmt.Sprintf("[%s]", a)
			} else {
				return fmt.Sprintf("\033[2m[%s]\033[0m", a)
			}
		}
		return s
	}
}

func ApplyFormatters(w *zerolog.ConsoleWriter, goModule string) {
	w.FormatCaller = CallerFormatter(w.NoColor, goModule)
	w.FormatLevel = LevelFormatter(w.NoColor)
	w.FormatPartValueByName = ValueFormatter(w.NoColor)
}
