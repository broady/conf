package conf

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestGetEnv(t *testing.T) {
	const notExist = "DOES_NOT_EXIST"
	_, err := Get(Env(notExist))

	if !strings.Contains(err.Error(), notExist) {
		t.Errorf("expected %q to contain %s", err, notExist)
	}

	if !strings.Contains(err.Error(), "environment variable") {
		t.Errorf("expected %q to contain %s", err, "environment variable")
	}
}

func TestGetEnvMultiple(t *testing.T) {
	const notExist = "DOES_NOT_EXIST"
	_, err := Get(Env(notExist+"0"), Env(notExist+"1"), Env(notExist+"2"))

	for i := 0; i < 3; i++ {
		want := fmt.Sprintf("%s%d", notExist, i)
		if !strings.Contains(err.Error(), want) {
			t.Errorf("expected %q to contain %s", err, want)
		}
	}
}

func TestError(t *testing.T) {
	want := errors.New("A")

	_, err := Get(alwaysError{want})
	if err == nil {
		t.Error("want err, got nil")
	}
	if _, ok := err.(EvalError); !ok {
		t.Errorf("want EvalError, got %#v", err)
	}
}

func TestMissing(t *testing.T) {
	want := "foo"

	got, err := Get(alwaysError{Missing}, Default(want))
	if err != nil {
		t.Error("want nil err, got %#v", err)
	}
	if want != got {
		t.Error("want %s, got %s", want, got)
	}
}

type alwaysError struct {
	e error
}

func (a alwaysError) Get() (string, error) {
	return "", a.e
}

func (a alwaysError) Usage() string {
	return fmt.Sprintf("alwaysError %v", a.e.Error())
}
