package pkg

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go routine manager", func() {
	It("Run go routine", func() {
		mgr := NewGoRoutineManager(1)
		err := mgr.Go("1", func() {
			fmt.Println("test")
		})
		Expect(err).ToNot(HaveOccurred())
	})
})
