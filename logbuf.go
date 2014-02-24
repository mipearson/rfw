package logbuf

import (
	"io"
	"os"
	"syscall"
)

type Logbuf struct {
	path  string
	mode  os.FileMode
	file  os.File
	inode uint64
}

func Open(path string, mode os.FileMode) (Logbuf, error) {
	var logbuf Logbuf
	logbuf.path = path
	logbuf.mode = mode
	err := logbuf.open()
	return logbuf, err
}

func (Logbuf l) Write(p []byte) (int, error) {
	inode, err := l.checkInode()
	if os.IsNotExist(err) || inode != l.inode {
		err = l.reopen()
		if err != nil {
			return 0, err
		}
	}
	return l.file.Write(p)
}

func (Logbuf l) Close() error {
	return l.file.Close()
}

func (Logbuf l) checkInode() (uint64, error) {
	var stat syscall.Stat_t
	err := syscall.Stat(path, &stat)
	return stat.Ino, err
}

func (Logbuf l) reopen() error {
	l.Close()
	l.open()
}

func (Logbuf l) open() error {
	l.file, err = os.FileOpen(l.path, os.O_APPEND|os.O_WRONLY, l.mode)
	if err != nil {
		return err
	}

	// TODO: Possible race here, but worst that will happen is next
	// write will close/reopen
	l.inode, err = logbuf.checkInode()
	if err != nil {
		_ = l.Close()
		return err
	}
	l.inode = stat.Ino

	return nil
}
