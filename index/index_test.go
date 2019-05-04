package index

import (
	"fmt"
	"testing"
)

var index = NewIndex("nnero")

func init() {
	index.Insert(ItemInfo{Name: "nnero", PrefixPath: "C:\\haha\\", Suffix: ".a"})
	index.Insert(ItemInfo{Name: "wahaha", PrefixPath: "C:\\haha\\hehe\\", Suffix: ""})
	index.Insert(ItemInfo{Name: "like", PrefixPath: "D:\\game\\ini\\", Suffix: ".jpg"})
	index.Insert(ItemInfo{Name: "like", PrefixPath: "F:\\code\\go\\", Suffix: ".png"})
	index.Insert(ItemInfo{Name: "like", PrefixPath: "C:\\haha\\wawa\\", Suffix: ".go"})
}

func TestIndexBasic(t *testing.T) {
	if index.AllSize != 5 {
		t.Errorf("AllSize should be: %d, but real is %d\n", 5, index.AllSize)
	}
	if index.DiffSize != 3 {
		t.Errorf("DiffSize should be: %d, but real is %d\n", 3, index.DiffSize)
	}

	items, err := index.Get("nnero")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	for _, i := range items {
		fmt.Printf("Name:%s\tPrefixPath:%s\n", i.Name, i.PrefixPath)
	}

	items, err = index.Get("like")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	for _, i := range items {
		fmt.Printf("Name:%s\tPrefixPath:%s\n", i.Name, i.PrefixPath)
	}
}

func TestPersist(t *testing.T) {
	ch := make(chan struct{})
	go Persist(index, ch)
	<-ch
}

func TestLoad(t *testing.T) {
	indexes, err := LoadDefault()
	if err != nil {
		t.Errorf("%v\n", err)
	}

	for _, in := range indexes {

		if in.AllSize != 5 {
			t.Errorf("AllSize should be: %d, but real is %d\n", 5, in.AllSize)
		}
		if in.DiffSize != 3 {
			t.Errorf("DiffSize should be: %d, but real is %d\n", 3, in.DiffSize)
		}

		items, err := in.Get("nnero")
		if err != nil {
			t.Errorf("%v\n", err)
		}

		for _, i := range items {
			fmt.Printf("Name:%s\tPrefixPath:%s\n", i.Name, i.PrefixPath)
		}

		items, err = in.Get("like")
		if err != nil {
			t.Errorf("%v\n", err)
		}

		for _, i := range items {
			fmt.Printf("Name:%s\tPrefixPath:%s\n", i.Name, i.PrefixPath)
		}
	}
}
