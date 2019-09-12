Shuttle
========

[![MIT](https://img.shields.io/badge/license-MIT-brightgreen.svg)](./LICENSE)

- [1 Background](#1-background)
- [2 Install](#2-install)
  - [2.1 Requirements](#21-requirements)
  - [2.2 Install bytom node](#22-install-bytom-node)
  - [2.3 Build from source code](#23-build-from-source-code)
- [3 Usage](#3-usage)
  - [3.1 Launch bytom node](#31-launch-bytom-node)
  - [3.2 Create account and issue your asset](#32-create-account-and-issue-your-asset)
  - [3.3 Deploy contract](#33-deploy-contract)
  - [3.5 Call contract](#35-call-contract)
- [4 Contributing](#4-contributing)
- [5 License](#5-license)

## 1 Background

Shuttle is designed to help swap different assets in bytom.

## 2 Install

### 2.1 Requirements

- [Go](https://golang.org/doc/install) version 1.12 or higher, with `$GOPATH` set to your preferred directory

### 2.2 Install bytom node

Firstly, you should install and configure bytom node, see also: [Bytom repository](https://github.com/Bytom/bytom).

### 2.3 Build from source code

This BTM swap tool is still in beta, so repository code will be changed frequently. You can build tool from source code directly.

```shell
$ git clone https://github.com/Bytom/shuttle.git $GOPATH/src/github.com/shuttle
$ cd $GOPATH/src/github.com/shuttle
$ make install
```

or remove shuttle:

```shell
$ cd $GOPATH/src/github.com/shuttle
$ make clean
```

## 3 Usage

### 3.1 Launch bytom node

For testing, you can launch bytom solonet node.

```shell
$ bytomd init --chain_id=solonet --home $HOME/bytom/solonet # init bytom solonet node
$ bytomd node --home $HOME/bytom/solonet --mining           # launch bytom solonet node and start mining
```

### 3.2 Create account and issue your asset

You should create several accounts and issue your asset for testing, more details:

- [Managing Accounts](https://github.com/Bytom/bytom/wiki/Managing-Accounts)
- [Assets registration](https://github.com/Bytom/bytom/wiki/Advanced-Transaction#assets-registration)

### 3.3 Deploy contract

```shell
$ cd $GOPATH/src/github.com/btm-swap-tool/cmd
$ ./cmd deploy 10CJPO1HG0A02 12345 --amountLocked 20000000000 --amountRequested 1000000000 --assetLocked bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a --assetRequested ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff --cancelKey 3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e --seller 00145dd7b82556226d563b6e7d573fe61d23bd461c1f --txFee 40000000
--> contractUTXOID: ad1e72021352ade9e0fef7b1c52cc070bfcb1429f885a540f00f8b957e941b2d
```

Then, wait about 2.5 minutes, and a new block will be mined, the contract will be confirmed.

### 3.5 Call contract

```shell
$ ./cmd call 10CKAD3000A02 12345 00140fdee108543d305308097019ceb5aec3da60ec66 ad1e72021352ade9e0fef7b1c52cc070bfcb1429f885a540f00f8b957e941b2d
--> txID: 762d6912e126ac3937cee54db8898af09abc5633c8b78682fa0cc23d89a518a9
```

When the transaction will be confirmed in a new block, the whole BTM swap is successful.

## 4 Contributing

Welcome to [open an issue](https://github.com/Bytom/btm-swap-tool/issues/new) or submit PRs. This project exists thanks to all the people who contribute.

## 5 License

[MIT](./LICENSE) © 2019 Bytom
