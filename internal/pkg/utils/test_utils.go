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
	"path/filepath"
	"strings"
)

// RunIntegrationTests checks whether integration tests should be executed.
func RunIntegrationTests() bool {
	var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")
	return runIntegration == "true"
}

// ProjectDir is a method to extract the folder of the project.
func ProjectDir(project string) string {
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, project) {
		wd = filepath.Dir(wd)
	}
	return wd
}
