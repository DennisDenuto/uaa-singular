package relying_party_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	rp "github.com/DennisDenuto/go-uaa-singular/oidc/session_management/iframes/relying_party"
	"github.com/DennisDenuto/go-uaa-singular/oidc/session_management/iframes/oidc_party/oidc_partyfakes"
	"github.com/DennisDenuto/go-uaa-singular/oidc/session_management/session_managementfakes"
	"github.com/DennisDenuto/go-uaa-singular/oidc/session_management"
)

var _ = Describe("SessionPoller", func() {
	var poller rp.SessionPoller
	var clientId string
	var fakeOpenIDProviderClient *oidc_partyfakes.FakeClient
	var fakeStore *session_managementfakes.FakeStore

	BeforeEach(func() {
		fakeOpenIDProviderClient = &oidc_partyfakes.FakeClient{}

		clientId = "clientId"
		fakeStore = &session_managementfakes.FakeStore{}
		poller = rp.SessionPoller{
			SessionKeyPrefix:     "storagekey",
			ClientID:             clientId,
			OpenIDProviderClient: fakeOpenIDProviderClient,
			SessionStateStore:    fakeStore,
		}
	})

	Context("session state has been stored", func() {
		BeforeEach(func() {
			fakeStore.GetReturns(session_management.StateValue("some-session-state"))
		})
		It("Should post a message to an OP iframe", func() {
			poller.CheckSession()

			Expect(fakeStore.GetCallCount()).To(Equal(1))
			Expect(fakeStore.GetArgsForCall(0)).To(Equal(session_management.StateKey("storagekey-session_state")))

			clientId, sessionState := fakeOpenIDProviderClient.CheckSessionArgsForCall(0)
			Expect(clientId).To(Equal("clientId"))
			Expect(sessionState).To(Equal("some-session-state"))
		})
	})
})
