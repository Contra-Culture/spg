package spg_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSpg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Spg Suite")
}
