# Contributing to CKIA

First off, thanks for taking the time to contribute!


#### Table Of Contents

[Code of Conduct](CODE_OF_CONDUCT.md)

[I just have a question!!!](#i-have-a-question)

[Project Overview](#project-overview)

[How Can I Contribute?](#how-can-i-contribute)
  * [Open an Issue](#open-an-issue)
  * [Your First Code Contribution](#your-first-code-contribution)
  * [Pull Requests](#pull-requests)

[Additional Notes](#additional-notes)

## Code of Conduct

This project and everyone participating in it is governed by the [Atom Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [atom@github.com](mailto:atom@github.com).

## I have a question

For now Questions should still be created in the form of a GitHub issue. 

## What should I know before I get started?

### Project Overview

CKIA (cloud know it all) is a open source command line tool that is intended to run opinionated checks against your cloud environment and provide recommendations. The full suite of AWS trusted advisor checks is the inspiration for this project, but we intend to support many cloud providers when the product has fully matured. 

Our current focus is feature parity with AWS trusted advisor checks for the AWS cloud.

The checks for each cloud provider have the following required structure:
1. Checks within the project are centralized in a map for a particular cloud provider. A function will be used to centrally initialize these checks and will be stored in the file `internal/<cloud_provider>/checks.go` file. Below is an example of the checks map for the `internal/aws/checks.go` AWS checks. 

```go
func BuildChecksMap() map[string]interface{} {
	checksMap := checkMapping{
		// Cost Checks go here
		cost.IdleDBInstanceCheckId:          new(cost.IdleDBInstanceCheck),
		cost.UnderutilizedEBSVolumesCheckId: new(cost.UnderutilizedEBSVolumesCheck),
		// Security checks go here
		security.RootAccountMFACheckId: new(security.RootAccountMFACheck),
		
		// Additional Checks omitted
	}
	return checksMap
}
```
2. A file containing the logic and structures of the check located in the `internal/<cloud_provider>/<check_category/` directory. This file requires you define the following:
    - A constant for field defined in the common Check structure defined in `internal/common/common.go`. 
    - A structure containing the fields for the particular check.
    - A structure with the combination of the cental check strict and a list of the check structure.
    - A `List()` method defined for your Check structure. This Method must set the common check values to the defined constants and return the structure. (This is currently enforced with a unit test.)
    - A `Run()` method defined for your Check structure. This Method contains the logic for performing the check and building the Check object and returning the object to the runner. (This is currently enforced with a unit test.)
    - A separate or multiple separate `expand` function for any logic performed for the check. We separate this logic out from API calls in order to allow for easier unit testing. 
3. A `_test` file containing unit tests for any `expand` functions defined for your check. These checks should include multiple cases to ensure your expand function is operating as intended. 

## How Can I Contribute?

### Open an Issue

An issue can be created for new feature requests, submitting a bug, or simply asking a question. 

Before creating bug reports or a new feature request, please check open issues as you might find one already exists. [Open Issues](https://github.com/brittandeyoung/ckia/issues)

When reporting a bug or creating a feature, please fill out all areas of the bug template. 

### Your First Code Contribution

Unsure where to begin contributing to CKIA? You can start by looking through open issues for a feature or bug that you would like to work on. Once you have Identified what you would like to work on, please reference the [Project Overview](#project-overview) for the layout of the project.

#### Local development

Developing locally requires the following:

1. The repo is cloned locally.
2. You are running the verison of go defined in `go.mod` or later.
3. Cloud credentials for the cloud provider you want to write checks for ( So you can test the check once written.)

### Pull Requests

The process described here has several goals:

- Maintain CKIA's quality
- Fix problems that are important to users
- Engage the community in working toward the best possible CKIA
- Enable a sustainable system for CKIA's maintainers to review contributions

Please follow these steps to have your contribution considered by the maintainers:

1. Follow all instructions in [the template](PULL_REQUEST_TEMPLATE.md)
2. Follow the [styleguides](#styleguides)
3. After you submit your pull request, verify that all pull request checks are passing. If not, review the failures and work on resolving. If you are unable to resolve failing checks, please comment that you would like some help from the maintainers. 

While the prerequisites above must be satisfied prior to having your pull request reviewed, the reviewer(s) may ask you to complete additional design work, tests, or other changes before your pull request can be ultimately accepted.


## Additional Notes

This projects CHANGELOG format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/). Maintainers will handle updating this during the review of your pull request.

This is a first attempt at our contributors guide following the examples of many other open source projects. If this document needs to be updated or you have questions, please feel free to raise an issue. 