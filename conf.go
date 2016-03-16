package conf

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func MustGet(g ...Value) string {
	v, err := Get(g...)
	if err != nil {
		panic(err)
	}
	return v
}

type Value interface {
	Value() (string, error)
	Usage() string
}

func Get(g ...Value) (string, error) {
	var usages []string
	for _, gg := range g {
		v, err := gg.Value()
		if err == Missing {
			usages = append(usages, gg.Usage())
			continue
		}
		if err != nil {
			return "", EvalError{gg, err}
		}
		return v, nil
	}
	if len(usages) == 1 {
		return "", fmt.Errorf("must set %s", usages[0])
	}
	return "", fmt.Errorf("must set one of: %s", strings.Join(usages, ", "))
}

func Env(name string) Value {
	return envValue(name)
}

type envValue string

func (e envValue) Value() (string, error) {
	v := os.Getenv(string(e))
	if v == "" {
		return "", Missing
	}
	return v, nil
}

func (e envValue) Usage() string {
	return fmt.Sprintf("environment variable %s", e)
}

var Missing error = errors.New("missing")

type EvalError struct {
	g   Value
	err error
}

func (e EvalError) Error() string {
	return fmt.Sprintf("%s: %v", e.g.Usage(), e.err)
}

type defaulter string

func (d defaulter) Value() (string, error) {
	return string(d), nil
}

func (d defaulter) Usage() string {
	return fmt.Sprintf("default value %s", d)
}

func Default(v string) Value {
	return defaulter(v)
}
