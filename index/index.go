package index

// index for items

import (
	"everything/common"
	"everything/datastruct/linklist"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

const (
	PersistenceDir     = "data"
	PersistenceSubbfix = ".data"
)

type Index struct {
	UniqueName string //diff from other Index Name for Global unique
	indexMap   map[string]*linklist.LinkList
	DiffSize   int
	AllSize    int
}

type ItemInfo struct {
	Name       string
	PrefixPath string
	IsDir      bool
	Suffix     string
}

func walkDir(dir string, n *sync.WaitGroup, items chan<- ItemInfo) {
	defer n.Done()
	for _, f := range dirents(dir) {
		name := f.Name()
		dotIndex := strings.LastIndex(name, ".")
		suffix := ""
		if dotIndex != -1 {
			suffix = name[dotIndex:]
			name = name[:dotIndex]
		}
		item := ItemInfo{Name: name, PrefixPath: dir, Suffix: suffix}
		if f.IsDir() {
			item.IsDir = true
			n.Add(1)
			go walkDir(dir+"\\"+f.Name(), n, items)
		}
		items <- item
	}

}

var sema = make(chan struct{}, 50)

func dirents(d string) []os.FileInfo {
	sema <- struct{}{}
	defer func() {
		<-sema
	}()
	files, err := ioutil.ReadDir(d)
	if err != nil {
		if !os.IsPermission(err) {
			fmt.Fprintf(os.Stderr, "ReadDir: %v\n", err)
		}
		return nil
	}
	return files
}

func CreateIndex(files chan<- int, indexes chan<- []*Index, persistCh chan<- struct{}) {
	items := make(chan ItemInfo)
	var n sync.WaitGroup

	dirs := common.GetWindowsDrives()
	for _, dir := range dirs {
		n.Add(1)
		go walkDir(dir, &n, items)
	}

	go func() {
		n.Wait()
		close(items)
	}()

	index := NewIndex("Files")
loop:
	for {
		select {
		case item, ok := <-items:
			if !ok {
				break loop
			}
			index.Insert(item)
			files <- 1
		}
	}
	//async to save
	go Persist(index, persistCh)

	var data []*Index
	data = append(data, index)
	indexes <- data
}

func NewIndex(name string) *Index {
	return &Index{indexMap: make(map[string]*linklist.LinkList), UniqueName: name}
}

func (index *Index) Insert(item ItemInfo) error {
	list, ok := index.indexMap[item.Name]
	if !ok {
		list = linklist.NewLinkList()
		index.indexMap[item.Name] = list
		index.DiffSize++
	}
	if err := list.Add(item); err != nil {
		return err
	}
	index.AllSize++
	return nil
}

func (index *Index) Get(name string) ([]ItemInfo, error) {
	list, ok := index.indexMap[name]
	if !ok {
		return nil, nil
	}
	var items []ItemInfo
	for i := 0; i < list.Size(); i++ {
		v, err := list.Get(i)
		if err != nil {
			return nil, err
		}
		itemInfo := v.(ItemInfo)
		items = append(items, itemInfo)
	}
	return items, nil
}

func Persist(index *Index, ch chan<- struct{}) {
	err := persist(index)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Persist: %v\n", err)
	}
	fmt.Println()
	ch <- struct{}{}
}

func LoadDefault() ([]*Index, error) {
	fileInfos, err := ioutil.ReadDir(PersistenceDir + "\\")
	if err != nil {
		return nil, err
	}
	if len(fileInfos) == 0 {
		return nil, fmt.Errorf("No indexes found in \"%s\"", PersistenceDir+"\\")
	}
	var indexes []*Index
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			index, err := load(PersistenceDir + "\\" + fileInfo.Name())
			if err != nil {
				return nil, err
			}
			indexes = append(indexes, index)
		}
	}
	return indexes, nil
}
