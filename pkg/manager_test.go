package pkg

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go routine manager", func() {
	It("Run go routine", func() {
		stopCh := make(chan struct{})
		mgr := NewGoRoutineManager(context.Background(), 1, stopCh)
		err := mgr.Go("1", func() {
			time.Sleep(time.Second)
			fmt.Println("test")
		})
		Expect(err).ToNot(HaveOccurred())
		time.Sleep(time.Second)
		close(stopCh)
	})
})
