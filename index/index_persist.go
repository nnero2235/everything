package index

import (
	"bufio"
	"bytes"
	"everything/common"
	"everything/datastruct/linklist"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func persist(index *Index) error {
	defer common.PrintExecTime("Index Persist...")()

	_, err := os.Stat(PersistenceDir)
	if err != nil { //means dir not exists
		if err = os.Mkdir(PersistenceDir, os.ModeDir); err != nil {
			return err
		}
	}
	var buf bytes.Buffer
	buf.WriteByte(0xff)
	buf.WriteByte(0xee)
	buf.WriteByte(0xdd)
	buf.WriteByte(0xcc)
	buf.WriteString(index.UniqueName)
	buf.WriteByte('|')
	fmt.Fprintf(&buf, "%d", index.AllSize)
	buf.WriteByte('|')
	fmt.Fprintf(&buf, "%d", index.DiffSize)
	buf.WriteString("\r\n")

	for key, value := range index.indexMap {
		buf.WriteByte('"')
		buf.WriteString(key)
		buf.WriteByte('"')
		buf.WriteByte(' ')
		for i := 0; i < value.Size(); i++ {
			data, err := value.Get(i)
			if err != nil {
				return err
			}
			if i > 0 {
				buf.WriteString("<-")
			}
			itemInfo := data.(ItemInfo)
			buf.WriteByte('"')
			buf.WriteString(itemInfo.Name)
			buf.WriteByte('"')
			buf.WriteByte('|')
			buf.WriteByte('"')
			buf.WriteString(itemInfo.PrefixPath)
			buf.WriteByte('"')
			buf.WriteByte('|')
			buf.WriteByte('"')
			buf.WriteString(itemInfo.Suffix)
			buf.WriteByte('"')
			buf.WriteByte('|')
			fmt.Fprintf(&buf, "%t", itemInfo.IsDir)
		}
		buf.WriteString("\r\n")
	}
	err = ioutil.WriteFile(PersistenceDir+"\\"+index.UniqueName+PersistenceSubbfix, buf.Bytes(), os.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}

func load(fileName string) (*Index, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	index := &Index{indexMap: make(map[string]*linklist.LinkList)}

	buf := bufio.NewReader(f)
	data, err := buf.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	//verify protocol
	if data[0] != 0xff && data[1] != 0xee && data[2] != 0xdd && data[3] != 0xcc {
		return nil, fmt.Errorf("Not index File! Prefix is:%x", data[0:4])
	}
	//firstLine parse
	firstLine := string(data[4 : len(data)-2])
	sArr := strings.Split(firstLine, "|")
	if len(sArr) != 3 {
		return nil, fmt.Errorf("index File Broken at FirstLine:%s", firstLine)
	}
	index.UniqueName = sArr[0]
	index.AllSize, err = strconv.Atoi(sArr[1])
	if err != nil {
		return nil, err
	}
	index.DiffSize, err = strconv.Atoi(sArr[2])
	if err != nil {
		return nil, err
	}

	//data parse
	for i := 0; ; i++ {
		s, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return index, nil
			}
			return nil, err
		}

		s = s[:len(s)-2]
		kvIndex := strings.Index(s, " ")
		if kvIndex == -1 {
			return nil, fmt.Errorf("KV err: index File Broken at line:%d -> lineContent:%s", i, s)
		}
		name := s[:kvIndex]
		name = name[1 : len(name)-1]
		lk := strings.Split(s[kvIndex+1:], "<-")
		list := linklist.NewLinkList()

		for _, item := range lk {
			fileds := strings.Split(item, "|")
			if len(fileds) != 4 {
				return nil, fmt.Errorf("fileds err: index File Broken at line:%d -> item:%s", i, item)
			}
			iName := fileds[0]
			iPrefixPath := fileds[1]
			iSuffix := fileds[2]
			iIsDir, err := strconv.ParseBool(fileds[3])
			if err != nil {
				return nil, fmt.Errorf("boolParse err: index File Broken at line:%d err:%v", i, err)
			}
			list.Insert(0, ItemInfo{
				Name:       iName[1 : len(iName)-1],
				PrefixPath: iPrefixPath[1 : len(iPrefixPath)-1],
				Suffix:     iSuffix[1 : len(iSuffix)-1],
				IsDir:      iIsDir,
			})
		}

		index.indexMap[name] = list
	}

}
