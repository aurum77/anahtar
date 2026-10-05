# Anahtar

A Redis-compatible in-memory key-value store written in Go, built from scratch as a learning project.

Anahtar started as a simple TCP key-value server and is gradually evolving toward a more complete Redis-like server, including RESP protocol support, key expiration, persistence, and performance testing.

## Features

### Currently implemented

* TCP server
* Concurrent client connections
* `SET`, `GET`, `DEL`, and `PING` commands
* Thread-safe in-memory storage
* Key expiration with:

  * `EX` — seconds
  * `PX` — milliseconds
* Lazy expiration
* Background expired-key cleanup
* RESP request parsing
* Multiple commands over a persistent TCP connection
* RESP response encoding
* Compatible with `redis-cli` for basic commands

### Planned

* Complete RESP protocol handling
* Additional Redis data types
* Append-only persistence / WAL
* Database recovery on startup
* Concurrent load testing
* Performance profiling and optimization
* Possible store sharding for improved concurrency

## Architecture

```text
                    TCP Connection
                          │
                          ▼
                    RESP Reader
                          │
                          ▼
                    []string command
                          │
                          ▼
                      Handler
                          │
                          ▼
                       Store
                          │
                          ▼
                    Handler Result
                          │
                          ▼
                    RESP Writer
                          │
                          ▼
                    TCP Connection
```

The storage layer is intentionally independent of the RESP protocol. The handler deals with command semantics, while the RESP package is responsible for protocol encoding and decoding.

## Supported Commands

### PING

```text
PING
```

Returns:

```text
PONG
```

### SET

```text
SET key value
```

Set a key without expiration.

Expiration can also be specified:

```text
SET key value EX 10
SET key value PX 5000
```

### GET

```text
GET key
```

Returns the stored value if the key exists.

### DEL

```text
DEL key
```

Deletes a key.

## Example

Anahtar currently runs on TCP port `7777`.

Using `redis-cli`:

```bash
redis-cli -p 7777
```

Then:

```text
127.0.0.1:7777> SET foo bar
OK

127.0.0.1:7777> GET foo
"bar"

127.0.0.1:7777> DEL foo
OK

127.0.0.1:7777> GET foo
(nil)
```

Multiple commands can be sent over the same TCP connection.

## RESP

Anahtar uses the Redis Serialization Protocol (RESP) for communication.

For example:

```text
SET foo bar
```

is represented as:

```text
*3\r\n
$3\r\n
SET\r\n
$3\r\n
foo\r\n
$3\r\n
bar\r\n
```

Responses are encoded using RESP types such as:

```text
+OK\r\n
```

for Simple Strings,

```text
$3\r\nbar\r\n
```

for Bulk Strings,

```text
-ERR message\r\n
```

for errors, and:

```text
$-1\r\n
```

for a null Bulk String.

## Expiration

Keys can have an optional expiration time.

Example:

```text
SET foo bar EX 10
```

The key will expire after 10 seconds.

Anahtar uses both:

* lazy expiration when a key is accessed
* a background sweeper for removing expired entries

The store is protected by a read/write mutex to allow concurrent access safely.

## Project Structure

```text
anahtar/
├── cmd/
├── internal/
│   ├── handler/
│   ├── resp/
│   ├── server/
│   └── store/
├── go.mod
└── README.md
```

The exact structure may evolve as the project grows.

## Goals

The goal of Anahtar is not simply to reproduce Redis. It is primarily a learning project focused on understanding how a networked data store works internally.

The project is being developed incrementally, with each stage introducing a new systems concept:

1. TCP networking
2. Concurrent connections
3. In-memory data structures
4. Synchronization and race conditions
5. Key expiration
6. Network protocol design
7. RESP parsing and encoding
8. Persistence
9. Load testing
10. Performance optimization

## Roadmap

* [x] TCP server
* [x] Basic command handling
* [x] Thread-safe key-value store
* [x] `SET` / `GET` / `DEL`
* [x] Key expiration
* [x] Background expiration cleanup
* [x] RESP request parsing
* [x] Persistent TCP command loop
* [x] RESP response encoding
* [ ] Clean RESP-independent handler architecture
* [ ] Additional Redis data type
* [ ] Append-only persistence
* [ ] Startup recovery
* [ ] Load testing
* [ ] Performance profiling
* [ ] Store optimization / sharding

## Why "Anahtar"?

**Anahtar** means **"key"** in Turkish — fitting for a key-value store.

