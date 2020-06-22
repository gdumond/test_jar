# Proxy Test

Tests that a connection may be made through a given SOCKS5 proxy

## Arguments

| Name     | Default  | Description                                                                  |
| -------- | -------- | ---------------------------------------------------------------------------- |
| proxy    |          | The host:port to use for the proxy config, default port is 1080              |
| url      |          | The URL to connect to once auth through proxy has succeeded                  |
| username | username | The username to use when attempting username/password auth through the proxy |
| password | password | The password to use when attempting username/password auth through the proxy |

## Sample Outcomes

### No Auth on proxy

```bash
$ ./proxyTest -proxy=192.168.56.121 -url=ingestion-lumberjack.condor.staging.axwaytest.net:453
Attempting connection through Proxy without authentication
Connection to proxy made, connecting to url now
connection without authentication was successful
Attempting connection through Proxy with username/password
Connection to proxy made, connecting to url now
connection with username/password authentication failed socks connect tcp 192.168.56.121:1080->ingestion-lumberjack.condor.staging.axwaytest.net:453: no acceptable authentication methods
Attempting connection through Proxy with gss-api
Connection to proxy made, connecting to url now
connection with gss-api authentication failed socks connect tcp 192.168.56.121:1080->ingestion-lumberjack.condor.staging.axwaytest.net:453: no acceptable authentication methods
```

### Username/Password on proxy

#### Authentication Details Incorrect

```bash
$ ./proxyTest -proxy=192.168.56.121 -url=ingestion-lumberjack.condor.staging.axwaytest.net:453
Attempting connection through Proxy without authentication
Connection to proxy made, connecting to url now
connection without authentication failed socks connect tcp 192.168.56.121:1080->ingestion-lumberjack.condor.staging.axwaytest.net:453: no acceptable authentication methods
Attempting connection through Proxy with username/password
Connection to proxy made, connecting to url now
connection with username/password authentication failed socks connect tcp 192.168.56.121:1080->ingestion-lumberjack.condor.staging.axwaytest.net:453: username/password authentication failed
Attempting connection through Proxy with gss-api
Connection to proxy made, connecting to url now
connection with gss-api authentication failed socks connect tcp 192.168.56.121:1080->ingestion-lumberjack.condor.staging.axwaytest.net:453: no acceptable authentication methods
```

#### Authentication Details Correct

```bash
$ ./proxyTest -proxy=192.168.56.121 -url=ingestion-lumberjack.condor.staging.axwaytest.net:453 -username=ubuntu -password=ubuntu
Attempting connection through Proxy without authentication
Connection to proxy made, connecting to url now
connection without authentication failed socks connect tcp 192.168.56.121:1080->ingestion-lumberjack.condor.staging.axwaytest.net:453: no acceptable authentication methods
Attempting connection through Proxy with username/password
Connection to proxy made, connecting to url now
connection with username/password authentication was successful
Attempting connection through Proxy with gss-api
Connection to proxy made, connecting to url now
connection with gss-api authentication failed socks connect tcp 192.168.56.121:1080->ingestion-lumberjack.condor.staging.axwaytest.net:453: no acceptable authentication methods
```

### GSS-API auth

This was not tested for success as it requires a lot of setup to enable in the test environment
