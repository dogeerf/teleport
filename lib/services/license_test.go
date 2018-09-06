/*
Copyright 2018 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package services

import (
	"testing"

	"github.com/gravitational/teleport/lib/fixtures"
	"github.com/gravitational/teleport/lib/utils"

	"github.com/gravitational/trace"
	"gopkg.in/check.v1"
)

func TestLicenseInfo(t *testing.T) { check.TestingT(t) }

type LicenseInfoSuite struct {
}

var _ = check.Suite(&LicenseInfoSuite{})

func (s *LicenseInfoSuite) SetUpSuite(c *check.C) {
	utils.InitLoggerForTests()
}

func (s *LicenseInfoSuite) TestUnmarshal(c *check.C) {
	type testCase struct {
		description string
		input       string
		expected    LicenseInfo
		err         error
	}
	testCases := []testCase{
		{
			description: "simple case",
			input:       `{"kind": "license_info", "version": "v3", "metadata": {"name": "Teleport Commercial"}, "spec": {"account_id": "accountID", "plan_id": "planID", "usage": true, "k8s": true, "aws_account": "123", "aws_pid": "4"}}`,
			expected: MustNew("Teleport Commercial", LicenseInfoSpecV3{
				ReportsUsage:       NewBool(true),
				SupportsKubernetes: NewBool(true),
				AWSAccountID:       "123",
				AWSProductID:       "4",
				AccountID:          "accountID",
				PlanID:             "planID",
			}),
		},
		{
			description: "simple case with string booleans",
			input:       `{"kind": "license_info", "version": "v3", "metadata": {"name": "license_info"}, "spec": {"account_id": "accountID", "plan_id": "planID", "usage": "yes", "k8s": "yes", "aws_account": "123", "aws_pid": "4"}}`,
			expected: MustNew("license_info", LicenseInfoSpecV3{
				ReportsUsage:       NewBool(true),
				SupportsKubernetes: NewBool(true),
				AWSAccountID:       "123",
				AWSProductID:       "4",
				AccountID:          "accountID",
				PlanID:             "planID",
			}),
		},
		{
			description: "failed validation - unknown version",
			input:       `{"kind": "license_info", "version": "v2", "metadata": {"name": "license_info"}, "spec": {"usage": "yes", "k8s": "yes", "aws_account": "123", "aws_pid": "4"}}`,
			err:         trace.BadParameter(""),
		},
		{
			description: "failed validation, bad types",
			input:       `{"kind": "license_info", "version": "v3", "metadata": {"name": "license_info"}, "spec": {"usage": 1, "k8s": "yes", "aws_account": 14, "aws_pid": "4"}}`,
			err:         trace.BadParameter(""),
		},
	}
	for _, tc := range testCases {
		comment := check.Commentf("test case %q", tc.description)
		out, err := UnmarshalLicenseInfo([]byte(tc.input))
		if tc.err == nil {
			c.Assert(err, check.IsNil, comment)
			fixtures.DeepCompare(c, tc.expected, out)
			data, err := MarshalLicenseInfo(out)
			c.Assert(err, check.IsNil, comment)
			out2, err := UnmarshalLicenseInfo(data)
			c.Assert(err, check.IsNil, comment)
			fixtures.DeepCompare(c, tc.expected, out2)
		} else {
			c.Assert(err, check.FitsTypeOf, tc.err, comment)
		}
	}
}

// MustNew is like New, but panics in case of error,
// used in tests
func MustNew(name string, spec LicenseInfoSpecV3) LicenseInfo {
	out, err := NewLicenseInfo(name, spec)
	if err != nil {
		panic(err)
	}
	return out
}
