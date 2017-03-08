package multipartbuilder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMultipartbuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Multipartbuilder Suite")
}
