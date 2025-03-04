// Copyright 2023 Upbound Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package templates contains utilities for rendering kuttl test cases using
// the templates contained in the package.
package templates

import (
	"strings"
	"text/template"

	"github.com/crossplane/crossplane-runtime/pkg/errors"

	"github.com/upbound/uptest/internal/config"
)

var fileTemplates = map[string]string{
	"00-apply.yaml":  inputFileTemplate,
	"00-assert.yaml": assertFileTemplate,
	"01-delete.yaml": deleteFileTemplate,
	"01-assert.yaml": assertDeletedFileTemplate,
}

// Render renders the specified list of resources as a test case
// with the specified configuration.
func Render(tc *config.TestCase, resources []config.Resource) (map[string]string, error) {
	data := struct {
		Resources []config.Resource
		TestCase  config.TestCase
	}{
		Resources: resources,
		TestCase:  *tc,
	}

	res := make(map[string]string, len(fileTemplates))
	for name, tmpl := range fileTemplates {
		t, err := template.New(name).Parse(tmpl)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot parse template %q", name)
		}

		var b strings.Builder
		if err := t.Execute(&b, data); err != nil {
			return nil, errors.Wrapf(err, "cannot execute template %q", name)
		}
		res[name] = b.String()
	}

	return res, nil
}
