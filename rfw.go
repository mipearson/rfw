// Package rfw implements a log rotation aware file writer.
// It will always write to the path that you give it, even if that file is
// deleted or moved out from under it.
package rfw

import (
	"os"
	"sync"
	"syscall"
)

// A file rotation aware writer
type Writer struct {
	path  string
	mode  os.FileMode
	file  *os.File
	inode uint64
	mutex sync.Mutex
}

/*
Open returns a new Writer at the specified path.

If the file does not exist it will be created with the specified mode.
If the file does exist it will be appended to.
*/
func Open(path string, mode os.FileMode) (*Writer, error) {
	var w Writer
	w.path = path
	w.mode = mode
	err := w.open()
	return &w, err
}

/*
Write p bytes to our file.

If our file has been deleted or has been moved out from under us,
a new file will be created with the mode specified at Open time.
*/
func (l *Writer) Write(p []byte) (int, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	inode, err := l.checkInode()
	if os.IsNotExist(err) || inode != l.inode {
		err = l.reopen()
		if err != nil {
			return 0, err
		}
	}
	return l.file.Write(p)
}

// Close our writer. Subsequent writes will fail.
func (l *Writer) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.file.Close()
}

func (l *Writer) checkInode() (uint64, error) {
	var stat syscall.Stat_t
	err := syscall.Stat(l.path, &stat)
	return uint64(stat.Ino), err
}

func (l *Writer) reopen() error {
	if err := l.file.Close(); err != nil {
		return err
	}
	return l.open()
}

func (l *Writer) open() error {
	var err error
	l.file, err = os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, l.mode)
	if err != nil {
		return err
	}

	// TODO: Possible race here, but worst that will happen is next
	// write will close/reopen
	l.inode, err = l.checkInode()
	if err != nil {
		_ = l.Close()
		return err
	}

	return nil
}
