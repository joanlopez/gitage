package fs

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Only a subset of bits is allowed to be changed. Documented under os.Chmod().
const chmodBits = os.ModePerm | os.ModeSetuid | os.ModeSetgid | os.ModeSticky

var _ FS = &MemFS{}

type MemFS struct {
	mu   sync.RWMutex
	data map[string]*MemFileData
	init sync.Once
}

func NewMemFS() FS {
	return &MemFS{}
}

func (m *MemFS) getData() map[string]*MemFileData {
	m.init.Do(func() {
		m.data = make(map[string]*MemFileData)
		// Root should always exist, right?
		// TODO: what about windows?
		root := CreateDir(string(os.PathSeparator))
		SetMode(root, os.ModeDir|0o755)
		m.data[string(os.PathSeparator)] = root
	})
	return m.data
}

func (m *MemFS) Create(name string) (File, error) {
	name = normalizePath(name)
	m.mu.Lock()
	file := CreateFile(name)
	m.getData()[name] = file
	m.registerWithParent(file, 0)
	m.mu.Unlock()
	return NewFileHandle(file), nil
}

func (m *MemFS) unRegisterWithParent(fileName string) error {
	f, err := m.lockFreeOpen(fileName)
	if err != nil {
		return err
	}
	parent := m.findParent(f)
	if parent == nil {
		log.Panic("parent of ", f.Name(), " is nil")
	}

	parent.Lock()
	RemoveFromMemDir(parent, f)
	parent.Unlock()
	return nil
}

func (m *MemFS) findParent(f *MemFileData) *MemFileData {
	parentDir, _ := filepath.Split(f.Name())
	parentDir = filepath.Clean(parentDir)
	parentFile, err := m.lockFreeOpen(parentDir)
	if err != nil {
		return nil
	}
	return parentFile
}

func (m *MemFS) registerWithParent(f *MemFileData, perm os.FileMode) {
	if f == nil {
		return
	}
	parent := m.findParent(f)
	if parent == nil {
		parentDir := filepath.Dir(filepath.Clean(f.Name()))
		err := m.lockFreeMkdir(parentDir, perm)
		if err != nil {
			return
		}
		parent, err = m.lockFreeOpen(parentDir)
		if err != nil {
			return
		}
	}

	parent.Lock()
	InitializeDir(parent)
	AddToMemDir(parent, f)
	parent.Unlock()
}

func (m *MemFS) lockFreeMkdir(name string, perm os.FileMode) error {
	name = normalizePath(name)
	x, ok := m.getData()[name]
	if ok {
		// Only return ErrFileExists if it's a file, not a directory.
		i := MemFileInfo{MemFileData: x}
		if !i.IsDir() {
			return ErrFileExists
		}
	} else {
		item := CreateDir(name)
		SetMode(item, os.ModeDir|perm)
		m.getData()[name] = item
		m.registerWithParent(item, perm)
	}
	return nil
}

func (m *MemFS) Mkdir(name string, perm os.FileMode) error {
	perm &= chmodBits
	name = normalizePath(name)

	m.mu.RLock()
	_, ok := m.getData()[name]
	m.mu.RUnlock()
	if ok {
		return &os.PathError{Op: "mkdir", Path: name, Err: ErrFileExists}
	}

	m.mu.Lock()
	// Double check that it doesn't exist.
	if _, ok := m.getData()[name]; ok {
		m.mu.Unlock()
		return &os.PathError{Op: "mkdir", Path: name, Err: ErrFileExists}
	}
	item := CreateDir(name)
	SetMode(item, os.ModeDir|perm)
	m.getData()[name] = item
	m.registerWithParent(item, perm)
	m.mu.Unlock()

	return m.setFileMode(name, perm|os.ModeDir)
}

func (m *MemFS) MkdirAll(path string, perm os.FileMode) error {
	err := m.Mkdir(path, perm)
	if err != nil {
		// nolint: forcetypeassert
		if err.(*os.PathError).Err == ErrFileExists {
			return nil
		}
		return err
	}
	return nil
}

// Handle some relative paths.
func normalizePath(path string) string {
	path = filepath.Clean(path)

	switch path {
	case ".":
		fallthrough
	case "..":
		return string(os.PathSeparator)
	default:
		return path
	}
}

func (m *MemFS) Open(name string) (File, error) {
	f, err := m.open(name)
	if f != nil {
		return NewReadOnlyFileHandle(f), err
	}
	return nil, err
}

func (m *MemFS) openWrite(name string) (File, error) {
	f, err := m.open(name)
	if f != nil {
		return NewFileHandle(f), err
	}
	return nil, err
}

func (m *MemFS) open(name string) (*MemFileData, error) {
	name = normalizePath(name)

	m.mu.RLock()
	f, ok := m.getData()[name]
	m.mu.RUnlock()
	if !ok {
		return nil, &os.PathError{Op: "open", Path: name, Err: ErrFileNotFound}
	}
	return f, nil
}

func (m *MemFS) lockFreeOpen(name string) (*MemFileData, error) {
	name = normalizePath(name)
	f, ok := m.getData()[name]
	if ok {
		return f, nil
	}
	return nil, ErrFileNotFound
}

func (m *MemFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	perm &= chmodBits
	chmod := false
	file, err := m.openWrite(name)
	if err == nil && (flag&os.O_EXCL > 0) {
		return nil, &os.PathError{Op: "open", Path: name, Err: ErrFileExists}
	}
	if os.IsNotExist(err) && (flag&os.O_CREATE > 0) {
		file, err = m.Create(name)
		chmod = true
	}
	if err != nil {
		return nil, err
	}
	if flag == os.O_RDONLY {
		// nolint: forcetypeassert
		file = NewReadOnlyFileHandle(file.(*MemFile).Data())
	}
	if flag&os.O_APPEND > 0 {
		_, err = file.(*MemFile).Seek(0, io.SeekEnd)
		if err != nil {
			if closeErr := file.Close(); closeErr != nil {
				return nil, fmt.Errorf("error seeking and closing file: %v, %v", err, closeErr)
			}
			return nil, err
		}
	}
	if flag&os.O_TRUNC > 0 && flag&(os.O_RDWR|os.O_WRONLY) > 0 {
		// nolint: forcetypeassert
		err = file.(*MemFile).Truncate(0)
		if err != nil {
			if closeErr := file.Close(); closeErr != nil {
				return nil, fmt.Errorf("error truncating and closing file: %v, %v", err, closeErr)
			}
			return nil, err
		}
	}
	if chmod {
		return file, m.setFileMode(name, perm)
	}
	return file, nil
}

func (m *MemFS) Remove(name string) error {
	name = normalizePath(name)

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.getData()[name]; ok {
		err := m.unRegisterWithParent(name)
		if err != nil {
			return &os.PathError{Op: "remove", Path: name, Err: err}
		}
		delete(m.getData(), name)
	} else {
		return &os.PathError{Op: "remove", Path: name, Err: os.ErrNotExist}
	}
	return nil
}

func (m *MemFS) RemoveAll(path string) error {
	path = normalizePath(path)
	m.mu.Lock()
	_ = m.unRegisterWithParent(path)
	m.mu.Unlock()

	m.mu.RLock()
	defer m.mu.RUnlock()

	for p := range m.getData() {
		if p == path || strings.HasPrefix(p, path+string(os.PathSeparator)) {
			m.mu.RUnlock()
			m.mu.Lock()
			delete(m.getData(), p)
			m.mu.Unlock()
			m.mu.RLock()
		}
	}
	return nil
}

func (m *MemFS) Lstat(name string) (os.FileInfo, error) {
	return m.Stat(name)
}

func (m *MemFS) Stat(name string) (os.FileInfo, error) {
	f, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	fi := GetFileInfo(f.(*MemFile).Data())
	return fi, nil
}

func (m *MemFS) setFileMode(name string, mode os.FileMode) error {
	name = normalizePath(name)

	m.mu.RLock()
	f, ok := m.getData()[name]
	m.mu.RUnlock()
	if !ok {
		return &os.PathError{Op: "chmod", Path: name, Err: ErrFileNotFound}
	}

	m.mu.Lock()
	SetMode(f, mode)
	m.mu.Unlock()

	return nil
}
