package marshal

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type String struct {
	val string
	fn  string
	opt string
}

func (s *String) UnmarshalYAML(data []byte) error {
	var str string
	if err := yaml.Unmarshal(data, &str); err != nil {
		return err
	}

	if after, found := strings.CutPrefix(str, "$"); found {
		fn, fnParam, found := strings.Cut(after, ":")
		if !found {
			s.val = str
			return nil
		}

		switch fn {
		case "env":
			s.val = os.Getenv(fnParam)

		case "file":
			data, err := os.ReadFile(fnParam)
			if err != nil {
				return err
			}
			s.val = string(data)

		default:
			s.val = str
			return nil
		}

		s.fn = fn
		s.opt = fnParam
		return nil
	}

	s.val = str
	return nil
}

func (s *String) MarshalYAML() ([]byte, error) {
	if s.fn == "" {
		return []byte(s.val), nil
	}

	value := fmt.Sprintf("$%s:%s", s.fn, s.opt)
	return []byte(value), nil
}

func (s *String) Value() string  { return s.val }
func (s *String) Func() string   { return s.fn }
func (s *String) Option() string { return s.opt }

var _ interface {
	yaml.BytesUnmarshaler
	yaml.BytesMarshaler
} = (*String)(nil)

func NewString(value string) String {
	return String{val: value}
}
