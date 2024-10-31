package main

import (
	"io"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"
)

type SortDir []fs.DirEntry

func (a SortDir) Len() int      { return len(a) }
func (a SortDir) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortDir) Less(i, j int) bool {
	first, _ := a[i].Info()
	second, _ := a[j].Info()
	res := strings.Compare(strings.Split(first.Name(), ".")[0], strings.Split(second.Name(), ".")[0])
	return res < 0
}

func findLastDir(dir []fs.DirEntry) int {
	lastDir := 0
	for i := range dir {
		if strings.Contains(fs.FormatDirEntry(dir[i]), "d") {
			lastDir = i
		}
	}
	return lastDir
}
func fileInfo(file fs.DirEntry, printFiles bool) string {
	info, _ := file.Info()
	fileinfo := ""
	if info.IsDir() {
		fileinfo = info.Name()
	} else if printFiles {
		checksize := info.Size()
		strchecksize := strconv.Itoa(int(checksize))
		if checksize != 0 {
			fileinfo = info.Name() + "(" + strchecksize + "b)"
		} else {
			fileinfo = info.Name() + "(empty)"
		}
	}
	return fileinfo
}
func checkIgnored(dir []fs.DirEntry) []fs.DirEntry {
	Ignored := []string{".DS_Store", ".git"}
	for i, f := range dir {
		for _, s := range Ignored {
			if strings.Contains(fs.FormatDirEntry(f), s) {
				dir = append(dir[:i], dir[i+1:]...)
			}

		}
	}
	return dir
}
func sortDir(dir []fs.DirEntry) []fs.DirEntry {
	comparator := SortDir(dir)
	sort.Sort(comparator)
	return dir
}
func showTree(path string, printFiles bool, tab string) string {
	indir := "├───"
	nextdir := "└───"
	pref := tab + indir
	result := ""
	islast := false
	dirPath := path + string(os.PathSeparator)
	dir, err := os.ReadDir(dirPath)
	if err == nil {
		dir = sortDir(checkIgnored(dir))
		islastDir := findLastDir(dir)
		for i, elem := range dir {
			file := fileInfo(elem, printFiles)
			if len(file) == 0 {
				continue
			}
			if i == len(dir)-1 || (i == islastDir && !printFiles) {
				pref = tab + nextdir
				islast = true
			}
			if !strings.Contains(file, ".") && !strings.HasSuffix(file, ")") {
				result += pref + file + "\n"
				if islast {
					result += showTree(dirPath+file, printFiles, tab+"\t")
				} else {
					result += showTree(dirPath+file, printFiles, tab+"│\t")
				}
			} else {
				result += pref + strings.ReplaceAll(file, "(", " (") + "\n"
			}
		}
	}
	return result
}
func dirTree(out io.Writer, path string, printFiles bool) error {
	result := showTree(path, printFiles, "")
	out.Write([]byte(result))
	return nil
}
func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
