# log-analyzer
Turn unstructured log files into dashboards.

## Example

An unstructured log file named `demo.log`:

```log
[2007-09-01 16:44:49.244 ADT] 192.168.2.10:ossecdb LOG:  duration: 4.550 ms  statement: SELECT id FROM location WHERE name = 'enigma->/var/log/messages' AND server_id = '1'
[2007-09-01 16:44:49.251 ADT] 192.168.2.10:ossecdb LOG:  duration: 5.252 ms  statement: INSERT INTO location(server_id, name) VALUES ('1', 'enigma->/var/log/messages')
[2007-09-01 16:44:49.252 ADT] 192.168.2.10:ossecdb LOG:  duration: 0.016 ms  statement: SELECT id FROM location WHERE name = 'enigma->/var/log/messages' AND server_id = '1'
[2007-09-27 11:02:51.611 ADT] 192.168.2.10:ossecdb LOG:  statement: INSERT INTO alert(id,server_id,rule_id,timestamp,location_id,src_ip) VALUES ('3577', '1', '50503','1190916566', '140', '0')
```

In `config.json`, build some simple patterns featured in your log files using a set of token values. These tokens can have any identifying name, and will be grouped and targeted for extraction from the log file.

- `patterns`: a list of string patterns found in your log file    
- `tokens`: all unique tokens used across your patterns, representing which values are important and should be targeted for extraction

```json
{
    "patterns": [
        "[timestamp] ip:database type:  statement: query",
        "[timestamp] ip:database type:  duration: elapsed ms  statement: query"
    ],
    "tokens": ["timestamp", "ip", "database", "type", "query", "elapsed"],
}
```

Run the log analyzer, providing the path to your log file.

```bash
> go build main.go
> ./main ./demo.log
```

The tokens are extracted from the log file and their data types inferred.

```text
line 0
        ip(string): 192.168.2.10
        database(string): ossecdb
        type(string): LOG
        elapsed(float64): 4.55
        query(string): SELECT id FROM location WHERE name = 'enigma->/var/log/messages' AND server_id = '1' 
        timestamp(time.Time): 2007-09-01 16:44:49.244 +0000 UTC
line 1
        elapsed(float64): 5.252
        query(string): INSERT INTO location(server_id, name) VALUES ('1', 'enigma->/var/log/messages')      
        timestamp(time.Time): 2007-09-01 16:44:49.251 +0000 UTC
        ip(string): 192.168.2.10
        database(string): ossecdb
        type(string): LOG
line 2
        type(string): LOG
        elapsed(float64): 0.016
        query(string): SELECT id FROM location WHERE name = 'enigma->/var/log/messages' AND server_id = '1' 
        timestamp(time.Time): 2007-09-01 16:44:49.252 +0000 UTC
        ip(string): 192.168.2.10
        database(string): ossecdb
line 3
        database(string): ossecdb
        type(string): LOG
        query(string): INSERT INTO alert(id,server_id,rule_id,timestamp,location_id,src_ip) VALUES ('3577', '1', '50503','1190916566', '140', '0')
        timestamp(time.Time): 2007-09-27 11:02:51.611 +0000 UTC
        ip(string): 192.168.2.10
```

Finally, your dashboard is generated. 

```text
Dashboard running at http://localhost:3000/
```
