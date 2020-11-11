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
	"errors"
	"os"
	"path/filepath"
)

var projectFiles = []string{"Makefile"}

// RunIntegrationTests checks whether integration tests should be executed.
func RunIntegrationTests(id string) bool {
	var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")
	if runIntegration == "all" {
		return true
	}
	return runIntegration == id
}

// ProjectDir is a method to extract the folder of the project.
func ProjectDir() (*string, error) {

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return findDirWithFiles(wd, projectFiles)
}

func findDirWithFiles(wd string, files []string) (*string, error) {

	path := wd
	prevPath := ""
	for !checkFiles(path, files) && path != prevPath {
		prevPath = path
		path = filepath.Dir(path)
	}
	if path == prevPath {
		return nil, errors.New("Not found directory")
	}
	return &path, nil
}

func checkFiles(dir string, files []string) bool {
	for _, file := range files {
		path := filepath.Join(dir, file)
		if !checkFile(path) {
			return false
		}
	}
	return true
}

func checkFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
