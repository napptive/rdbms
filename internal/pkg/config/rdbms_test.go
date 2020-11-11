/**
 * Copyright 2020 Napptive
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Config RDBMS tests", func() {

	cfgValid := RDBMS{
		ConnString:        "connectionstring",
		PingRetries:       3,
		PingWaitingPeriod: 5 * time.Second,
		SkipPing:          false,
	}

	ginkgo.It("Should be valid", func() {
		gomega.Expect(cfgValid.IsValid()).To(gomega.Succeed())
	})

	cfgInvalidConnString := RDBMS{
		ConnString:        "",
		PingRetries:       3,
		PingWaitingPeriod: 5 * time.Second,
		SkipPing:          false,
	}

	ginkgo.It("Should be invalid", func() {
		gomega.Expect(cfgInvalidConnString.IsValid()).ToNot(gomega.Succeed())
	})

	cfgInvalidrRetries := RDBMS{
		ConnString:        "connectionstring",
		PingRetries:       -1,
		PingWaitingPeriod: 5 * time.Second,
		SkipPing:          false,
	}

	ginkgo.It("Should be invalid", func() {
		gomega.Expect(cfgInvalidrRetries.IsValid()).ToNot(gomega.Succeed())
	})
})
