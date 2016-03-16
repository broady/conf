package main

import (
	"log"

	"golang.org/x/net/context"

	"github.com/broady/conf"
	"github.com/broady/conf/gcloudconf"
)

func main() {
	ctx := context.Background()

	foo := conf.MustGet(
		conf.Env("FOO"),
		gcloudconf.Metadata(ctx, "foo"),
		conf.Default("bar"))

	log.Print(foo)

	xx, err := conf.Get(
		conf.Env("FOO1"),
		conf.Env("FOO2"))

	// must set one of: environment variable FOO1, environment variable FOO2
	log.Print(xx, err)
}
