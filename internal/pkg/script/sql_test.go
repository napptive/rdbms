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

package script

import (
	"time"

	"github.com/napptive/rdbms/internal/pkg/utils"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("SQL Timeout Duration tests", func() {

	ginkgo.It("Should be default value", func() {
		var step = SQLStep{
			Name:    "test",
			Timeout: "",
			Queries: nil,
		}
		var defaultDuration = 5 * time.Second
		dur, err := step.TimeoutDuration(defaultDuration)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(dur).To(gomega.Equal(defaultDuration))
	})
	ginkgo.It("Should be specific value", func() {
		var step = SQLStep{
			Name:    "test",
			Timeout: "17s",
			Queries: nil,
		}
		var defaultDuration = 5 * time.Second
		dur, err := step.TimeoutDuration(defaultDuration)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(dur).To(gomega.Equal(17 * time.Second))
	})

	ginkgo.It("Should fail", func() {
		var step = SQLStep{
			Name:    "test",
			Timeout: "q7s",
			Queries: nil,
		}
		var defaultDuration = 5 * time.Second
		_, err := step.TimeoutDuration(defaultDuration)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

})

var _ = ginkgo.Describe("SQL File Parse tests", func() {

	var basepath = utils.ProjectDir("napptive/rdbms")

	ginkgo.It("Should parse the file", func() {
		var filepath = basepath + "/test/data/ValidSQLScript.yaml"
		sql, err := SQLFileParse(filepath)

		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(sql).NotTo(gomega.BeNil())
		gomega.Expect(sql.Steps).To(gomega.HaveLen(3))
	})
	ginkgo.It("Should not parse anything", func() {
		var filepath = basepath + "/test/data/InvalidSQLScript.yaml"
		sql, err := SQLFileParse(filepath)

		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(sql).NotTo(gomega.BeNil())
		gomega.Expect(sql.Steps).To(gomega.BeEmpty())
	})
	ginkgo.It("Should fail incorrect format", func() {
		var filepath = basepath + "/test/data/RandomFile.txt"
		_, err := SQLFileParse(filepath)

		gomega.Expect(err).NotTo(gomega.Succeed())
	})
	ginkgo.It("Should fail not found file", func() {
		var filepath = "test/data/NotFound.yaml"
		_, err := SQLFileParse(filepath)

		gomega.Expect(err).NotTo(gomega.Succeed())
	})

})
