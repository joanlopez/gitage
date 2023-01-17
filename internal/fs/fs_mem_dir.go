package fs

import "sort"

func RemoveFromMemDir(dir *MemFileData, f *MemFileData) {
	dir.memDir.Remove(f)
}

func AddToMemDir(dir *MemFileData, f *MemFileData) {
	if f.name != dir.name {
		dir.memDir.Add(f)
	}
}

func InitializeDir(d *MemFileData) {
	if d.memDir == nil {
		d.dir = true
		d.memDir = &MemDir{}
	}
}

type MemDir map[string]*MemFileData

func (d MemDir) Len() int {
	return len(d)
}

func (d MemDir) Add(f *MemFileData) {
	d[f.name] = f
}

func (d MemDir) Remove(f *MemFileData) {
	delete(d, f.name)
}

func (d MemDir) Files() (files []*MemFileData) {
	for _, f := range d {
		files = append(files, f)
	}
	sort.Sort(filesSorter(files))
	return files
}

var _ sort.Interface = filesSorter{}

type filesSorter []*MemFileData

func (s filesSorter) Len() int {
	return len(s)
}

func (s filesSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s filesSorter) Less(i, j int) bool {
	return s[i].name < s[j].name
}

func (d MemDir) Names() (names []string) {
	for x := range d {
		names = append(names, x)
	}
	return names
}
