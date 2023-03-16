# CKIA

**Warning**
This project is currently a very early Work In Progress.

\[SEE\] + \[KEE\] + \[UH\]

ckia (cloud know it all) is a open source command line tool that is intended to run opinionated checks against your cloud environment and provide recommendations. The full suite of AWS trusted advisor checks is the inspiration for this project, but we intend to support many cloud providers when the product has fully matured. 

The key features of CKIA are:

- ** Cloud Auditing **: Audit your cloud configuration for best practices, cost optimizations, performance Improvements, security misconfigurations, Fault Tolerance recommendation, and Service limits. 

# Available Checks

We are currently focused on duplicating the AWS Trusted advisor checks, but are willing to accept contributions to any cloud provider.

| Id | Provider | Check Category |  Name | Rule Description|
|--|----------|------------|--------------------------------------|------------------------------------------------------------------------|
| ckia:aws:cost:IdleDBInstanceCheck | AWS      | Cost Optimization | IdleDBInstances |  Any RDS DB instance that has not had a connection in the last 7 days is considered idle. |                                                     |

# Contributing

As this project matures the contribution process will be updated. For now please feel free to fork the repo and provide a pull request for review. 

# Getting Started & Documentation

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