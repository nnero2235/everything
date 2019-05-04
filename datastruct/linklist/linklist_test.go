package linklist

import (
	"fmt"
	"testing"
)

func TestLinkList(t *testing.T) {
	list := NewLinkList()
	if err := list.Add(5); err != nil {
		t.Errorf("%v\n", err)
	}
	if err := list.Add(10); err != nil {
		t.Errorf("%v\n", err)
	}
	if err := list.Add(11); err != nil {
		t.Errorf("%v\n", err)
	}
	if err := list.Add(110); err != nil {
		t.Errorf("%v\n", err)
	}
	if err := list.Add(99); err != nil {
		t.Errorf("%v\n", err)
	}
	v, err := list.Get(4)
	if err != nil {
		t.Errorf("%v\n", err)
	}
	f := v.(int)
	a := 10
	fmt.Printf("%d\n", f+a)
	list.PrintList()

	list.Insert(3, -18)
	list.PrintList()

	list.Insert(4, "wahaha")
	list.PrintList()

	list.Remove(4)
	list.PrintList()
	fmt.Println(list.Size())
}
