package fs

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

var _ os.FileInfo = &MemFileInfo{}

type MemFileData struct {
	sync.Mutex
	name    string
	data    []byte
	memDir  *MemDir
	dir     bool
	mode    os.FileMode
	modtime time.Time
}

func (d *MemFileData) Name() string {
	d.Lock()
	defer d.Unlock()
	return d.name
}

func GetFileInfo(f *MemFileData) *MemFileInfo {
	return &MemFileInfo{MemFileData: f}
}

type MemFileInfo struct {
	*MemFileData
}

func (s *MemFileInfo) Name() string {
	s.Lock()
	_, name := filepath.Split(s.name)
	s.Unlock()
	return name
}

func (s *MemFileInfo) Mode() os.FileMode {
	s.Lock()
	defer s.Unlock()
	return s.mode
}

func (s *MemFileInfo) ModTime() time.Time {
	s.Lock()
	defer s.Unlock()
	return time.Now()
}

func (s *MemFileInfo) IsDir() bool {
	s.Lock()
	defer s.Unlock()
	return s.dir
}

func (s *MemFileInfo) Sys() interface{} {
	return nil
}

func (s *MemFileInfo) Size() int64 {
	if s.IsDir() {
		return int64(42)
	}
	s.Lock()
	defer s.Unlock()
	return int64(len(s.data))
}
