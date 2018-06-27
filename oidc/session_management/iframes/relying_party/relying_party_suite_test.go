package relying_party_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRelyingParty(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RelyingParty Suite")
}
