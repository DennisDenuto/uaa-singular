package oidc_party_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/DennisDenuto/go-uaa-singular/oidc/session_management/iframes/oidc_party"
	"github.com/onsi/gomega/ghttp"
	"github.com/DennisDenuto/go-uaa-singular/oidc/session_management"
	"net/http"
)

var _ = Describe("SessionReceiver", func() {
	Describe("OpenIdProviderClient", func() {
		var opClient oidc_party.OpenIDProviderClient
		var UAAServer *ghttp.Server
		var userAuthChannel chan session_management.UserAuthenticated
		var userAuthenticatedResponse session_management.UserAuthenticated

		BeforeEach(func() {
			UAAServer = ghttp.NewServer()
			UAAServer.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/oauth/authorize", "client_id=admin&response_type=id_token&scope=openid&state=some-session-state&prompt=none&redirect_uri=http://some-redirecturl.com"),
				ghttp.RespondWith(302, ``, http.Header{"Location": {"http://some-redirecturl.com#error=login_required&session_state=15e5d9e438b14888d4ec85eec2518e3e0705714a7add8ab7bb994e5ba72da2bc.e665fbc86636dd72c4f50ed315259ca969bfd57183e160ea6105866ea6cb99b3"}}),
			))

			userAuthChannel = make(chan session_management.UserAuthenticated, 1)
			userAuthChannel <- userAuthenticatedResponse

			opClient = oidc_party.OpenIDProviderClient{
				SilentAuthURL: UAAServer.URL(),
				RedirectURL:   "http://some-redirecturl.com",
				AuthResponse:  userAuthChannel,
			}
		})

		It("should be able to post a message to OP", func() {
			userAuthenticated, err := opClient.CheckSession("admin", "some-session-state")
			Expect(err).NotTo(HaveOccurred())

			Expect(UAAServer.ReceivedRequests()).To(HaveLen(1))
			Expect(userAuthenticated).To(Equal(userAuthenticatedResponse))
		})

	})
})
