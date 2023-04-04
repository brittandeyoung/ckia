# CKIA
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md)
> **Warning**
> This project is currently a very early Work In Progress. Please pay attention to releases for potential breaking changes. At this time we are not providing Cost Savings Estimations, we are still working out how to best provide these numbers in the most accurate manner with the given AWS apis. 

\[SEE\] + \[KEE\] + \[UH\]

ckia (cloud know it all) is a open source command line tool that is intended to run opinionated checks against your cloud environment and provide recommendations. The full suite of AWS trusted advisor checks is the inspiration for this project, but we intend to support many cloud providers when the product has fully matured. Our current focus is providing check parity with the current AWS Trusted Advisor offerings. 

The key features of CKIA are:

- **Cloud Auditing and Recomendations**: Audit your cloud configuration for:

1. best practices
2. cost optimizations
3. performance Improvements
4. security misconfigurations
5. Fault Tolerance recommendation
6. Service limits

Each Check additionally provides recommended actions for when a check is failing. 

# Available Checks

We are currently focused on duplicating the AWS Trusted advisor checks, but are willing to accept contributions to any cloud provider.

[Current List of Available Checks](AVAILABLE_CHECKS.md)

# Contributing

Please refer to our [Contributors Guide](CONTRIBUTING.md)

# Installation Guide

Documentation is currently a work in progress. This README will contain the majority of the documentation for now. 

To install the latest version of CKIA, download the correct package for your operating system and architecture from the latest release. Below is an example for automating this process for Linux. Modify this to match your OS and architecture:

```shell
LATEST_VERSION=$(wget -O - "https://api.github.com/repos/brittandeyoung/ckia/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | cut -c 2-)
wget "https://github.com/brittandeyoung/ckia/releases/download/v${LATEST_VERSION}/ckia_${LATEST_VERSION}_Linux_amd64.tar.gz"
tar xzf ckia_${LATEST_VERSION}_Linux_amd64.tar.gz
chmod +x ckia
sudo mv ckia /usr/local/bin
```

Then you can run the application:

```shell
✗ ckia
An open source tool for making recommendations for target cloud account.

Usage:
  ckia [command]

Available Commands:
  aws         Checks related to the aws cloud.
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.ckia.yaml)
  -h, --help            help for ckia

Use "ckia [command] --help" for more information about a command.
```

```shell
✗ ckia aws -h
Checks related to the aws cloud.

Usage:
  ckia aws [flags]
  ckia aws [command]

Available Commands:
  check       Run available checks for aws
  list        List available checks for aws

Flags:
  -h, --help   help for aws

Global Flags:
      --config string   config file (default is $HOME/.ckia.yaml)

Use "ckia aws [command] --help" for more information about a command.
```

## License

[Mozilla Public License v2.0](https://github.com/brittandeyoung/ckia/blob/main/LICENSE)