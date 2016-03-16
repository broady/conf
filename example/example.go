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
}
