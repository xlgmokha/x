# proxy-server

A tiny proxy server.


## Installation

```bash
$ go install github.com/xlgmokha/proxy-server@latest
```

## Usage

Start the server:

```bash
$ proxy-server
Listening and serving HTTP on http://127.0.0.1:8080
```

Use the proxy server:

```bash
$ curl -k --proxy 127.0.0.1:8080 https://www.eff.org/
```
