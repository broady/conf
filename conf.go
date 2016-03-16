package conf

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func MustGet(g ...Getter) string {
	v, err := Get(g...)
	if err != nil {
		panic(err)
	}
	return v
}

type Getter interface {
	Get() (string, error)
	Usage() string
}

func Get(g ...Getter) (string, error) {
	var usages []string
	for _, gg := range g {
		v, err := gg.Get()
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

func Env(name string) Getter {
	return envGetter(name)
}

type envGetter string

func (e envGetter) Get() (string, error) {
	v := os.Getenv(string(e))
	if v == "" {
		return "", Missing
	}
	return v, nil
}

func (e envGetter) Usage() string {
	return fmt.Sprintf("environment variable %s", e)
}

var Missing error = errors.New("missing")

type EvalError struct {
	g   Getter
	err error
}

func (e EvalError) Error() string {
	return fmt.Sprintf("%s: %v", e.g.Usage(), e.err)
}

type defaulter string

func (d defaulter) Get() (string, error) {
	return string(d), nil
}

func (d defaulter) Usage() string {
	return fmt.Sprintf("default value %s", d)
}

func Default(v string) Getter {
	return defaulter(v)
}
