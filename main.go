package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func main() {
	var f, k string

	flag.StringVar(&f, "f", ".dedup", "storage file")
	flag.StringVar(&k, "k", "id", "identify for the key")
	flag.Parse()

	store, err := leveldb.OpenFile(f, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	defer store.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var v map[string]interface{}
		line := scanner.Text()
		err = json.Unmarshal([]byte(line), &v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
			continue
		}
		fv, ok := v[k]
		if !ok {
			fmt.Fprintf(os.Stderr, "%v: %q not found\n", os.Args[0], k)
			continue
		}

		_, err = store.Get([]byte(k), &opt.ReadOptions{DontFillCache: true})
		if err != nil {
			if err != leveldb.ErrNotFound {
				fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
				continue
			}
			err = store.Put([]byte(k), []byte(fmt.Sprint(fv)), &opt.WriteOptions{Sync: true})
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
				continue
			}
			fmt.Println(line)
		}
	}
}
