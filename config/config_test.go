package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/shinji62/gregistrar/config"
)

var _ = Describe("Config", func() {
	var config *Config
	BeforeEach(func() {
		config = DefaultConfig()
	})

	Describe("Initialize", func() {

		It("sets nats config", func() {
			var b = []byte(`
nats:
  - host: remotehost
    port: 4223
    user: user
    pass: pass
`)
			config.Initialize(b)

			Ω(config.Nats).To(HaveLen(1))
			Ω(config.Nats[0].Host).To(Equal("remotehost"))
			Ω(config.Nats[0].Port).To(Equal(uint16(4223)))
			Ω(config.Nats[0].User).To(Equal("user"))
			Ω(config.Nats[0].Pass).To(Equal("pass"))
		})

		It("sets urls config", func() {
			var b = []byte(`
uris:
  - uri: testurl
 `)

			config.Initialize(b)
			Ω(config.Uris).To(HaveLen(1))
			Ω(config.Uris[0].Uri).To(Equal("testurl"))
		})

		It("sets register time interval config", func() {
			var b = []byte(`
register_message_interval: 5
 `)

			config.Initialize(b)
			Ω(config.RegisterInterval).To(Equal(5))
		})

		It("sets Maximun Go procs", func() {
			var b = []byte(`
max_go_procs: 5
 `)

			config.Initialize(b)
			Ω(config.GoMaxProcs).To(Equal(5))
		})

		It("set logging level", func() {
			var b = []byte(`
logging_level: debug
`)

			config.Initialize(b)
			Ω(config.LoggingLevel).To(Equal("debug"))

		})

		It("set port", func() {
			var b = []byte(`
port: 8080
`)

			config.Initialize(b)
			Ω(config.Port).To(Equal(uint16(8080)))

		})

	})
})
