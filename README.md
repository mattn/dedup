# dedup

de-duplicate line from stdin

## Usage

```
$ cat foo.json
{"key": 123, value: "foo1"}
{"key": 124, value: "foo2"}

$ cat foo.json | dedup -k key
{"key": 123, value: "foo1"}
{"key": 124, value: "foo2"}

$ cat foo.json | dedup -k key

$ cat bar.json
{"key": 123, value: "foo3"}
{"key": 125, value: "foo4"}

$ cat bar.json | dedup -k key
{"key": 125, value: "foo4"}
```

## Tutorial

Do something for the twitter statuses getting from crontab.

```
* */1 * * * twty | dedup -k id_str | do-something
```

## Installation

```
$ go get github.com/mattn/dedup
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
