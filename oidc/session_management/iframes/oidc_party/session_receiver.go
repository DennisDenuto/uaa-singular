package oidc_party

import (
	"github.com/DennisDenuto/go-uaa-singular/oidc/session_management"
	"github.com/DennisDenuto/go-uaa-singular/oidc"
	"honnef.co/go/js/dom"
)

type CheckSessionHandlerFunc func(changed bool)

//go:generate counterfeiter . Client
type Client interface {
	CheckSession(clientId string, sessionState string) (session_management.UserAuthenticated, error)
}

type Browser interface {
	Location(string) error
}

type OpenIDProviderClient struct {
	SilentAuthURL string
	RedirectURL   string
	AuthResponse  chan session_management.UserAuthenticated
	Browser       Browser
}

func (c OpenIDProviderClient) CheckSession(clientId string, sessionState string) (session_management.UserAuthenticated, error) {
	authURL, err := oidc.NewSilentAuthURL(c.SilentAuthURL, c.RedirectURL, clientId, sessionState)
	if err != nil {
		return session_management.UserAuthenticated{}, err
	}

	//c.Browser.Location(authURL)
	dom.GetWindow().Open(authURL, "name", "features")
	//http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
	//	println(req.URL.String())
	//	return http.ErrUseLastResponse
	//}
	//_, err = http.Get(authURL)
	//if err != nil {
	//	panic(err)
	//}

	var userAuth = <-c.AuthResponse
	return userAuth, nil
}

type ServerReceiever interface {
	Receieve(clientId, sessionState, salt string) error
}

type OIDCServerReceiever struct {
}

func (OIDCServerReceiever) Receieve(clientId, sessionState, salt string) error {
	panic("implement me")
}
