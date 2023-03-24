package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	checkFile()

	object := make(map[string]bool)

	object["lorem"] = true

	for k, v := range object {
		fmt.Println(fmt.Sprintf("%v %v", k, v))
	}

	list := make([]string, 0)

	list = append(list, "ipsum", "xzy", "abcde")

	sort.Strings(list)

	for _, v := range list {
		fmt.Println(v)
	}
}

func checkFile() {
	f, err := os.Open("C:\\Users\\ws-91\\Canvas\\some-api\\cmd\\coding\\scratch.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	f.Close()

	for _, line := range lines {
		fmt.Println(line)
	}
}
