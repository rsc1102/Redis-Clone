# Redis Clone
The repository is a simple In-Memory Database similar to Redis, made using Go. It accepts commands via the redis-client.

## Installation
Clone the git repo using the command:
```bash
git clone
```

Move to the directory and start the server using the command:
```bash
cd redis-clone
go run .
```

Run the client in a different terminal:
```bash
redis-cli
```

## Supported Commands
| Command | Syntax | Example | Output |
| :--------: | :------- | :-------- | :-------: |
| PING | PING [message] | PING | "PONG" |
| SET | SET key value | SET myKey "Hello" | "OK" |
| GET | GET key | GET myKey | "Hello" |
| EXISTS | EXISTS key [key ...] | EXISTS myKey | (integer) 1 |
| DEL | DEL key [key ...] | DEL myKey | (integer) 1 |
| HSET | HSET key field value | HSET myhash field1 "Hello" | "OK" |
| HGET | HGET key field | HGET myhash field1 | "Hello" |
| HGETALL | HGETALL key | HGET myhash field1 | "Hello" |
| HEXISTS | HEXISTS key field  |  HEXISTS myhash field1 | (integer) 1 |
| HDEL | HDEL key field [field ...] | HDEL myhash field1 | (integer) 1 |

## Data Persistence
The data is stored in an AOF (Append only file). In this method, the Redis Clone records each command in the file as Redis serialization protocol (RESP). When a restart occurs, the Redis Cone reads all the RESP commands from the AOF file and executes them in memory. 
