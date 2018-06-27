package oidc

import (
	"github.com/cloudfoundry-community/go-uaa"
	"golang.org/x/oauth2"
	"net/url"
)

func NewSilentAuthURL(target, redirectURL, clientID, state string) (string, error) {
	url, err := uaa.BuildTargetURL(target)
	if err != nil {
		return "", err
	}

	authorizationURL, err := urlWithPath(*url, "oauth/authorize")
	if err != nil {
		return "", err
	}

	authorizationURL.Query().Add("scope", "openid")

	c := &oauth2.Config{
		ClientID:     clientID,
		Endpoint: oauth2.Endpoint{
			AuthURL: authorizationURL.String(),
		},
		Scopes:      []string{"openid"},
		RedirectURL: redirectURL,

	}

	codeURL := c.AuthCodeURL(state, oauth2.SetAuthURLParam("prompt", "none"), oauth2.SetAuthURLParam("response_type", "id_token"))
	return codeURL, nil
}



// urlWithPath copies the URL and sets the path on the copy.
func urlWithPath(u url.URL, path string) (*url.URL, error) {
	urlPath, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	return u.ResolveReference(urlPath), nil
}