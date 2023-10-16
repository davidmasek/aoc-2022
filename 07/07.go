package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Node interface {
	Size() int
	CalculateSize()
	Name() string
}

type File struct {
	name   string
	parent *Directory
	size   int
}

type Directory struct {
	parent   *Directory
	name     string
	size     int
	children []Node
}

func (t *File) Name() string {
	return t.name
}

func (d *Directory) Name() string {
	return d.name
}

func (d *Directory) addChild(c Node) {
	d.children = append(d.children, c)
}

func (d *Directory) CalculateSize() {
	s := 0
	for _, c := range d.children {
		c.CalculateSize()
		s += c.Size()
	}
	d.size = s
}

func (d *Directory) GetChild(name string) Node {
	for _, c := range d.children {
		if c.Name() == name {
			return c
		}
	}
	panic("not found")
}

func (f *File) CalculateSize() {
	// noop
}

func (d *Directory) Size() int {
	return d.size
}

func (f *File) Size() int {
	return f.size
}

func (d *Directory) Iterate() []*Directory {
	children := make([]*Directory, 0)
	children = append(children, d)
	for _, c := range d.children {
		if d, ok := c.(*Directory); ok {
			children = append(children, d.Iterate()...)
		}
	}
	return children
}

func startswith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func main() {
	fmt.Println(len(os.Args), os.Args)
	fmt.Println("-----")
	if len(os.Args) != 2 {
		fmt.Println("Usage: XY.go XY.in")
		os.Exit(1)
	}
	var filename string
	test := true
	if test {
		filename = strings.Replace(os.Args[1], ".ex", ".in", 1)
	} else {
		filename = os.Args[1]
	}
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	root := Directory{name: "/"}
	c := &root
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		// assume start here and skip
		if line == "$ cd /" {
			// fmt.Println("Cwd: /")
			continue
		}
		// handle this implicitly below
		if line == "$ ls" {
			continue
		}
		if startswith(line, "$ cd ") {
			next_name := line[len("$ cd "):]
			if next_name == ".." {
				c = c.parent
			} else {
				c = c.GetChild(next_name).(*Directory)
			}
			continue
		}
		if startswith(line, "dir ") {
			dir_name := line[len("dir "):]
			c.addChild(&Directory{name: dir_name, parent: c})
			// fmt.Printf("added %s to %s\n", dir_name, c.name)
			continue
		}
		var size int
		var name string
		_, err = fmt.Sscanf(line, "%d %s", &size, &name)
		check(err)
		c.addChild(&File{name: name, size: size, parent: c})
		// fmt.Printf("added %s to %s\n", name, c.name)
	}
	fmt.Println("-----")
	root.CalculateSize()
	fmt.Printf("%#v\n", root)
	fmt.Println("-----")
	sum := 0
	for _, c := range root.Iterate() {
		if c.size <= 100000 {
			sum += c.size
		}
		fmt.Printf("%s %d\n", c.name, c.size)
	}
	fmt.Println("-----")
	fmt.Println("A:", sum)
	fmt.Println("-----")
	const TOTAL_SIZE = 70000000
	const REQUIRED_FREE_SIZE = 30000000
	totalUsed := root.size
	totalFree := TOTAL_SIZE - totalUsed
	neededFree := REQUIRED_FREE_SIZE - totalFree
	fmt.Println(totalUsed, totalFree, neededFree)
	minValue := TOTAL_SIZE
	for _, c := range root.Iterate() {
		if c.size < minValue && c.size >= neededFree {
			minValue = c.size
		}
	}
	fmt.Println("-----")
	fmt.Println("B:", minValue)
}
