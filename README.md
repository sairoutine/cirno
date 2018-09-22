# Cirno

Cirno is a stand alone application to generate unique ID.

## Example

```
$ telnet localhost 11212
Trying ::1...
Connected to localhost.
Escape character is '^]'.
GET new
VALUE new 0 20
bejajfghs7icvi3okm70
END
```

## Installation

Build from source code.
```
$ go get github.com/sairoutine/cirno
$ cd $GOPATH/github.com/sairoutine/cirno
go install
```

## Usage

```
$ cd $GOPATH/github.com/sairoutine/cirno/cmd/cirno_server
go build
./cirno_server -port=7238
./cirno_server -sock=/path/to/unix-domain.sock
```

## Protocol

Cirno uses the protocol which is compatible with memcached.

### API

#### GET, GETS

```
GET id1 id2
VALUE id1 0 20
bejajfghs7icvi3okm70
VALUE id2 0 20
bejajfghs7icvi3okm71
END
```

VALUE(s) are unique IDs.

#### VERSION

Returns a version of Cirno.

```
VERSION 1.0.0
```

#### QUIT

Disconnect the established connection.

## Algorithm

Cirno uses xid algorithm to generate ID.
Xid is alod used also Mongo Object ID algorithm.

## Commandline Options

### -port

Optional.
Port number used for connection.
Default value is `11212`.

### -sock

Optional.
Path of unix doamin socket.

### -timeout

Optional.
Connection idle timeout in seconds.
`0` means infinite.
Default value is `5`.

## Test
```
# benchmark
go test -bench . -benchmem -benchtime 30s -v ./cmd/cirno_server
```

## Licence

[MIT](https://github.com/sairoutine/cirno/blob/master/LICENSE)


