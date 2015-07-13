testspy
---

A silly program that scans a directory for "test" files and
runs `go test ./...` when they are changed.

Work in progress, don't expect much

Install with

```bash
go install github.com/julien/testspy
```

If your [GOPATH](https://golang.org/doc/code.html#GOPATH) i
is configured properly, you should be able to do this

```shell
testspy
```

You can optionally specify a path and a coverage file name if you want to,

```shell
testspy -path=/home/foo/code/superduper -coverfile=bananas.out
```

