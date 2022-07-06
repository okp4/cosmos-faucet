# CØSMOS Faucet

[![version](https://img.shields.io/github/v/release/okp4/cosmos-faucet)](https://github.com/okp4/cosmos-faucet/releases)
[![build](https://github.com/okp4/cosmos-faucet/actions/workflows/build.yml/badge.svg)](https://github.com/okp4/cosmos-faucet/actions/workflows/build.yml)
[![lint](https://github.com/okp4/cosmos-faucet/actions/workflows/lint.yml/badge.svg)](https://github.com/okp4/cosmos-faucet/actions/workflows/lint.yml)
[![test](https://github.com/okp4/cosmos-faucet/actions/workflows/test.yml/badge.svg)](https://github.com/okp4/cosmos-faucet/actions/workflows/test.yml)
[![maintainability](https://api.codeclimate.com/v1/badges/b2b9effa4c2f43ffbf3d/maintainability)](https://codeclimate.com/github/okp4/cosmos-faucet/maintainability)
[![coverage](https://codecov.io/gh/okp4/cosmos-faucet/branch/main/graph/badge.svg?token=0VQHJDMY5B)](https://codecov.io/gh/okp4/cosmos-faucet)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![license](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

This faucet app allows anyone to request tokens for a [Cosmos](https://cosmos.network/) account address, either in command line or in service mode, with optional [ReCaptcha](https://www.google.com/recaptcha/about/) security.

The faucet app supports any Cosmos based blockchain and is intended to be configured on testnet networks.

## Install

```sh
go install github.com/okp4/cosmos-faucet@latest

cosmos-faucet --help
```

## Config

 Configuration can be passed as CLI argument, configuration file (`config.yml`) or by environment variable prefixed by `FAUCET` (i.e.: `FAUCET_MNEMONIC`).

```yml
grpc-address: 127.0.0.1:9090
mnemonic: "kiwi nuclear blast wet badge..."
chain-id: testnet-1
denom: uatom
prefix: cosmos
fee-amount: 0
amount-send: 1
memo: "Sent by økp4 faucet"
gas-limit: 200000
```

## Usage

### CLI

#### Send

```shell
Send tokens to a given address

Usage:
  cosmos-faucet send <address> [flags]

Flags:
  -h, --help   help for send

Global Flags:
      --amount-send int       Amount send value (default 1)
      --chain-id string       The network chain ID (default "localnet-okp4-1")
      --denom string          Token denom (default "know")
      --fee-amount int        Fee amount
      --gas-limit uint        Gas limit (default 200000)
      --grpc-address string   The grpc okp4 server url (default "127.0.0.1:9090")
      --memo string           The memo description (default "Sent by økp4 faucet")
      --mnemonic string       
      --no-tls                No encryption with the GRPC endpoint
      --prefix string         Address prefix (default "okp4")
      --tls-skip-verify       Encryption with the GRPC endpoint but skip certificates verification
```

#### Start

```shell
Start the GraphQL api

Usage:
  cosmos-faucet start [flags]

Flags:
      --address string              graphql api address (default ":8080")
      --captcha                     enable captcha verification
      --captcha-min-score float     set Captcha min score (default 0.5)
      --captcha-secret string       set Captcha secret
      --captcha-verify-url string   set Captcha verify URL (default "https://www.google.com/recaptcha/api/siteverify")
      --health                      enable health endpoint
  -h, --help                        help for start
      --metrics                     enable metrics endpoint

Global Flags:
      --amount-send int       Amount send value (default 1)
      --chain-id string       The network chain ID (default "localnet-okp4-1")
      --denom string          Token denom (default "know")
      --fee-amount int        Fee amount
      --gas-limit uint        Gas limit (default 200000)
      --grpc-address string   The grpc okp4 server url (default "127.0.0.1:9090")
      --memo string           The memo description (default "Sent by økp4 faucet")
      --mnemonic string       
      --no-tls                No encryption with the GRPC endpoint
      --prefix string         Address prefix (default "okp4")
      --tls-skip-verify       Encryption with the GRPC endpoint but skip certificates verification
```

### GraphQL

Start GraphQL server with captcha verification for the send token mutation.

```shell
cosmos-faucet start --captcha --captcha-secret $CAPCTHA_SECRET
```

Access on playground and documentation at the root of server.

## Build

The project comes with a convenient `Makefile` which depends on [Docker](https://www.docker.com). Please verify that Docker is properly installed and if not, follow the instructions:

- for macOS: <https://docs.docker.com/docker-for-mac/install/>
- for Windows: <https://docs.docker.com/docker-for-windows/install/>
- for Linux: <https://docs.docker.com/engine/install/>

To build the app, invoke the following goal `build` of the `Makefile`:

```sh
make build
```

The app will be generated under the folder `target/dist`.

## :heart: Supporting this project & Contributing

A simple star in this repository is enough to make us happy!

But you're also welcome to contribute! We appreciate any help you're willing to give. Don't hesitate to open issues and/or submit pull requests.

Remember that this is the Cosmos Token Faucet we use at [OKP4](http://okp4.network). This is why we may have to refuse change requests simply because they do not comply with our internal requirements, and not because they are not relevant.
