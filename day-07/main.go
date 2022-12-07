package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	maxSize           int64 = 100000
	totalDiskSpace    int64 = 70000000
	requiredDiskSpace int64 = 30000000
)

type fileSystem struct {
	root *directory
	cwd  *directory
}

func (fs *fileSystem) addDir(name string) {
	newDir := directory{
		name:     strings.TrimSpace(name),
		parent:   fs.cwd,
		children: make(map[string]diskEntity),
	}
	fs.cwd.children[name] = &newDir
}

func (fs *fileSystem) addFile(size int64, name string) {
	newFile := file{
		name: strings.TrimSpace(name),
		size: size,
	}
	fs.cwd.children[name] = &newFile
}

func (fs *fileSystem) cd(dir string) {
	if dir == "/" {
		fs.cwd = fs.root
		return
	}
	if dir == ".." {
		fs.cwd = fs.cwd.parent
		return
	}

	diskDir := fs.cwd.children[dir]
	fs.cwd = diskDir.(*directory)
}

func (fs *fileSystem) ls(list []string) {
	for _, entry := range list {
		words := strings.Split(entry, " ")
		switch words[0] {
		case "dir":
			fs.addDir(words[1])
		default:
			size, _ := strconv.ParseInt(words[0], 10, 64)
			fs.addFile(size, words[1])
		}
	}
}

func (fs fileSystem) getSizeSum() int64 {
	return fs.root.getInnerSizeSumMax()
}

func (fs fileSystem) getSmallestEligibleDirectory() int64 {
	used := fs.root.Size()
	remaining := totalDiskSpace - used
	needed := requiredDiskSpace - remaining

	return fs.root.findClosestNeededSize(needed, used)
}

func newRootFileSystem() *fileSystem {
	rootDir := directory{
		name:     "/",
		parent:   nil,
		children: make(map[string]diskEntity),
	}

	return &fileSystem{
		root: &rootDir,
		cwd:  &rootDir,
	}
}

type diskEntity interface {
	Size() int64
}

type directory struct {
	name     string
	parent   *directory
	children map[string]diskEntity
}

func (d directory) Size() int64 {
	var size int64 = 0
	for _, c := range d.children {
		size += c.Size()
	}
	return size
}

func (d directory) getInnerSizeSumMax() int64 {
	var total int64 = 0

	size := d.Size()
	if size <= maxSize {
		total += size
	}

	for _, entity := range d.children {
		if dir, ok := entity.(*directory); ok {
			total += dir.getInnerSizeSumMax()
		}
	}

	return total
}

func (d directory) findClosestNeededSize(needed, closest int64) int64 {
	size := d.Size()
	if size > needed && size < closest {
		closest = size
	}
	for _, ch := range d.children {
		if dir, ok := ch.(*directory); ok {
			chClosest := dir.findClosestNeededSize(needed, closest)
			if chClosest < closest {
				closest = chClosest
			}
		}
	}
	return closest
}

type file struct {
	name string
	size int64
}

func (f file) Size() int64 {
	return f.size
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	input := strings.TrimSpace(string(bytes))
	fs := mapFileSystem(input)

	log.Println("Total size of directories < 100000:", fs.getSizeSum())
	log.Println(
		"Size of the smallest eligible directoy to delete:",
		fs.getSmallestEligibleDirectory(),
	)
}

func mapFileSystem(input string) *fileSystem {
	fs := newRootFileSystem()

	commands := strings.Split(input, "$")
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		lines := strings.Split(cmd, "\n")

		cmdLine := strings.TrimSpace(lines[0])
		args := strings.Split(cmdLine, " ")
		switch args[0] {
		case "cd":
			fs.cd(args[1])
		case "ls":
			fs.ls(lines[1:])
		}
	}

	return fs
}
