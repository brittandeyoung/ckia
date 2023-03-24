# CKIA
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md)
> **Warning**
> This project is currently a very early Work In Progress.

\[SEE\] + \[KEE\] + \[UH\]

ckia (cloud know it all) is a open source command line tool that is intended to run opinionated checks against your cloud environment and provide recommendations. The full suite of AWS trusted advisor checks is the inspiration for this project, but we intend to support many cloud providers when the product has fully matured. 

The key features of CKIA are:

- **Cloud Auditing**: Audit your cloud configuration for best practices, cost optimizations, performance Improvements, security misconfigurations, Fault Tolerance recommendation, and Service limits. 

# Available Checks

We are currently focused on duplicating the AWS Trusted advisor checks, but are willing to accept contributions to any cloud provider.

[Current List of Available Checks](AVAILABLE_CHECKS.md)

# Contributing

Please refer to our [Contributors Guide](CONTRIBUTORS.md)

# Installation Guide

Documentation is currently a work in progress. This README will contain the majority of the documentation for now. 

Currently to run the available checks:

1. clone the repo.
2. ensure you have aws credentials configured.
3. run the main.go file with the needed subcommands.

```
go run main.go aws check
```

## License

[Mozilla Public License v2.0](https://github.com/brittandeyoung/ckia/blob/main/LICENSE)