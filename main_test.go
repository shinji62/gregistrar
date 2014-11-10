package main_test

import (
	"github.com/cloudfoundry/gorouter/test_util"
	"github.com/cloudfoundry/gunk/natsrunner"
	"github.com/fraenkel/candiedyaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"github.com/shinji62/gregistrar/config"
	"github.com/shinji62/gregistrar/mbus"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var _ = Describe("Nats Registration", func() {

	var tmpdir string
	var natsPort uint16
	var natsRunner *natsrunner.NATSRunner
	var gregistrarSession *Session

	BeforeEach(func() {
		var err error
		tmpdir, err = ioutil.TempDir("", "gregistrar")
		Ω(err).ShouldNot(HaveOccurred())

		natsPort = test_util.NextAvailPort()
		natsRunner = natsrunner.NewNATSRunner(int(natsPort))
		natsRunner.Start()
	})

	createConfig := func(cfgFile string) *config.Config {

		c := config.DefaultConfig()

		c.Nats = []config.NatsConfig{
			config.NatsConfig{
				Host: "localhost",
				Port: natsPort,
				User: "",
				Pass: "",
			},
		}

		c.Uris = []config.UriConfig{
			config.UriConfig{
				Uri: "local.local",
			},
		}

		c.LoggingLevel = "info"
		c.RegisterInterval = 5
		c.IpRegexp = "^127\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$"
		cfgBytes, err := candiedyaml.Marshal(c)
		Ω(err).ShouldNot(HaveOccurred())

		ioutil.WriteFile(cfgFile, cfgBytes, os.ModePerm)
		return c
	}

	startGregistrar := func(cfgFile string) *Session {
		gregistrarCmd := exec.Command(GregistarPath, "-c", cfgFile)
		session, err := Start(gregistrarCmd, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session, 5).Should(Say("gorouter.started"))
		gregistrarSession = session

		return session
	}

	stopGregistrar := func(gregistrarSession *Session) {
		err := gregistrarSession.Command.Process.Signal(syscall.SIGTERM)
		Ω(err).ShouldNot(HaveOccurred())
		Expect(gregistrarSession.Wait(5 * time.Second)).Should(Exit(0))
	}

	AfterEach(func() {
		if natsRunner != nil {
			natsRunner.Stop()
		}

		os.RemoveAll(tmpdir)

		if gregistrarSession != nil {
			stopGregistrar(gregistrarSession)
		}
	})

	Describe("NatsMessageTest", func() {

		c := createConfig("test.yml")
		gregistrar := startGregistrar("test.yml")

		natsClient, err := mbus.NewMessageBusConnection(c)

	})

})
