package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const name = "dedup"

const version = "0.0.1"

var revision = "HEAD"

func run() int {
	var f, k string
	var dump bool
	var showVersion bool

	flag.StringVar(&f, "f", ".dedup", "storage file")
	flag.StringVar(&k, "k", "id", "identify for the key")
	flag.BoolVar(&dump, "dump", false, "dump stored keys")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	store, err := leveldb.OpenFile(f, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
		return 1
	}
	defer store.Close()

	if dump {
		snapshot, err := store.GetSnapshot()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: %q not found\n", os.Args[0], k)
			return 1
		}
		it := snapshot.NewIterator(nil, nil)
		for it.Next() {
			fmt.Println(string(it.Key()))
		}
		return 0
	}

	scanner := bufio.NewScanner(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for scanner.Scan() {
		var v map[string]interface{}
		text := scanner.Text()
		err = json.Unmarshal([]byte(text), &v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
			continue
		}
		vk, ok := v[k]
		if !ok {
			fmt.Fprintf(os.Stderr, "%v: %q not found\n", os.Args[0], k)
			continue
		}
		bk := []byte(fmt.Sprint(vk))

		_, err = store.Get(bk, &opt.ReadOptions{DontFillCache: true})
		if err != nil {
			if err != leveldb.ErrNotFound {
				fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
				continue
			}
			err = store.Put(bk, []byte(fmt.Sprint(time.Now().Unix())), &opt.WriteOptions{Sync: true})
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
				continue
			}
			enc.Encode(v)
		}
	}
	return 0
}

func main() {
	os.Exit(run())
}
