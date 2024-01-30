# tf-caf-terratest-common

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)

## Overview

Terratest support utitities and test runners supporting Common Automation Framework (CXAF) terraform modules running automated tests in pipelines.

Goals:
1. To keep infra test code DRY and composable, reusable functions have been extracted into this dedicated repo which can be included by TF module tests
2. Tests are configuration driven, to be reusable and aggregatable for higher level customer project specific integration testing
3. Configuration is shared across infra deployment (terraform) and infra test (terratest) automation
3. Automated pipeline friendly. Configuration switches are driven by OS env vars

## Usage

By default test suites are pointed to `<Module Repo>/examples` and expect configuration variables in test.tfvars

```
make check
```

To point to customer project specific terraform code:

```
DSO_INFRA_TEST_CONFIG_FOLDER=/projects/abc/ make check
```

To override default test.tfvars

```
DSO_INFRA_TEST_CONFIG_FOLDER=/projects/abc/ DSO_INFRA_TEST_CONFIG_TFVAR_FILENAME=project.tfvars make check
```

### Stages

Tests and individual test stages can be skipped. Example:

```
<tf module repo>/
	examples/
		ecs_example/
			main.tf
			test.tfvars
		eks_example/
			main.tf
			test.tfvars
```

To skip a test:

```
DSO_INFRA_TEST_SKIP_TEST_<name of test TF folder> make check
# from layout above:
DSO_INFRA_TEST_SKIP_TEST_ecs_example make check
```

To disable selected stage(s) of the test:

```
SKIP_teardown_test_eks_example=y make check
```

TODO - to pickup from multi tfvars to align to any project file naming convention

## Reuse

We want reuse the same test implementation for a TF module development and for regression testing of a project that includes that module, probably among multi other ones.

Not every test can be reused. A low-level primitive TF module, like "DNS record", has to have an extensive test fixture.
We solve it by introducing a naming convention for GoLang tests. Those safe to be composed/reused from higher level pipelines have a `TestComposable` prefix in their GoLang test name.

Example:

```
=== ECS-Application-module/tests/testimpl.go ===
func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	...
	assert.Equal(t, ctx.TestConfig.(*ThisTFModuleConfig).dockerImage, getAWSEcsAPI().FargateApp(appArn).Container().ImageName)

}
====
```

## Examples

### Launch test suite in ReadOnly mode - part of after deployment regression test
No cloud resources will be created nor teared down
```
tf-module-skeleton $ make go/readonly_test
```

### Launch test suite in Regular mode
```
tf-module-skeleton $ make go/test
```

### Many to many relation between tests and IaC being tested

```
<repo>/
	xyz_project_test/
		private_network/
			main.tf
			test.tfvars
		private_network_and_no_egress/
			main.tf
			test.tfvars
		private_network_and_abc/
			main.tf
			test.tfvars
		public_network/
			main.tf
			test.tfvars
```

```

func TestFeatureABC_1(t *testing.T, ctx types.TestContext) {
	t.Run("OnlyPrivateNetworks/TestIfAPPisUP", func(t *testing.T) {
		ctx.EnabledOnlyForTests(t, "private_network_and_no_egress","private_network_and_abc")
		//^ this test will be run only for terraform code in folders "private_network_and_no_egress" or  "private_network_and_abc"
		remoteAgent := launchAgentInsidePrivateNetwork( ctx.TestConfig.(*ThisTFModuleConfig).network)
		assertHTTP_200_OK(remoteAgent.sendHTTPRequest2Target( ctx.TestConfig.(*ThisTFModuleConfig).InternalURL).getStatusCode)

	})
}
func TestFeatureABC_2(t *testing.T, ctx types.TestContext) {
	t.Run("Basic/TestIfAPPisUP", func(t *testing.T) {
		ctx.EnabledOnlyForTests(t, "public_network")
		// This test code requires infra be in public network
		assertHTTP_200_OK(sendHTTPRequest2Target( ctx.TestConfig.(*ThisTFModuleConfig).PublicURL).getStatusCode)
	})
}
```

### Enable/disable subset of tests

Leveraging GoLang test utilities inherited by this "framework"
https://pkg.go.dev/testing#hdr-Subtests_and_Sub_benchmarks
```
go test -run ''        # Run all tests.
go test -run Foo       # Run top-level tests matching "Foo", such as "TestFooBar".
go test -run Foo/A=    # For top-level tests matching "Foo", run subtests matching "A=".
go test -run /A=1      # For all top-level tests, run subtests matching "A=1".
go test -fuzz FuzzFoo  # Fuzz the target matching "FuzzFoo"
```
```
// tests/post_deploy_functional/main_tests.go
func TestCommon(t *testing.T) {

	ctx := types.TestContext{
		TestConfig: &testimpl.ThisTFModuleConfig{},
	}
	lib.RunSetupTestTeardown(t, testConfigsFolder, infraTFVarFileNameDefault, ctx,
		testimpl.TestXYZ)
}
...
// tests/testimpl/test_impl.go
func TestXYZ(t *testing.T, ctx types.TestContext) {
	t.Run("Basic/AzureManagedIdentityON/abc", func(t *testing.T) {
		...
	})
	t.Run("Basic/AzureManagedIdentityOFF/abc", func(t *testing.T) {
		...
	})

```

```
$ cd tests/post_deploy_functional
go test -run Common # runs all tests from "Common"
go test -run /Basic/AzureManagedIdentityON # runs all subtests from Basic category that requires Azure Managed Identity be enabled
go test -run /AzureManagedIdentityON # runs all subtests any category that requires Azure Managed Identity be enabled

```

### One-to-one relationship between examples and tests
There would be a few scenarios where users would like to run specific tests for each example. Although this framework natively doesn't support doing that. However, as a work-around we can achieve in our `main_test.go` as follows

```go
func TestKubernetesModule(t *testing.T) {
    // Provide a map of examples to the tests
	examplesToTestsMap := map[string]lib.TestFunc{
		"private-cluster": testimpl.TestPrivateCluster,
		"public-cluster":  testimpl.TestPublicCluster,
	}
	ctx := types.TestContext{
		TestConfig:                 &testimpl.ThisTFModuleConfig{},
		IsTerraformIdempotentApply: false,
	}
	// Loop through the examples
	for example, testFunction := range examplesToTestsMap {
		lib.RunSetupTestTeardown(t, testConfigsExamplesFolderDefault+"/"+example, infraTFVarFileNameDefault, ctx, testFunction)
	}

}
```

Currently, these tests run sequentially. Making to run them in parallel can be a future optimization.

### Run non-idempotent terraform apply
There are a few scenarios where users would like to run `terraform.initAndApply()` instead of `terraform.InitAndApplyAndIdempotent()`, mostly because of bugs in the providers which doesn't support idempotent applies. This can be done by setting the flag `IsTerraformIdempotentApply` in the context as shown below
```go
ctx := types.TestContext{
    TestConfig:                 &testimpl.ThisTFModuleConfig{},
    IsTerraformIdempotentApply: false,
}
```

### Set timeout for go test
The default timeout of go test is `20 mins` which may not be enough for running some heavy tests. If timeout is reached, it may leave resources provisioned in the cloud and cost us money. Simple way is to increase the timeout during running go tests
```go
go test main_test.go -timeout 1h
```
## References
[Terratest best practices](https://terratest.gruntwork.io/docs/#testing-best-practices)

[GoLang test framework ](https://pkg.go.dev/testing)

## Diagrams

[Overview](doc/Overview.svg)


### local development

To test amendments to the terratest helper before those committed to github, use GoLang "replace". Example
```
module github.com/nexient-llc/tf-aws-module-private_dns_namespace

go 1.20

replace github.com/nexient-llc/tf-caf-terratest-common => /Home/user/CAF/NOT_CHECKED_IN_YET/tf-caf-terratest-common

require (
	github.com/nexient-llc/tf-caf-terratest-common v0.0.0-00010101000000-000000000000
)
```

### GoLang

To use "github.com/nexient-llc" private repository when developing or running GoLang code:

```
go env -w GOPRIVATE='github.com/nexient-llc/'
```

### Pipeline integration

For unattended CI//CD pipelines, you must pre-authenticate to Github.

#### HTTPS authentication

```
git config --add --global url."https://oauth2:$GITHUB_PTA_TOKEN@github.com/".insteadOf "https://github.com/"
```

#### SSH authentication

```
git config --add --global url."ssh://git@github.com/".insteadOf "https://github.com/"
```

## Prerequisites

- [asdf](https://github.com/asdf-vm/asdf) used for tool version management
- [make](https://www.gnu.org/software/make/) used for automating various functions of the repo
- [repo](https://android.googlesource.com/tools/repo) used to pull in all components to create the full repo template

### Repo Init

Run the following commands to prep repo and enable all `Makefile` commands to run

```shell
asdf plugin add conftest
asdf plugin add golang
asdf plugin add golangci-lint
asdf plugin add pre-commit
asdf plugin add terraform
asdf plugin add terraform-docs
asdf plugin add tflint

asdf install
```

## Pre-Commit hooks

A [.pre-commit-config.yaml](.pre-commit-config.yaml) file defines certain `pre-commit` hooks that are relevant to terraform, golang and common linting tasks. There are no custom hooks added.

`commitlint` hook enforces that commit messages in a certain format (see [Conventional Commits](https://www.conventionalcommits.org/)). The commit message must contain the following structural elements, to communicate intent to the consumers of your commits:

- **fix**: a commit of the type `fix` patches a bug in your codebase (this correlates with PATCH in Semantic Versioning).
- **feat**: a commit of the type `feat` introduces a new feature to the codebase (this correlates with MINOR in Semantic Versioning).
- **BREAKING CHANGE**: a commit that has a footer `BREAKING CHANGE:`, or appends a `!` after the type/scope, introduces a breaking API change (correlating with MAJOR in Semantic Versioning). A BREAKING CHANGE can be part of commits of any type.
footers other than BREAKING CHANGE: <description> may be provided and follow a convention similar to git trailer format.
- **build**: a commit of the type `build` adds changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
- **chore**: a commit of the type `chore` adds changes that don't modify src or test files
- **ci**: a commit of the type `ci` adds changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)
- **docs**: a commit of the type `docs` adds documentation only changes
- **perf**: a commit of the type `perf` adds code change that improves performance
- **refactor**: a commit of the type `refactor` adds code change that neither fixes a bug nor adds a feature
- **revert**: a commit of the type `revert` reverts a previous commit
- **style**: a commit of the type `style` adds code changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **test**: a commit of the type `test` adds missing tests or correcting existing tests

Base configuration used for this project is [commitlint-config-conventional (based on the Angular convention)](https://github.com/conventional-changelog/commitlint/tree/master/@commitlint/config-conventional#type-enum)

If you are a developer using vscode, the [commitlint](https://marketplace.visualstudio.com/items?itemName=joshbolduc.commitlint) plugin may be helpful.

`detect-secrets-hook` prevents new secrets from being introduced into the baseline. [TODO: INSERT DOC LINK ABOUT HOOKS]

In order for `pre-commit` hooks to work properly:

- You need to have the pre-commit package manager installed. [Here](https://pre-commit.com/#install) are the installation instructions.
- `pre-commit` would install all the hooks when commit message is added by default except for `commitlint` hook. `commitlint` hook would need to be installed manually using the command below

```
pre-commit install --hook-type commit-msg
```

## To run a local quality check

1. For development/enhancements to this module locally, you'll need to install all of its components. This is controlled by the `configure` target in the project's [`Makefile`](./Makefile). Before you can run `configure`, familiarize yourself with the variables in the `Makefile` and ensure they're pointing to the right places.

```
make configure
```

This adds in several files and directories that are ignored by `git`. They expose many new Make targets.

2. The first target you care about is `check`.
If the `make check` target is successful, the developer can commit the code to git.

`make check` target

- runs `terraform commands` to `lint`, `validate` and `plan` terraform code.
- runs `conftests`. `conftests` make sure `policy` checks are successful.
- runs `terratest`. This is integration test suite.
