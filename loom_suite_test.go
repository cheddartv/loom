package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLoom(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Loom Suite")
}
