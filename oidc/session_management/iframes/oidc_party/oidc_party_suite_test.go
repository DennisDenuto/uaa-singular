package oidc_party_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOidcParty(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OidcParty Suite")
}
