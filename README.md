# CyberArk Client

`cac` is a
simple [CyberArk Central Credentials Provider REST client](https://docs.cyberark.com/Product-Doc/OnlineHelp/AAM-CP/Latest/en/Content/CCP/Calling-the-Web-Service-using-REST.htm?tocpath=Developer%7CCentral%20Credential%20Provider%7CCall%20the%20Central%20Credential%20Provider%20Web%20Service%20from%20Your%20Application%20Code%7C_____2)
written in Go.

![build](https://github.com/MartyHub/cac/actions/workflows/go.yml/badge.svg)

## Authentication

Authentication to CCP is done via client certificate/key files.

## Usage

```text
Usage of ./cac:
  -appId string
    	CyberArk Application Id
  -certFile string
    	Certificate file
  -config string
    	Config name
  -host string
    	CyberArk CCP REST Web Service Host
  -json
    	JSON output
  -keyFile string
    	Key file
  -maxConns int
    	Max connections (default 4)
  -maxTries int
    	Max tries (default 3)
  -object value
    	CyberArk Object (at least one required)
  -safe string
    	CyberArk Safe
  -timeout duration
    	Timeout (default 30s)
  -version
    	Display version information
  -wait duration
    	Wait before retry (default 100ms)
```

* Multiple objects can be queried in one go: `--object o1 --object o2`
* Pipe mode:
  ```shell
  $ echo "KEY=${CYBER_ARK:OBJECT}" | cac -config MyConfig
  KEY=VALUE
  ```
* Use `maxConns` to set max parallel HTTP connections to CCP server
* Retries happened on error 500 (internal server error), 502 (bad gateway), 503 (service unavailable) and 504 (gateway
  timeout)
* A config file can be setup in `$XDG_CONFIG_HOME/cac/config.json` (default to `~/.config/cac/config.json`)
  * command line parameters overwrite config
* `appId`, `certFile`, `host`, `keyFile`, `safe` are required, either in config or in command line parameters
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