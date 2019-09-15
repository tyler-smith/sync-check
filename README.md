# sync-check
[![Build Status](https://ci.tcrypt.dev/api/badges/tyler-smith/sync-check/status.svg)](https://ci.tcrypt.dev/tyler-smith/sync-check)
[![license](https://img.shields.io/github/license/tyler-smith/sync-check.svg?maxAge=2592000)](https://github.com/tyler-smith/sync-check/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tyler-smith/sync-check)](https://goreportcard.com/report/github.com/tyler-smith/sync-check)
[![GitHub issues](https://img.shields.io/github/issues/tyler-smith/sync-check.svg)](https://github.com/tyler-smith/sync-check/issues)

A health check appliance for Bitcoin Cash nodes. It checks a list of given nodes against each other and bitcoin.com and fails if any node is more than 3 blocks
behind the best block seen.

# Usage

```
sync-check grpc://bchd.greyh.at8335 rpc://localhost:8334
```

# Response

On success: Exits with code 0 and no output.

On failure: Exits with code equal to number of failing nodes and prints failing
nodes to stdout.