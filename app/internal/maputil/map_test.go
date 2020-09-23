package maputil

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gin test", func() {

	Context("happy path tests", func() {

		var source, destination map[string]string
		BeforeEach(func() {
			source = make(map[string]string)
			destination = make(map[string]string)
		})

		It("join maps should contain all elements", func() {
			source["a"] = "a"
			destination["b"] = "b"

			finalMap := JoinMaps(source, destination)

			Expect(len(finalMap)).To(BeEquivalentTo(2))
		})

		It("join maps should contain all elements without repetition", func() {
			source["a"] = "a"
			source["b"] = "b"
			destination["b"] = "b"

			finalMap := JoinMaps(source, destination)

			Expect(len(finalMap)).To(BeEquivalentTo(2))
		})
	})

	Context("not happy path tests", func() {

		It("join maps should not crash with source nil maputil", func() {
			destination := make(map[string]string)
			destination["a"] = "a"

			finalMap := JoinMaps(nil, destination)

			Expect(len(finalMap)).To(BeEquivalentTo(1))
		})

		It("join maps should not crash with destiny nil maputil", func() {
			source := make(map[string]string)
			source["a"] = "a"

			finalMap := JoinMaps(source, nil)

			Expect(len(finalMap)).To(BeEquivalentTo(1))
		})
	})

})
