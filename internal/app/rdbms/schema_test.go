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

var _ = ginkgo.Describe("Load Schema test", func() {

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var basepath = utils.ProjectDir("napptive/rdbms")

	var connstring = "host=localhost user=postgres password=Pass2020! port=5432"

	var defaultDuration = 5 * time.Second

	ginkgo.It("Should work", func() {
		var filepath = basepath + "/test/data/ValidSQLScript.yaml"
		cfg := config.Config{
			Commit:  "111",
			Version: "0.10",

			RDBMS: config.RDBMS{
				ConnString:        connstring,
				PingRetries:       3,
				PingWaitingPeriod: 5 * time.Second,
				SkipPing:          false,
			},
		}
		result, err := Load(filepath, defaultDuration, []string{}, cfg)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result.ExecutedSteps).To(gomega.HaveLen(3))
		gomega.Expect(result.SkippedSteps).To(gomega.HaveLen(0))
		result.Print()
	})

	ginkgo.It("Should work", func() {
		var filepath = basepath + "/test/data/ValidSQLScript.yaml"
		cfg := config.Config{
			Commit:  "111",
			Version: "0.10",

			RDBMS: config.RDBMS{
				ConnString:        connstring,
				PingRetries:       3,
				PingWaitingPeriod: 5 * time.Second,
				SkipPing:          false,
			},
		}
		result, err := Load(filepath, defaultDuration, []string{"creation-step", "drop-step"}, cfg)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result.ExecutedSteps).To(gomega.HaveLen(2))
		gomega.Expect(result.SkippedSteps).To(gomega.HaveLen(1))
		result.Print()
	})

})
