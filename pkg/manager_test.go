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
		err := mgr.Go("1", func(stopCh <-chan struct{}) {
			for {
				select {
				case <-stopCh:
					fmt.Println("stopping")
					return
				default:
					fmt.Println("running")
					time.Sleep(time.Second)
				}
			}
		})
		Expect(err).ToNot(HaveOccurred())
		time.Sleep(3 * time.Second)
		close(stopCh)
		time.Sleep(3 * time.Second)
	})
})
