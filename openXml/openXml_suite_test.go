package openXml_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOpenXml(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OpenXml Suite")
}
