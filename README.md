# CyberArk Client

`cac` is a simple [CyberArk Central Credentials Provider REST client](https://docs.cyberark.com/Product-Doc/OnlineHelp/AAM-CP/Latest/en/Content/CCP/Calling-the-Web-Service-using-REST.htm?tocpath=Developer%7CCentral%20Credential%20Provider%7CCall%20the%20Central%20Credential%20Provider%20Web%20Service%20from%20Your%20Application%20Code%7C_____2) written in Go.

![build](https://github.com/MartyHub/cac/actions/workflows/go.yml/badge.svg)

## Authentication

Authentication to CCP is done via client certificate/key files.

## Usage

```text
Usage:
  cac get <object>... [flags]

Flags:
      --app-id string         CyberArk Application Id
      --cert-file string      Certificate file
  -c, --config string         Config name
  -h, --help                  help for get
      --host string           CyberArk CCP REST Web Service Host
      --json                  JSON output
      --key-file string       Key file
      --max-connections int   Max connections (default 4)
      --max-tries int         Max tries (default 3)
      --safe string           CyberArk Safe
      --timeout duration      Timeout (default 30s)
      --wait duration         Wait before retry (default 100ms)
```

* Multiple objects can be queried in one go
* Pipe mode:
  ```shell
  $ echo "KEY=${CYBERARK:OBJECT}" | cac get -c MyConfig
  KEY=VALUE
  ```
* Use `max-connections` to set max parallel HTTP connections to CCP server
* Retries happened on error 500 (internal server error), 502 (bad gateway), 503 (service unavailable) and 504 (gateway
  timeout)
* A config file can be setup in `$XDG_CONFIG_HOME/cac/` (default to `~/.config/cac/`)
  * [Viper](https://github.com/spf13/viper#what-is-viper) supports JSON, TOML, YAML, HCL, envfile and Java properties config files
  * the config name parameter should be the name of the config file (without the extension)
  * command line parameters overwrite config
* `app-id`, `cert-file`, `host`, `key-file`, `safe` are required, either in config or in command line parameters
* To generate the autocompletion script for your favorite shell:
  ```shell
  $ cac completion (bash | fish | zsh)
  ```
* Default output is "shell" but JSON is also supported:
  ```shell
  o1='value1'
  o2='value2'
  ```

  ```json
  [
    {
      "object": "o1",
      "value": "value1"
    },
    {
      "object": "o2",
      "value": "value2"
    }
  ]
  ```
