# log-analyzer

Turn unstructured log files into a dashboard.

## Example

An unstructured log file named `demo.log`:

```log
[11/Dec/2023:11:01:28] 220.203.23.174 "GET /blog/home HTTP/1.1" 200 182 "Mozilla/5.0 Chrome/60.0.3112.113"
[11/Dec/2023:11:01:29] 89.238.65.53 "POST /new-user/ HTTP/1.1" 201 182 "Mozilla/5.0 (Linux; Android 13; SM-S901B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36"
[11/Dec/2023:11:01:29] 209.51.141.74 "GET /test HTTP/1.1" 404 182 "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"
[11/Dec/2023:11:01:32] 122.161.56.36 [error] request failed: error reading the headers
[11/Dec/2023:11:01:34] 74.6.8.121 "GET /api/data HTTP/1.1" 200 182 "python-requests/2.26.0"
```

In `config.json`, define some simple patterns featured in your log files using a set of token values to control the extraction. These tokens can have any identifying name, and will be grouped across all patterns and targeted for extraction from the log file.

- `patterns`: a list of string patterns found in your log file    
- `tokens`: all unique tokens used across your patterns, representing which values are important and should be targeted for extraction

```json
{
    "patterns": [
        "[timestamp] ip \"method endpoint http\" status bytes \"user_agent\"",
        "[timestamp] ip [error] message"
    ],
    "tokens": ["timestamp", "ip", "method", "endpoint", "http", "status", "bytes", "user_agent", "message"]
}
```

Run the log analyzer, providing the path to your log file.

```bash
> go build main.go
> ./main ./tests/data/logs/demo.log
```

The tokens are extracted from the log file and their data types inferred.

```text
line 1
        timestamp(time.Time): 2023-12-11 11:01:28 +0000 UTC
        ip(string): 220.203.23.174
        method(string): GET
        endpoint(string): /blog/home
        http(string): HTTP/1.1
        status(int): 200
        bytes(int): 182
        user_agent(string): Mozilla/5.0 Chrome/60.0.3112.113
line 2
        timestamp(time.Time): 2023-12-11 11:01:29 +0000 UTC
        ip(string): 89.238.65.53
        method(string): POST
        endpoint(string): /new-user/
        http(string): HTTP/1.1
        status(int): 201
        bytes(int): 182
        user_agent(string): Mozilla/5.0 (Linux; Android 13; SM-S901B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36
line 3
        timestamp(time.Time): 2023-12-11 11:01:29 +0000 UTC
        ip(string): 209.51.141.74
        method(string): GET
        endpoint(string): /test
        http(string): HTTP/1.1
        status(int): 404
        bytes(int): 182
        user_agent(string): Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36
line 4
        timestamp(time.Time): 2023-12-11 11:01:32 +0000 UTC
        ip(string): 122.161.56.36
        message(string): request failed: error reading the headers
line 5
        timestamp(time.Time): 2023-12-11 11:01:34 +0000 UTC
        ip(string): 74.6.8.121
        method(string): GET
        endpoint(string): /api/data
        http(string): HTTP/1.1
        status(int): 200
        bytes(int): 182
        user_agent(string): python-requests/2.26.0
```

Finally, your dashboard is generated. 

```text
Dashboard running at http://localhost:3000/
```

![dashboard](https://user-images.githubusercontent.com/41476809/279524886-acebc2a1-519b-42e5-a31e-043efb9d012c.png)

## Advanced Config

### Wildcard

Unimportant values within your log file that may still vary in value can be excluded from your dashboard by using wildcard tokens `*` or `_` within your pattern.

For example, if you don't want thread pool number or thread ID featured in your dashboard:

```log
[2023-10-25T16:24:31+00:00] pool-1-thread-2 INFO: getUserID() duration 14.29 ms
[2023-10-25T16:24:32+00:00] pool-1-thread-1 INFO: getUserID() duration 13.11 ms
[2023-10-25T16:24:37+00:00] pool-1-thread-2 INFO: getStatus() duration 3.87 ms
```

```json
{
    "patterns": [
        "[timestamp] pool-*-thread-* INFO: function() duration elapsed ms"
    ],
    "tokens": ["timestamp", "function", "elapsed"]
}
```

### Dependencies

Within the config, you can specify any tokens that are dependent upon other tokens. For example, elapsed running time may depend on the function name. Each dependency specified allows for deeper analysis with a richer dashboard that considers this relationship.

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

Some token values with a numeric value may be equivalent to other token values after performing a conversion to account for different units. In order to group these tokens together in your dashboard, your config needs to specify how to convert them from one to another. The example below states that `elapsed_ns` can be converted into `elapsed_ms` by multiplying by `0.001`, and `elapsed_s` can be converted into `elapsed_ms` by multiplying by 1000. With this config, the dashboard will convert and group all time recording values into milliseconds.

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

You can specify a path to the config file containing your patterns following the `-c` or `--config` flag. The config path defaults to `./config.json`

### Config Test

Once you have your patterns together, you can perform a test run by including the `-t` or `--test` flag and the extracted result of each line will be saved to json file in the current directory.

```bash
> ./main ./tests/data/logs/demo.log --test
```
