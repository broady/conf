// +build appengine

import "google.golang.org/appengine"

func init() {
	appengineProject = appengine.AppID
}
