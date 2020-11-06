package config

import (
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Config tests", func() {

	cfgValid := Config{
		Commit:  "111",
		Version: "0.10",

		RDBMS: RDBMS{
			ConnString:        "connectionstring",
			PingRetries:       3,
			PingWaitingPeriod: 5 * time.Second,
			SkipPing:          false,
		},
	}

	ginkgo.It("Should be valid", func() {
		gomega.Expect(cfgValid.IsValid()).Should(gomega.Succeed())
	})

	cfgInvalidCommit := Config{
		Commit:  "",
		Version: "0.10",

		RDBMS: RDBMS{
			ConnString:        "connectionstring",
			PingRetries:       3,
			PingWaitingPeriod: 5 * time.Second,
			SkipPing:          false,
		},
	}

	ginkgo.It("Should be invalid", func() {
		gomega.Expect(cfgInvalidCommit.IsValid()).ShouldNot(gomega.Succeed())
	})

	cfgInvalidVersion := Config{
		Commit:  "111",
		Version: "",

		RDBMS: RDBMS{
			ConnString:        "connectionstring",
			PingRetries:       3,
			PingWaitingPeriod: 5 * time.Second,
			SkipPing:          false,
		},
	}

	ginkgo.It("Should be invalid", func() {
		gomega.Expect(cfgInvalidVersion.IsValid()).ShouldNot(gomega.Succeed())
	})

})
