package queryParamsCtrl_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQueryParamsCtrl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Query Params Ctrl Suite")
}
