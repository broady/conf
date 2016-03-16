package gcloudconf

import (
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/cloud/compute/metadata"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"github.com/broady/conf"
)

// Metadata retrieves the value from the project metadata, stored at the given key.
func Metadata(ctx context.Context, key string) conf.Source {
	return &metadataSource{ctx, key}
}

type metadataSource struct {
	ctx context.Context
	key string
}

var appengineProject func(ctx context.Context) string

func (e *metadataSource) Evaluate() (string, error) {
	if metadata.OnGCE() {
		return metadata.ProjectAttributeValue(e.key)
	}
	if appengineProject != nil {
		proj := appengineProject(e.ctx)
		hc, err := google.DefaultClient(e.ctx, compute.ComputeScope)
		if err != nil {
			return "", err
		}
		svc, err := compute.New(hc)
		if err != nil {
			return "", err
		}
		p, err := svc.Projects.Get(proj).Context(e.ctx).Do()
		if err != nil {
			return "", err
		}
		for _, item := range p.CommonInstanceMetadata.Items {
			if item.Key == e.key {
				return *item.Value, nil
			}
		}
		return "", conf.Missing
	}
	// TODO(cbro): figure out a way to get the project ID locally?
	//return "", errors.New("not running on Cloud")
	return "", conf.Missing
}

func (e *metadataSource) Usage() string {
	return fmt.Sprintf("Google Cloud project metadata variable %s", e.key)
}
