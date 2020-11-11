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

package rdbms

import (
	"time"

	"github.com/napptive/rdbms/internal/pkg/config"
	"github.com/napptive/rdbms/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

var _ = ginkgo.Describe("Ping Test", func() {
	
	if !utils.RunIntegrationTests("ping_test") {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var connstring = "host=localhost user=postgres password=Pass2020! port=5432"
	var invconnstring = "host=localhost user=postgres password=Pass2020! port=5431"

	ginkgo.It("Should work", func() {

		cfg := config.Config{
			Commit:  "111",
			Version: "0.10",

			RDBMS: config.RDBMS{
				ConnString:        connstring,
				PingRetries:       3,
				PingWaitingPeriod: 30 * time.Second,
				SkipPing:          false,
			},
		}

		gomega.Expect(Ping(cfg)).To(gomega.Succeed())
	})

	ginkgo.It("Should not work", func() {

		cfg := config.Config{
			Commit:  "111",
			Version: "0.10",

			RDBMS: config.RDBMS{
				ConnString:        invconnstring,
				PingRetries:       2,
				PingWaitingPeriod: 1 * time.Second,
				SkipPing:          false,
			},
		}

		gomega.Expect(Ping(cfg)).NotTo(gomega.Succeed())
	})
})
