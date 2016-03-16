package main

import (
	"log"
	"net/http"

	"google.golang.org/api/plus/v1"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"github.com/broady/conf"
	"github.com/broady/conf/gcloudconf"
)

func main() {
	ctx := context.Background()

	oauth2Config := oauth2.Config{
		ClientID:     conf.MustGet(conf.Env("OAUTH2_CLIENT"), gcloudconf.Metadata(ctx, "oauth2-client")),
		ClientSecret: conf.MustGet(conf.Env("OAUTH2_SECRET"), gcloudconf.Metadata(ctx, "oauth2-secret")),
		RedirectURL: conf.MustGet(
			conf.Env("OAUTH2_REDIRECT"),
			gcloudconf.Metadata(ctx, "oauth2-redirect"),
			conf.Default("http://localhost:8080/oauth2callback")),
		Scopes: []string{plus.UserinfoEmailScope},
	}

	xx, err := conf.Get(
		conf.Env("FOO1"),
		conf.Env("FOO2"))

	// must set one of: environment variable FOO1, environment variable FOO2
	log.Print(xx, err)

	log.Print(oauth2Config.AuthCodeURL(""))

	log.Print("serving on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
