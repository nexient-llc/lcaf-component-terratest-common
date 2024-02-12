package configure

import "github.com/gruntwork-io/terratest/modules/terraform"

func ConfigureTerraform(terraformDir string, varFiles []string) *terraform.Options {
	return &terraform.Options{
		TerraformDir: terraformDir,
		VarFiles:     varFiles,
	}
}
