package connector

import (
	"net/http"
	"net/url"
	"path"
)

type client struct {
	base *url.URL
	http *http.Client
}

func (c *client) Do(r *http.Request) (*http.Response, error) {
	return c.http.Do(r)
}

func (c *client) Resolve(endpoint *url.URL) (string, error) {
	base := *c.base

	if base.Path != "" {
		endpoint.Path = path.Join(base.Path, endpoint.Path)
		base.Path = ""
	}

	resolved := base.ResolveReference(endpoint)
	return resolved.String(), nil
}
