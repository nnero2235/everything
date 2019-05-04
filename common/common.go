package common

import (
	"fmt"
	"os"
	"strings"
	"time"
)

//return name and suffix
func GetNameAndSuffix(fileName string) (string, string) {
	i := strings.LastIndex(fileName, ".")
	if i == -1 {
		return fileName, ""
	} else {
		suffix := fileName[i:]
		return fileName[:i], suffix
	}
}

//ugly way to get windows drives. But simple enough
func GetWindowsDrives() []string {
	var r []string
	for _, drive := range "abcdefghijklmnopqrstovwxyz" {
		_, err := os.Open(string(drive) + ":\\")
		if err == nil {
			r = append(r, string(drive)+":\\")
		}
	}
	return r
}

//use defer func() to print method's running time
func PrintExecTime(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s running: %.3fs\n", name, time.Since(start).Seconds())
	}
}
