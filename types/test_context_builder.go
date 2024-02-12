package types

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

type TestContextBuilder struct {
	context *TestContext
}

func CreateTestContextBuilder() *TestContextBuilder {
	return &TestContextBuilder{
		context: NewTestContext(),
	}
}

func (b TestContextBuilder) SetTestConfig(config interface{}) *TestContextBuilder {
	b.context.testConfig = config

	return &b
}
func (b TestContextBuilder) SetTestSpecificFlags(flags map[string]TestFlags) *TestContextBuilder {
	b.context.testSpecificFlags = flags

	return &b
}

func (b TestContextBuilder) SetTestConfigFolderName(folderName string) *TestContextBuilder {
	b.context.testConfigFolderName = folderName

	return &b
}

func (b TestContextBuilder) SetTestConfigFileName(fileName string) *TestContextBuilder {
	b.context.testConfigFileName = fileName

	return &b
}

func (b TestContextBuilder) SetTerraformOptions(options *terraform.Options) *TestContextBuilder {
	b.context.terratestTerraformOptions = options

	return &b
}

func (b TestContextBuilder) SetCurrentTestName(testName string) *TestContextBuilder {
	b.context.currentTestName = testName

	return &b
}

func (b TestContextBuilder) validateContext() error {
	// Checks if the flags set by the client are in the allowed flags
	for test, testFlags := range b.context.testSpecificFlags {
		for flag := range testFlags {
			if !b.context.AllowedTestFlags().contains(flag) {
				return fmt.Errorf("test specific flag: %s is not allowed for test: %s", flag, test)
			}
		}
	}
	if len(b.context.TestConfigFolderName()) == 0 {
		return fmt.Errorf("TestConfigFolderName is not set")
	}

	if len(b.context.TestConfigFileName()) == 0 {
		return fmt.Errorf("TestConfigFileName is not set")
	}

	return nil
}

func (b TestContextBuilder) Build() *TestContext {
	err := b.validateContext()
	if err != nil {
		panic(err)
	}
	return b.context
}
