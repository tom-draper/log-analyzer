# Log Analyzer

Turn unstructured log files into a dashboard in seconds.

![Log Analyzer](https://github.com/tom-draper/log-analyzer/assets/41476809/384162a6-58aa-4c9d-b22e-7544baba5b88)

## Example

A log file named `demo.log`:

```log
[11/Dec/2023:11:01:28] 220.203.23.174 "GET /blog/home HTTP/1.1" 200 182 "Mozilla/5.0 Chrome/60.0.3112.113"
[11/Dec/2023:11:01:29] 89.238.65.53 "POST /new-user/ HTTP/1.1" 201 182 "Mozilla/5.0 (Linux; Android 13; SM-S901B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36"
[11/Dec/2023:11:01:29] 209.51.141.74 "GET /test HTTP/1.1" 404 182 "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"
[11/Dec/2023:11:01:32] 122.161.56.36 [error] request failed: error reading the headers
[11/Dec/2023:11:01:34] 74.6.8.121 "GET /api/data HTTP/1.1" 200 182 "python-requests/2.26.0"
```

In `config/config.json`, define some simple patterns featured in your log files using a set of tokens to control which values should be extracted. These tokens can have any identifying name, and will be grouped across all patterns and targeted for extraction from the log file.

- `patterns`: a list of string patterns found in your log file    
- `tokens`: all unique tokens used across your patterns, representing the values that should be targeted for extraction

```json
{
    "patterns": [
        "[timestamp] ip_address \"method endpoint http\" status bytes \"user_agent\"",
        "[timestamp] ip_address [error] message"
    ],
    "tokens": ["timestamp", "ip_address", "method", "endpoint", "http", "status", "bytes", "user_agent", "message"]
}
```

Run the log analyzer, providing the path to your log file.

```bash
> go build main.go
> ./main ./tests/data/logs/demo.log
```

The tokens are extracted from the log file and their data types inferred. Any timestamps are parsed without needing to specify a format.

```text
line 1
        timestamp(time): 2023-12-11 11:01:28 +0000 UTC
        ip_address(ip): 220.203.23.174
        method(str): GET
        endpoint(str): /blog/home
        http(str): HTTP/1.1
        status(int): 200
        bytes(int): 182
        user_agent(str): Mozilla/5.0 Chrome/60.0.3112.113
line 2
        timestamp(time): 2023-12-11 11:01:29 +0000 UTC
        ip_address(ip): 89.238.65.53
        method(str): POST
        endpoint(str): /new-user/
        http(str): HTTP/1.1
        status(int): 201
        bytes(int): 182
        user_agent(str): Mozilla/5.0 (Linux; Android 13; SM-S901B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36
line 3
        timestamp(time): 2023-12-11 11:01:29 +0000 UTC
        ip_address(ip): 209.51.141.74
        method(str): GET
        endpoint(str): /test
        http(str): HTTP/1.1
        status(int): 404
        bytes(int): 182
        user_agent(str): Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36
line 4
        timestamp(time): 2023-12-11 11:01:32 +0000 UTC
        ip_address(ip): 122.161.56.36
        message(str): request failed: error reading the headers
line 5
        timestamp(time): 2023-12-11 11:01:34 +0000 UTC
        ip_address(ip): 74.6.8.121
        method(str): GET
        endpoint(str): /api/data
        http(str): HTTP/1.1
        status(int): 200
        bytes(int): 182
        user_agent(str): python-requests/2.26.0
```

Finally, your dashboard is generated. 

```text
Dashboard running at http://localhost:3000/
```

![Dashboard](https://github.com/tom-draper/log-analyzer/assets/41476809/6fd6cf34-f5f5-4148-9d40-e86288751c6b)

## Advanced Config

### Data Types

Data types are inferred by default, but it may not always be possible to correctly identify the intended type of a value. You can provide explicit data types to help the parser with one of the following prefixes on your token name: `int_`, `float_`, `str_`, `bool_`, `time_` and `ip_`.

If a token value fails to be converted into the explicit data type specified by its prefix, the value is excluded from the dashboard instead of reverting to the default string type.

Data type warnings can be found at the bottom of your dashboard, highlighting tokens with inconsistently inferred data types. This can be used to help you decide whether your patterns need improving or you could benefit from explicit data types. 

### Multi-line Patterns

Patterns that can match multiple lines of a log file can be created simply by including newline characters `\n` in your pattern. Although the reliability of this pattern may be variable if the log file is being written to by multiple threads.

```log
2022-01-11T12:15:06+00:00 - DEBUG thread1: success
2022-01-11T12:15:09+00:00 - DEBUG thread2: error
  --> critical error on thread2: index 44 out of range
2022-01-11T12:15:12+00:00 - DEBUG thread3: success
```

```json
{
    "patterns": [
        "timestamp - log_level thread_id: message",
        "timestamp - log_level thread_id: error\n  --> critical error on threadthread_id: error_message"
    ],
    "tokens": ["timestamp", "log_level", "thread_id", "message", "error_message"]
}
```

### Wildcards

Any unimportant yet variable values within your log file can be excluded from your dashboard by using the wildcard token `*` within your pattern.

For example, if you don't want thread pool number or thread ID featured in your dashboard:

```log
[2023-10-25T16:24:31+00:00] pool-1-thread-2 INFO: getUserID() duration 14.29 ms
[2023-10-25T16:24:32+00:00] pool-1-thread-1 INFO: getUserID() duration 13.11 ms
[2023-10-25T16:24:37+00:00] pool-1-thread-2 INFO: getStatus() duration 3.87 ms
```

```json
{
    "patterns": [
        "[<timestamp>] pool-*-thread-* INFO: <function>() duration <elapsed> ms"
    ],
    "tokens": ["<timestamp>", "<function>", "<elapsed>"]
}
```

### Dependencies

Within the config, you can specify any tokens that are dependent upon other tokens. For example, elapsed running time may depend on the function name. Each dependency specified allows for a richer dashboard that considers this relationship and explores the combinations of values.

```log
2023-10-25 20:44:07 - LOG: genUUID() duration 44.29 ms
2023-10-25 20:44:08 - LOG: createUser() duration 51.13 ms
2023-10-25 20:44:11 - LOG: recordDataAccess() duration 6.87 ms
```

```json
{
    "patterns": [
        "timestamp - LOG: function() duration elapsed ms"
    ],
    "tokens": ["timestamp", "function", "elapsed"],
    "dependencies": {
        "elapsed": ["function"]
    }
}
```

### Conversions

Some tokens with a numeric value may be equivalent to other tokens after performing a simple conversion to account for different units. In order to merge these tokens together in your dashboard, your config needs to specify how to convert them from one to another. The example below describes that `elapsed_ns` and `elapsed_s` can both be converted into `elapsed_ms` by multiplying by 0.001 and 1000 respectively. With this config, the dashboard will convert and group all elapsed time values into milliseconds and display them under the `elapsed_ms` token.

```log
[November 26, 2017 at 7:25p PST] LOG: function 'createAccount' took 51.13 ms
[November 26, 2017 at 7:25p PST] LOG: function 'encryptData' took 1.29 s
[November 26, 2017 at 7:25p PST] LOG: function 'checkOnline' took 904.87 ns
```

```json
{
    "patterns": [
        "[timestamp] INFO: function 'function_name' took elapsed_ms ms"
        "[timestamp] INFO: function 'function_name' took elapsed_ns ns"
        "[timestamp] INFO: function 'function_name' took elapsed_s s"
    ],
    "tokens": ["timestamp", "function_name", "elapsed_ms", "elapsed_ns", "elapsed_s"],
    "conversions": {
        "elapsed_ns": {
            "token": "elapsed_ms",
            "multiplier": 0.001
        },
        "elapsed_s": {
            "token": "elapsed_ms",
            "multiplier": 1000
        }
    }
}
```

### Config Path

You can specify a path to the config file containing your patterns following the `-c` or `--config` flag. The config path defaults to `config/config.json`

```bash
> ./main -c ./tests/data/configs/config.json ./tests/data/logs/demo.log
```

### Config Test

Once you have your patterns together, you can perform a test run by including the `-t` or `--test` flag and the extracted result of each line will be written to json file in the current directory.

```bash
> ./main ./tests/data/logs/demo.log --test
```

## Additional Tools

### Pattern Identification

With larger, more complex log files, identifying the patterns to extract may be challenging. To help create these patterns, the similarity tool can be used to group the lines of a log file by similarity.

```text
go get github.com/tom-draper/log-analyzer/pkg/similarity
```

```go
package main

import (
    "os"
    "github.com/tom-draper/log-analyzer/pkg/similarity"
)

func main() {
    body, err := os.ReadFile("./tests/data/logs/demo.log")
    if err != nil {
        panic(err)
    }
    logtext := string(body)

    groups := similarity.FindGroups(logtext)
}
```

### Parser

The project can also be used without the dashboard, as a parser for your log files to extract values.

```text
go get github.com/tom-draper/log-analyzer/pkg/parse
```

```go
package main

import "github.com/tom-draper/log-analyzer/pkg/parse"

func main() {
    config, err := parse.LoadConfig("./config/config.json")
    if err != nil {
        panic(err)
    }

    extraction, err := parse.ParseFile("./tests/data/logs/demo.log", &config)
    if err != nil {
        panic(err)
    }

    parse.DisplayLines(extraction)
}
```

## Contributions

Contributions, issues and feature requests are welcome.

- Fork it (https://github.com/tom-draper/log-analyzer)
- Create your feature branch (`git checkout -b my-new-feature`)
- Commit your changes (`git commit -am 'Add some feature'`)
- Push to the branch (`git push origin my-new-feature`)
- Create a new Pull Request
