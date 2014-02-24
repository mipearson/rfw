package logbuf

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestLogbuf(t *testing.T) {
	tmpdir := makeTempDir(t)
	defer rmTempDir(tmpdir)

	l, err := Open(path.Join(tmpdir, "logbuf"), 0644)
	chkerr(t, err)

	// Basic writer tests
	_, err = l.Write([]byte("Hello there!\n"))
	chkerr(t, err)
	assertFileEquals(t, l, []byte("Hello there!\n"))

	_, err = l.Write([]byte("Goodbye.\n"))
	chkerr(t, err)
	assertFileEquals(t, l, []byte("Hello there!\nGoodbye.\n"))

	// Move the file
	err = os.Rename(l.path, l.path+".1")
	chkerr(t, err)
	_, err = l.Write([]byte("New content\n"))
	chkerr(t, err)
	assertFileEquals(t, l, []byte("New content\n"))

	// Delete the file
	err = os.Remove(l.path)
	chkerr(t, err)
	_, err = l.Write([]byte("More new content\n"))
	chkerr(t, err)
	assertFileEquals(t, l, []byte("More new content\n"))

	// Close & Re-open does not destroy the file
	err = l.Close()
	chkerr(t, err)
	l, err = Open(l.path, l.mode)
	chkerr(t, err)
	_, err = l.Write([]byte("foo\n"))
	chkerr(t, err)
	assertFileEquals(t, l, []byte("More new content\nfoo\n"))

}

/* Helpers */

func chkerr(t *testing.T, err error) {
	if err != nil {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Error encountered at line %d: %s", line, err)
	}
}

func makeTempDir(t *testing.T) string {
	tmpdir, err := ioutil.TempDir("", "logstash-config-test")
	chkerr(t, err)
	return tmpdir
}

func rmTempDir(tmpdir string) {
	_ = os.RemoveAll(tmpdir)
}

func assertFileEquals(t *testing.T, l *Logbuf, expected []byte) {
	contents, err := ioutil.ReadFile(l.path)
	chkerr(t, err)
	if !bytes.Equal(contents, expected) {
		_, _, line, _ := runtime.Caller(1)
		t.Fatalf("line %d: Expected %v, got %v from logbuf file.", line, expected, contents)
	}
}
