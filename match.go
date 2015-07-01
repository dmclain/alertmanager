// Copyright 2015 Prometheus Team
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

package main

import (
	"regexp"

	"github.com/prometheus/common/model"
)

// Matcher defines a matching rule for the value of a given label.
type Matcher struct {
	Name  model.LabelName
	Value string

	isRegex bool
	regex   *regexp.Regexp
}

// IsRegex returns true of the matcher compares against a regular expression.
func (m *Matcher) IsRegex() bool {
	return m.isRegex
}

// Match checks whether the label of the matcher has the specified
// matching value.
func (m *Matcher) Match(lset model.LabelSet) bool {
	// Unset labels are treated as unset labels globally. Thus, if a
	// label is not set we retrieve the empty label which is correct
	// for the comparison below.
	v := lset[m.Name]

	if m.isRegex {
		return m.regex.MatchString(string(v))
	}
	return string(v) == m.Value
}

// NewMatcher returns a new matcher that compares against equality of
// the given value.
func NewMatcher(name model.LabelName, value string) *Matcher {
	return &Matcher{
		Name:  name,
		Value: value,
	}
}

// NewRegexMatcher returns a new matcher that treats value as a regular
// expression which is used for matching.
func NewRegexMatcher(name model.LabelName, value string) (*Matcher, error) {
	re, err := regexp.Compile(value)
	if err != nil {
		return nil, err
	}
	m := &Matcher{
		Name:    name,
		Value:   value,
		isRegex: true,
		regex:   re,
	}
	return m, nil
}

type Matchers []*Matcher

// MatchAll checks whether all matchers are fulfilled against the given label set.
func (ms Matchers) MatchAll(lset model.LabelSet) bool {
	for _, m := range ms {
		if !m.Match(lset) {
			return false
		}
	}
	return true
}
