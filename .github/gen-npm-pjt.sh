#! /bin/bash

set -eux

VERSION=$(cat version)

cat <<EOF > package.json
{
  "name": "cosmos-faucet-schema",
  "version": "${VERSION}",
  "private": false,
  "description": "Cosmos Faucet GraphQL schema.",
  "repository": {
    "type": "git",
    "url": "git@github.com:okp4/cosmos-faucet.git"
  },
  "keywords": [
    "cosmos",
    "faucet",
    "graphql"
  ],
  "author": {
    "name": "okp4",
    "web": "https://okp4.network"
  },
  "license": "BSD-3-Clause",
  "bugs": {
    "url": "https://github.com/okp4/cosmos-faucet/issues"
  },
  "homepage": "https://github.com/okp4/cosmos-faucet#readme",
  "files": [
    "graph/schema.graphqls"
  ]
}
EOF
