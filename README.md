# CyberArk Client

`cac` is a
simple [CyberArk Central Credentials Provider REST client](https://docs.cyberark.com/Product-Doc/OnlineHelp/AAM-CP/Latest/en/Content/CCP/Calling-the-Web-Service-using-REST.htm?tocpath=Developer%7CCentral%20Credential%20Provider%7CCall%20the%20Central%20Credential%20Provider%20Web%20Service%20from%20Your%20Application%20Code%7C_____2)
written in Go.

![build](https://github.com/MartyHub/cac/actions/workflows/go.yml/badge.svg)

## Authentication

Authentication to CCP is done via client certificate/key files.

## Features

* Handle multiple configurations
* Multiple connections
* Automatic retries on error 500 (internal server error), 502 (bad gateway), 503 (service unavailable) and 504 (gateway
  timeout)
* Manage a cache
* Shell, JSON or file output

## Installation

To generate the autocompletion script for your favorite shell:

```shell
cac completion (bash | fish | zsh)
```

## Configuration

To add or update a configuration:

```text
cac config set <config> [flags]

Flags:
--aliases strings       Aliases
--app-id string         CyberArk Application Id
--cert-file string      Certificate file
--expiry duration       Cache expiry (default 12h0m0s)
--host string           CyberArk CCP REST Web Service Host
--key-file string       Key file
--max-connections int   Max connections (default 4)
--max-tries int         Max tries (default 3)
--safe string           CyberArk Safe
--skip-verify           Skip server certificate verification
--timeout duration      Timeout (default 30s)
--wait duration         Wait before retry (default 100ms)
```

A configuration has a main `<config>` name but can also have aliases

## Usage

To get accounts from CyberArk:

```text
cac get <config> <account>... [flags]

Flags:
  -j, --json            Output JSON
  -o, --output string   Generate files in given output path
```

Using pipe, the behavior is to look for accounts using a regular expression `${CYBERARK:XXX}`:

```shell
$  echo 'KEY=${CYBERARK:MY_ACCOUNT}' | cac get test
KEY=MY_ACCOUNT_PASSWORD
```
