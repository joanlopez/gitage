package fs

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

var _ File = &MemFile{}

type MemFile struct {
	// atomic requires 64-bit alignment for struct field access
	at           int64
	readDirCount int64
	closed       bool
	readOnly     bool
	fileData     *MemFileData
}

func NewFileHandle(data *MemFileData) *MemFile {
	return &MemFile{fileData: data}
}

func NewReadOnlyFileHandle(data *MemFileData) *MemFile {
	return &MemFile{fileData: data, readOnly: true}
}

func (f *MemFile) Data() *MemFileData {
	return f.fileData
}

func CreateFile(name string) *MemFileData {
	return &MemFileData{name: name, mode: os.ModeTemporary, modtime: time.Now()}
}

func CreateDir(name string) *MemFileData {
	return &MemFileData{name: name, memDir: &MemDir{}, dir: true, modtime: time.Now()}
}

func SetMode(f *MemFileData, mode os.FileMode) {
	f.Lock()
	f.mode = mode
	f.Unlock()
}

func (f *MemFile) Close() error {
	return nil
}

func (f *MemFile) Stat() (os.FileInfo, error) {
	return &MemFileInfo{MemFileData: f.fileData}, nil
}

func (f *MemFile) Readdirnames(n int) (names []string, err error) {
	fi, err := f.readDir(n)
	names = make([]string, len(fi))
	for i, f := range fi {
		_, names[i] = filepath.Split(f.Name())
	}
	return names, err
}

func (f *MemFile) readDir(count int) (res []os.FileInfo, err error) {
	if !f.fileData.dir {
		return nil, &os.PathError{Op: "readdir", Path: f.fileData.name, Err: errors.New("not a dir")}
	}
	var outLength int64

	f.fileData.Lock()
	files := f.fileData.memDir.Files()[f.readDirCount:]
	if count > 0 {
		if len(files) < count {
			outLength = int64(len(files))
		} else {
			outLength = int64(count)
		}
		if len(files) == 0 {
			err = io.EOF
		}
	} else {
		outLength = int64(len(files))
	}
	f.readDirCount += outLength
	f.fileData.Unlock()

	res = make([]os.FileInfo, outLength)
	for i := range res {
		res[i] = &MemFileInfo{MemFileData: files[i]}
	}

	return res, err
}

func (f *MemFile) Read(b []byte) (n int, err error) {
	f.fileData.Lock()
	defer f.fileData.Unlock()
	if f.closed {
		return 0, ErrFileClosed
	}
	if len(b) > 0 && int(f.at) == len(f.fileData.data) {
		return 0, io.EOF
	}
	if int(f.at) > len(f.fileData.data) {
		return 0, io.ErrUnexpectedEOF
	}
	if len(f.fileData.data)-int(f.at) >= len(b) {
		n = len(b)
	} else {
		n = len(f.fileData.data) - int(f.at)
	}
	copy(b, f.fileData.data[f.at:f.at+int64(n)])
	atomic.AddInt64(&f.at, int64(n))
	return
}

func (f *MemFile) Truncate(size int64) error {
	if f.closed {
		return ErrFileClosed
	}
	if f.readOnly {
		return &os.PathError{Op: "truncate", Path: f.fileData.name, Err: errors.New("file handle is read only")}
	}
	if size < 0 {
		return ErrOutOfRange
	}
	f.fileData.Lock()
	defer f.fileData.Unlock()
	if size > int64(len(f.fileData.data)) {
		diff := size - int64(len(f.fileData.data))
		f.fileData.data = append(f.fileData.data, bytes.Repeat([]byte{0o0}, int(diff))...)
	} else {
		f.fileData.data = f.fileData.data[0:size]
	}
	setModTime(f.fileData, time.Now())
	return nil
}

func (f *MemFile) Seek(offset int64, whence int) (int64, error) {
	if f.closed {
		return 0, ErrFileClosed
	}
	switch whence {
	case io.SeekStart:
		atomic.StoreInt64(&f.at, offset)
	case io.SeekCurrent:
		atomic.AddInt64(&f.at, offset)
	case io.SeekEnd:
		atomic.StoreInt64(&f.at, int64(len(f.fileData.data))+offset)
	}
	return f.at, nil
}

func (f *MemFile) Write(b []byte) (n int, err error) {
	if f.closed {
		return 0, ErrFileClosed
	}
	if f.readOnly {
		return 0, &os.PathError{Op: "write", Path: f.fileData.name, Err: errors.New("file handle is read only")}
	}
	n = len(b)
	cur := atomic.LoadInt64(&f.at)
	f.fileData.Lock()
	defer f.fileData.Unlock()
	diff := cur - int64(len(f.fileData.data))
	var tail []byte
	if n+int(cur) < len(f.fileData.data) {
		tail = f.fileData.data[n+int(cur):]
	}
	if diff > 0 {
		f.fileData.data = append(f.fileData.data, append(bytes.Repeat([]byte{0o0}, int(diff)), b...)...)
		f.fileData.data = append(f.fileData.data, tail...)
	} else {
		f.fileData.data = append(f.fileData.data[:cur], b...)
		f.fileData.data = append(f.fileData.data, tail...)
	}
	setModTime(f.fileData, time.Now())

	atomic.AddInt64(&f.at, int64(n))
	return
}

func setModTime(f *MemFileData, mtime time.Time) {
	f.modtime = mtime
}

var (
	ErrFileClosed        = errors.New("MemFile is closed")
	ErrOutOfRange        = errors.New("out of range")
	ErrTooLarge          = errors.New("too large")
	ErrFileNotFound      = os.ErrNotExist
	ErrFileExists        = os.ErrExist
	ErrDestinationExists = os.ErrExist
)
