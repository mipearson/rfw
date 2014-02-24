package logbuf

import (
	"os"
	"syscall"
)

type Logbuf struct {
	path  string
	mode  os.FileMode
	file  *os.File
	inode uint64
}

/*
Open returns a new Logbuf at the specified path.

If the file does not exist it will be created with the specified mode.
If the file does exist it will be appended to.
*/
func Open(path string, mode os.FileMode) (*Logbuf, error) {
	var logbuf Logbuf
	logbuf.path = path
	logbuf.mode = mode
	err := logbuf.open()
	return &logbuf, err
}

/*
Write p bytes to our file.

If our file has been deleted or has been moved out from under us,
a new file will be created.
*/
func (l *Logbuf) Write(p []byte) (int, error) {
	inode, err := l.checkInode()
	if os.IsNotExist(err) || inode != l.inode {
		err = l.reopen()
		if err != nil {
			return 0, err
		}
	}
	return l.file.Write(p)
}

/* Close our file */
func (l *Logbuf) Close() error {
	return l.file.Close()
}

func (l *Logbuf) checkInode() (uint64, error) {
	var stat syscall.Stat_t
	err := syscall.Stat(l.path, &stat)
	return stat.Ino, err
}

func (l *Logbuf) reopen() error {
	if err := l.Close(); err != nil {
		return err
	}
	return l.open()
}

func (l *Logbuf) open() error {
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
