# TFTP Server

A sample RFC 1350 server, implemented from scratch according to the [spec](https://tools.ietf.org/html/rfc1350) in Go.

## Running the server

Clone this repo and execute `./runServer.sh` from a POSIX-compatible shell. Superuser permissions are needed due to it listening on the default port 69.

## Sample use case

Using the default _tftp_ client available from most package managers:

```
echo "Hello, world!" >> hello.txt
tftp 127.0.0.1
mode octet
put hello.txt
get hello.txt
```

## Building the project

Dependencies: Go 1 (developed with Go 1.11.2 (linux/amd64))

```
cd cmd/tftpd && go build
sudo ./tftpd
```

## Request log

Read (RRQ), Write (WRQ) requests and unknown requests will be outputted to `requestLog.txt` in the current working directory.

## Tests

for the server are located in `lib/tftpserver_test.go` and are integration-style tests. When ran within Visual Studio Code, these will output the requestLog to `lib/`.

## Author

Callum Gavin