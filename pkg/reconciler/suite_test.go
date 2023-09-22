package reconciler_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestApplicationControllerManagerStateFunctions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ACM StateFunction Suite")
}
