package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func run() int {
	var f, k string
	var dump bool

	flag.StringVar(&f, "f", ".dedup", "storage file")
	flag.StringVar(&k, "k", "id", "identify for the key")
	flag.BoolVar(&dump, "dump", false, "dump stored keys")
	flag.Parse()

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

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for {
		var v map[string]interface{}
		err = dec.Decode(&v)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
			return 1
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
