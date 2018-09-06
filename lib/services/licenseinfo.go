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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gravitational/teleport/lib/defaults"
	"github.com/gravitational/teleport/lib/utils"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
)

// LicenseInfo defines teleport License Information
type LicenseInfo interface {
	Resource
	// GetReportsUsage returns true if teleport cluster reports usage
	// to control plane
	GetReportsUsage() Bool

	// SetReportsUsage sets usage report
	SetReportsUsage(Bool)

	// GetAWSProductID returns product id that limits usage to AWS instance
	// with a similar product ID
	GetAWSProductID() string

	// SetAWSProductID sets AWS product ID
	SetAWSProductID(string)

	// GetAWSAccountID limits usage to AWS instance within account ID
	GetAWSAccountID() string

	// SetAWSAccountID sets AWS account ID that will be limiting
	// usage to AWS instance
	SetAWSAccountID(accountID string)

	// GetSupportsKubernetes returns kubernetes support flag
	GetSupportsKubernetes() Bool

	// SetSupportsKubernetes sets kubernetes support flag
	SetSupportsKubernetes(Bool)

	// SetLabels sets metadata labels
	SetLabels(labels map[string]string)

	// GetAccountID returns Account ID
	GetAccountID() string

	// GetPlanID returns Plan ID
	GetPlanID() string

	// CheckAndSetDefaults sets and default values and then
	// verifies the constraints for LicenseInfo.
	CheckAndSetDefaults() error
}

// NewLicenseInfo is a convenience method to to create LicenseInfoV3.
func NewLicenseInfo(name string, spec LicenseInfoSpecV3) (LicenseInfo, error) {
	return &LicenseInfoV3{
		Kind:    KindLicenseInfo,
		Version: V3,
		Metadata: Metadata{
			Name:      name,
			Namespace: defaults.Namespace,
		},
		Spec: spec,
	}, nil
}

// LicenseInfoV3 represents LicenseInfo resource version V3
type LicenseInfoV3 struct {
	// Kind is a resource kind - always resource.
	Kind string `json:"kind"`

	// Version is a resource version.
	Version string `json:"version"`

	// Metadata is metadata about the resource.
	Metadata Metadata `json:"metadata"`

	// Spec is the specification of the resource.
	Spec LicenseInfoSpecV3 `json:"spec"`
}

// GetName returns the name of the resource
func (c *LicenseInfoV3) GetName() string {
	return c.Metadata.Name
}

// SetLabels sets metadata labels
func (c *LicenseInfoV3) SetLabels(labels map[string]string) {
	c.Metadata.Labels = labels
}

// GetLabels returns metadata labels
func (c *LicenseInfoV3) GetLabels() map[string]string {
	return c.Metadata.Labels
}

// SetName sets the name of the resource
func (c *LicenseInfoV3) SetName(name string) {
	c.Metadata.Name = name
}

// Expiry returns object expiry setting
func (c *LicenseInfoV3) Expiry() time.Time {
	return c.Metadata.Expiry()
}

// SetExpiry sets object expiry
func (c *LicenseInfoV3) SetExpiry(t time.Time) {
	c.Metadata.SetExpiry(t)
}

// SetTTL sets Expires header using current clock
func (c *LicenseInfoV3) SetTTL(clock clockwork.Clock, ttl time.Duration) {
	c.Metadata.SetTTL(clock, ttl)
}

// GetMetadata returns object metadata
func (c *LicenseInfoV3) GetMetadata() Metadata {
	return c.Metadata
}

// GetReportsUsage returns true if teleport cluster reports usage
// to control plane
func (c *LicenseInfoV3) GetReportsUsage() Bool {
	return c.Spec.ReportsUsage
}

// SetReportsUsage sets usage report
func (c *LicenseInfoV3) SetReportsUsage(reports Bool) {
	c.Spec.ReportsUsage = reports
}

// CheckAndSetDefaults verifies the constraints for LicenseInfo.
func (c *LicenseInfoV3) CheckAndSetDefaults() error {
	return c.Metadata.CheckAndSetDefaults()
}

// GetAWSProductID returns product ID that limits usage to AWS instance
// with a similar product ID
func (c *LicenseInfoV3) GetAWSProductID() string {
	return c.Spec.AWSProductID
}

// SetAWSProductID sets AWS product ID
func (c *LicenseInfoV3) SetAWSProductID(pid string) {
	c.Spec.AWSProductID = pid
}

// GetPlanID sets AWS product ID
func (c *LicenseInfoV3) GetPlanID() string {
	return c.Spec.PlanID
}

// GetAccountID sets AWS product ID
func (c *LicenseInfoV3) GetAccountID() string {
	return c.Spec.AccountID
}

// GetAWSAccountID limits usage to AWS instance within account ID
func (c *LicenseInfoV3) GetAWSAccountID() string {
	return c.Spec.AWSAccountID
}

// SetAWSAccountID sets AWS account ID that will be limiting
// usage to AWS instance
func (c *LicenseInfoV3) SetAWSAccountID(accountID string) {
	c.Spec.AWSAccountID = accountID
}

// GetSupportsKubernetes returns kubernetes support flag
func (c *LicenseInfoV3) GetSupportsKubernetes() Bool {
	return c.Spec.SupportsKubernetes
}

// SetSupportsKubernetes sets kubernetes support flag
func (c *LicenseInfoV3) SetSupportsKubernetes(supportsK8s Bool) {
	c.Spec.SupportsKubernetes = supportsK8s
}

// String represents a human readable version of authentication settings.
func (c *LicenseInfoV3) String() string {
	var features []string
	if !c.Expiry().IsZero() {
		features = append(features, fmt.Sprintf("expires at %v", c.Expiry()))
	}
	if c.Spec.ReportsUsage.Value() {
		features = append(features, "reports usage")
	}
	if c.Spec.SupportsKubernetes.Value() {
		features = append(features, "supports kubernetes")
	}
	if c.Spec.AWSProductID != "" {
		features = append(features, fmt.Sprintf("is limited to AWS product ID %q", c.Spec.AWSProductID))
	}
	if c.Spec.AWSAccountID != "" {
		features = append(features, fmt.Sprintf("is limited to AWS account ID %q", c.Spec.AWSAccountID))
	}
	if len(features) == 0 {
		return ""
	}
	return strings.Join(features, ",")
}

// LicenseInfoSpecV3 is the actual data we care about for LicenseInfoV3.
type LicenseInfoSpecV3 struct {
	// AccountID is a customer account ID
	AccountID string `json:"account_id,omitempty"`
	// PlanID is a Plan ID
	PlanID string `json:"plan_id,omitempty"`
	// ReportsUsage is turned on when system reports usage
	ReportsUsage Bool `json:"usage,omitempty"`
	// AWSProductID limits usage to AWS instance with a product ID
	AWSProductID string `json:"aws_pid,omitempty"`
	// AWSAccountID limits usage to AWS instance within account ID
	AWSAccountID string `json:"aws_account,omitempty"`
	// SupportsKubernetes turns kubernetes support on or off
	SupportsKubernetes Bool `json:"k8s"`
}

// LicenseInfoSpecV3Template is a template for V3 LicenseInfo JSON schema
const LicenseInfoSpecV3Template = `{
  "type": "object",
  "additionalProperties": false,
  "properties": {
	"account_id": {
		"type": ["string"]
	},
	"plan_id": {
		"type": ["string"]
	},
	"usage": {
		"type": ["string", "boolean"]
	},
	"aws_pid": {
		"type": ["string"]
	},
	"aws_account": {
		"type": ["string"]
	},
	"k8s": {
		"type": ["string", "boolean"]
	}
  }
}`

// UnmarshalLicenseInfo unmarshals LicenseInfo from JSON or YAML
// and validates schema
func UnmarshalLicenseInfo(bytes []byte) (LicenseInfo, error) {
	var licenseInfo LicenseInfoV3

	if len(bytes) == 0 {
		return nil, trace.BadParameter("missing resource data")
	}

	schema := fmt.Sprintf(V2SchemaTemplate, MetadataSchema, LicenseInfoSpecV3Template, DefaultDefinitions)

	err := utils.UnmarshalWithSchema(schema, &licenseInfo, bytes)
	if err != nil {
		return nil, trace.BadParameter(err.Error())
	}

	if licenseInfo.Version != V3 {
		return nil, trace.BadParameter("unsupported version %v, expected version %v", licenseInfo.Version, V3)
	}

	if err := licenseInfo.CheckAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}

	return &licenseInfo, nil
}

// MarshalLicenseInfo marshals role to JSON or YAML.
func MarshalLicenseInfo(licenseInfo LicenseInfo, opts ...MarshalOption) ([]byte, error) {
	return json.Marshal(licenseInfo)
}
