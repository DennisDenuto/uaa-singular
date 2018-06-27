package relying_party

import (
	"github.com/DennisDenuto/go-uaa-singular/oidc/session_management/iframes/oidc_party"
	"github.com/DennisDenuto/go-uaa-singular/oidc/session_management"
)

const sessionStateKey = "session_state"

type SessionPoller struct {
	OpenIDProviderClient oidc_party.Client
	ClientID             string
	SessionStateStore    session_management.Store
	SessionKeyPrefix     string
}

func (sp SessionPoller) CheckSession() {
	stateValue := sp.SessionStateStore.Get(session_management.StateKey(sp.SessionKeyPrefix + "-" + sessionStateKey))
	sp.OpenIDProviderClient.CheckSession(sp.ClientID, string(stateValue))
}
