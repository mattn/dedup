# dedup

de-duplicate line from stdin

## Usage

```
Usage of dedup:
  -f string
    	storage file (default ".dedup")
  -k string
    	identify for the key (default "id")
```

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
* */1 * * * twty -json | dedup -k id_str -f /tmp/twty | jsonargs do-something
```

## Installation

```
$ go get github.com/mattn/dedup
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
