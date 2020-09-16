package awslambda_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAwslambda(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Awslambda Suite")
}
