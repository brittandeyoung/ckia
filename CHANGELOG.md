# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- **New:** `Add progress bar when running checks.` [#57](https://github.com/brittandeyoung/ckia/issues/57)
- **New Flag:** `aws check --exclude-checks` [#49](https://github.com/brittandeyoung/ckia/issues/49)
- **New Flag:** `aws check --include-checks` [#49](https://github.com/brittandeyoung/ckia/issues/49)

## [0.1.0] - 2023-04-04
### Added
- **New Check:** `ckia:aws:cost:IdleLoadBalancers` [#44](https://github.com/brittandeyoung/ckia/issues/44)
- **New Check:** `ckia:aws:cost:UnassociatedElasticIPAddresses` [#22](https://github.com/brittandeyoung/ckia/issues/22)

## [0.0.1] - 2023-03-24
### Added
- **New Check:** `ckia:aws:security:RootAccountMissingMFA`
- **New Check:** `ckia:aws:cost:UnderutilizedEBSVolumes` [#9](https://github.com/brittandeyoung/ckia/issues/9)
- **New Check:** `ckia:aws:cost:IdleDBInstance`

[Unreleased]: https://github.com/brittandeyoung/ckia/compare/v0.1.0..HEAD
[0.1.0]: https://github.com/brittandeyoung/ckia/compare/v0.0.1..v0.1.0
[0.0.1]: https://github.com/brittandeyoung/ckia/tree/v0.0.1