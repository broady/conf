package conf

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// MustGet is the same as Get, but panics if an error is returned.
func MustGet(s ...Source) string {
	v, err := Get(s...)
	if err != nil {
		panic(err)
	}
	return v
}

// Source represents a source of a configuration value.
type Source interface {
	// Evaluate returns the value from the source.
	// The "Missing" error is returned if the value is missing.
	// Any other error is returned if something unexpected occurred while evaluating the value.
	Evaluate() (string, error)

	// A short string describing the source of this value.
	// For example, Env("FOO") has a usage string of "environment variable FOO".
	// This string is used to trace back the source of an error or the source of a value.
	Usage() string
}

// Get resolves a list of possible sources for the value, trying the next one if the value is missing.
//
// For example, in the following code, unless the "FLAG" environment variable is set, val will be set to "false".
//	val, _ := conf.Get(conf.Env("FLAG"), conf.Default("false"))
func Get(s ...Source) (string, error) {
	var usages []string
	for _, source := range s {
		v, err := source.Evaluate()
		if err == Missing {
			usages = append(usages, source.Usage())
			continue
		}
		if err != nil {
			return "", EvalError{source, err}
		}
		return v, nil
	}
	if len(usages) == 1 {
		return "", fmt.Errorf("must set %s", usages[0])
	}
	return "", fmt.Errorf("must set one of: %s", strings.Join(usages, ", "))
}

// Env gets the value from the environment variable.
func Env(name string) Source {
	return envSource(name)
}

type envSource string

func (e envSource) Evaluate() (string, error) {
	v := os.Getenv(string(e))
	if v == "" {
		return "", Missing
	}
	return v, nil
}

func (e envSource) Usage() string {
	return fmt.Sprintf("environment variable %s", e)
}

// Missing is returned by a source when it can't resolve the value.
var Missing = errors.New("missing")

// EvalError is an error that occurs when evaluating a Source.
type EvalError struct {
	s   Source
	err error
}

func (e EvalError) Error() string {
	return fmt.Sprintf("%s: %v", e.s.Usage(), e.err)
}

type defaulter string

func (d defaulter) Evaluate() (string, error) {
	return string(d), nil
}

func (d defaulter) Usage() string {
	return fmt.Sprintf("default value %s", d)
}

// Default always returns the given value.
func Default(v string) Source {
	return defaulter(v)
}
