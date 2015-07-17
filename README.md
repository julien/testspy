testspy
---

A silly program that scans a directory for "test" files and
runs `go test -cover ./...` when they're changed.

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

You can optionally user the -path flag to specify the directory
you want to watch, but it defaults to the current working directory.

```shell
testspy -path=/home/foo/code/superduper
```

