package main

import (
	"everything/common"
	"everything/index"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

var flagDir = flag.String("d", "", "find dirent")
var flagUpdate = flag.Bool("u", false, "update indexes")

func main() {
	defer common.PrintExecTime("main")()

	flag.Parse()
	targets := flag.Args()

	if len(targets) != 1 {
		log.Fatalf("Only support one file search! but given %d\n", len(targets))
	}

	var indexes []*index.Index
	if *flagUpdate {
		indexes = createIndex()
	} else {
		oldIndexes, err := getIndex()
		if err != nil {
			fmt.Printf("%v\n", err)
			oldIndexes = createIndex()
		}
		indexes = oldIndexes
	}

	items, err := searchIndex(targets[0], indexes, *flagDir)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	printResult(items)

}

func printResult(items []index.ItemInfo) {
	for _, item := range items {
		var fType string
		if item.IsDir {
			fType = "dir"
		} else {
			fType = "file"
		}
		fmt.Printf("Find (%4s): %s\\%s%s\n", fType, item.PrefixPath, item.Name, item.Suffix)
	}
}

func searchIndex(target string, indexes []*index.Index, dir string) ([]index.ItemInfo, error) {
	if len(indexes) == 0 {
		return nil, fmt.Errorf("No indexes.Something went wrong!")
	} else {
		for _, in := range indexes {
			name, suffix := common.GetNameAndSuffix(target)
			items, err := in.Get(name)
			if err != nil {
				return nil, err
			}
			if len(items) == 0 {
				return nil, fmt.Errorf("Nothing found in indexes!")
			} else {
				var result []index.ItemInfo
				for _, item := range items {
					if dir != "" && !strings.HasPrefix(item.PrefixPath, dir) {
						continue
					}
					if suffix != "" && item.Suffix != suffix {
						continue
					}
					result = append(result, item)
				}
				if len(result) == 0 {
					return nil, fmt.Errorf("Nothing found in indexes!")
				}
				return result, nil
			}
		}
		panic("searchIndex: improssible reach")
	}
}

func getIndex() ([]*index.Index, error) {
	items, err := index.LoadDefault()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func createIndex() []*index.Index {
	defer common.PrintExecTime("createIndex")()

	fmt.Printf("Update indexes: starting...\n")

	var items []*index.Index
	currentFiles := 0
	files := make(chan int)
	indexes := make(chan []*index.Index)
	persistCh := make(chan struct{})
	startTime := time.Now()

	var tick <-chan time.Time
	tick = time.Tick(500 * time.Millisecond)

	go index.CreateIndex(files, indexes, persistCh)
loop:
	for {
		select {
		case <-tick:
			fmt.Printf("indexing... files:%d \t running:%.3fs\n", currentFiles, time.Since(startTime).Seconds())
		case n := <-files:
			currentFiles += n
		case data := <-indexes:
			items = data
			break loop
		}
	}
	//block until persit commplete
	<-persistCh
	fmt.Printf("indexes create successful! Totally created: %d indexes\n", len(items))
	return items
}
