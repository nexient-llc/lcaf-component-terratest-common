package types

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

type GenericTFModuleConfig struct {
	//the framework standard subset of attributes
	Naming_prefix      string            `json:"naming_prefix"`
	Environment        string            `json:"environment"`
	Environment_number string            `json:"environment_number"`
	Resource_number    string            `json:"resource_number"`
	Tags               map[string]string `json:"tags"`
	//to be extended by the TF module specific attrs
}

type TestContext struct {
	testConfig                any // pointer to a TF module specific inheritance of GenericTFModuleConfig
	testConfigFolderName      string
	testConfigFileName        string
	terratestTerraformOptions *terraform.Options
	currentTestName           string
	testSpecificFlags         map[string]TestFlags
	allowedTestFlags          AllowedTestFlags
}

type TestFlags map[string]bool
type AllowedTestFlags []string

func (t AllowedTestFlags) contains(element string) bool {
	for _, value := range t {
		if value == element {
			return true
		}
	}
	return false
}

func defaultAllowedTestFlags() AllowedTestFlags {
	return []string{"SKIP_TEST", "IS_TERRAFORM_IDEMPOTENT_APPLY"}

}

func NewTestContext() *TestContext {
	return &TestContext{
		allowedTestFlags: defaultAllowedTestFlags(),
	}
}

func (ctx *TestContext) TestSpecificFlags() map[string]TestFlags {
	return ctx.testSpecificFlags
}

func (ctx *TestContext) TestConfig() any {
	return ctx.testConfig
}

func (ctx *TestContext) SetTestConfig(config any) {
	ctx.testConfig = config
}

func (ctx *TestContext) TestConfigFolderName() string {
	return ctx.testConfigFolderName
}

func (ctx *TestContext) TestConfigFileName() string {
	return ctx.testConfigFileName
}

func (ctx *TestContext) TerratestTerraformOptions() *terraform.Options {
	return ctx.terratestTerraformOptions
}

func (ctx *TestContext) SetTerratestTerraformOptions(options *terraform.Options) {
	ctx.terratestTerraformOptions = options

}

func (ctx *TestContext) AllowedTestFlags() AllowedTestFlags {
	return ctx.allowedTestFlags
}

func (ctx *TestContext) CurrentTestName() string {
	return ctx.currentTestName
}

func (ctx *TestContext) SetCurrentTestName(testName string) {
	ctx.currentTestName = testName
}

func (ctx *TestContext) IsCurrentTest(testName string) bool {
	return ctx.currentTestName == testName
}

func (ctx *TestContext) EnabledOnlyForTests(t *testing.T, testName ...string) {
	for _, testName := range testName {
		if ctx.currentTestName == testName {
			return
		}
	}
	t.SkipNow()
}

type SecurityGroupT struct {
	EgressWithCidrBlocks []struct {
		CidrBlocksCommaSeparated string `json:"cidr_blocks"`
		CidrBlocks               []string
		FromPort                 int    `json:"from_port"`
		Protocol                 string `json:"protocol"`
		ToPort                   int    `json:"to_port"`
	} `json:"egress_with_cidr_blocks"`
	IngressCidrBlocks []string `json:"ingress_cidr_blocks"`
	IngressRules      []string `json:"ingress_rules"`
}
