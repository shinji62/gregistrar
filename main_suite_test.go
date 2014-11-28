package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var GregistarPath string

func TestGregistrar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gregistrar Suite")
}

var _ = BeforeSuite(func() {
	path, err := gexec.Build("github.com/shinji62/gregistrar", "-race")
	Ω(err).ShouldNot(HaveOccurred())
	GregistarPath = path
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
