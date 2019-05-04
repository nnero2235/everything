package linklist

import (
	"bytes"
	"fmt"
)

type LinkList struct {
	root *Node
	tail *Node
	size int
}

type Node struct {
	data interface{}
	pre  *Node
	next *Node
}

func NewLinkList() *LinkList {
	return &LinkList{root: nil, size: 0}
}

func (list *LinkList) Size() int {
	return list.size
}

func (list *LinkList) Add(data interface{}) error {
	return list.Insert(list.size-1, data)
}

func (list *LinkList) Get(index int) (interface{}, error) {
	if list.size == 0 || list.size <= index {
		return nil, fmt.Errorf("Out of Index: size is %d, index is %d", list.size, index)
	}
	mid := list.size / 2
	if index <= mid {
		node := list.root
		for i := 0; i <= mid; i++ {
			if i == index {
				return node.data, nil
			}
			node = node.next
		}
		panic("improssible Get reach: aesc")
	} else {
		node := list.tail
		for i := list.size - 1; i > mid; i-- {
			if i == index {
				return node.data, nil
			}
			node = node.pre
		}
		panic("Improssible Get reach: desc")
	}
}

func (list *LinkList) Insert(index int, data interface{}) error {
	if data == nil {
		return fmt.Errorf("Insert: data can't be nil")
	}
	if list.size == 0 {
		list.root = &Node{data: data, pre: nil, next: nil}
		list.tail = list.root
		list.size++
		return nil
	}

	newNode := &Node{data: data, pre: nil, next: nil}

	if index == 0 {
		list.root.pre = newNode
		newNode.next = list.root
		list.root = newNode
		list.size++
		return nil
	} else if index == list.size-1 {
		list.tail.next = newNode
		newNode.pre = list.tail
		list.tail = newNode
		list.size++
		return nil
	}

	mid := list.size / 2
	if index <= mid {
		node := list.root
		for i := 0; i <= mid; i++ {
			if i == index {
				node.pre.next = newNode
				newNode.pre = node.pre
				newNode.next = node
				node.pre = newNode
				list.size++
				return nil
			}
			node = node.next
		}
		panic("improssible Get reach: aesc")
	} else {
		node := list.tail
		for i := list.size - 1; i > mid; i-- {
			if i == index {
				node.pre.next = newNode
				newNode.pre = node.pre
				newNode.next = node
				node.pre = newNode
				list.size++
				return nil
			}
			node = node.pre
		}
		panic("Improssible Get reach: desc")
	}
	return nil
}

func (list *LinkList) Remove(index int) error {
	if list.size == 0 {
		return nil
	}

	if index == 0 {
		list.root = list.root.next
		list.root.pre = nil
		list.size--
		return nil
	} else if index == list.size-1 {
		list.tail = list.tail.pre
		list.tail.next = nil
		list.size--
		return nil
	}

	mid := list.size / 2
	if index <= mid {
		node := list.root
		for i := 0; i <= mid; i++ {
			if i == index {
				node.pre.next = node.next
				node.next.pre = node.pre
				list.size--
				return nil
			}
			node = node.next
		}
		panic("improssible Get reach: aesc")
	} else {
		node := list.tail
		for i := list.size - 1; i > mid; i-- {
			if i == index {
				node.pre.next = node.next
				node.next.pre = node.pre
				list.size--
				return nil
			}
			node = node.pre
		}
		panic("Improssible Get reach: desc")
	}
	return nil
}

func (list *LinkList) PrintList() {
	if list.size == 0 {
		fmt.Println("[]")
		return
	}
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := 0; i < list.size; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		if v, err := list.Get(i); err != nil {
			fmt.Fprintf(&buf, "Error:%v\n", err)
			return
		} else {
			fmt.Fprintf(&buf, "%v", v)
		}
	}
	buf.WriteByte(']')
	fmt.Println(buf.String())
}
