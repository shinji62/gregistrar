package main_test

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry/gunk/natsrunner"
	"github.com/cloudfoundry/yagnats"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/shinji62/gregistrar/config"
	"net/url"
	"os"
)

var _ = Describe("Nats Registration", func() {

	var tmpdir string
	var natsPort uint16
	var natsRunner *natsrunner.NATSRunner

	BeforeEach(func() {
		var err error
		tmpdir, err = ioutil.TempDir("", "gorouter")
		Ω(err).ShouldNot(HaveOccurred())

		natsPort = test_util.NextAvailPort()
		natsRunner = natsrunner.NewNATSRunner(int(natsPort))
		natsRunner.Start()
	})

	createConfig := func(cfgFile string) *config.Config {

		config.NATS = []NatsConfig{
			Host: "localhost",
			Port: natsPort,
			User: "",
			Pass: "",
		}

		config.URL = []UrlConfig{
			Url: "local.local",
		}

		config.LoggingLevel = "info"
		config.RegisterInterval = 5
		cfgBytes, err := candiedyaml.Marshal(config)
		Ω(err).ShouldNot(HaveOccurred())

		ioutil.WriteFile(cfgFile, cfgBytes, os.ModePerm)
		return config
	}

})
