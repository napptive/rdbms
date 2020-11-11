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

package utils

import (
	"os"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Project Dir tests", func() {

	ginkgo.It("Should find the project dir", func() {
		path, err := ProjectDir()
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(*path).To(gomega.BeAnExistingFile())
	})
})

var _ = ginkgo.Describe("Find Dir With Files tests", func() {

	ginkgo.It("Should find the  dir", func() {
		wd, err := os.Getwd()
		gomega.Expect(err).To(gomega.Succeed())

		path, err := findDirWithFiles(wd, []string{"Makefile", "README.md"})
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(*path).To(gomega.BeAnExistingFile())
	})

	ginkgo.It("Should not find wrong files", func() {
		wd, err := os.Getwd()
		gomega.Expect(err).To(gomega.Succeed())

		path, err := findDirWithFiles(wd, []string{"Makefi", "README.md"})
		gomega.Expect(err).NotTo(gomega.Succeed())
		gomega.Expect(path).To(gomega.BeNil())
	})

	ginkgo.It("Should not find local dir", func() {
		wd := "."

		path, err := findDirWithFiles(wd, []string{"Makefi", "README.md"})
		gomega.Expect(err).NotTo(gomega.Succeed())
		gomega.Expect(path).To(gomega.BeNil())
	})
	ginkgo.It("Should not find random dir", func() {
		wd := "/folder/_rr/croq"

		path, err := findDirWithFiles(wd, []string{"Makefi", "README.md"})
		gomega.Expect(err).NotTo(gomega.Succeed())
		gomega.Expect(path).To(gomega.BeNil())
	})
	ginkgo.It("Should not find windows directory", func() {
		wd := "c:\\folder\\_rr\\croq"

		path, err := findDirWithFiles(wd, []string{"Makefi", "README.md"})
		gomega.Expect(err).NotTo(gomega.Succeed())
		gomega.Expect(path).To(gomega.BeNil())
	})
})
