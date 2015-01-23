package util

import (
	"bytes"
	"io"
	"testing"
)

func makeStreams() (*ReadSeekCloseWrapper, io.WriteCloser) {
	r, w := io.Pipe()
	s := WrapReadSeekClose(r)
	return s, w
}

func TestRead(t *testing.T) {
	r, w := makeStreams()

	buf := make([]byte, 10)
	go func() {
		if _, err := w.Write([]byte("hello")); err != nil {
			t.Error(err)
		}
		if err := w.Close(); err != nil {
			t.Error(err)
		}
	}()

	count, err := r.Read(buf)
	if err != nil {
		t.Error(err)
	}
	if count != 5 {
		t.Errorf("Read returned %v, expected %v", count, 5)
	}
}

func TestClose(t *testing.T) {
	r, w := makeStreams()

	if err := r.Close(); err != nil {
		t.Error(err)
	}

	if count, err := w.Write([]byte("hello")); err == nil {
		t.Errorf("Successfully wrote %v bytes after closing", count)
	}
}

func TestSeekSet(t *testing.T) {
	r, w := makeStreams()

	if _, err := w.Write([]byte("hello")); err != nil {
		t.Error(err)
	}

	if pos, err := r.Seek(2, 0); err == nil {
		t.Errorf("Seek via SEEK_SET should not be implemented; got %v", pos)
	}
}

func TestSeekCur(t *testing.T) {
	r, w := makeStreams()

	if _, err := w.Write([]byte("hello")); err != nil {
		t.Error(err)
	}

	if pos, err := r.Seek(2, 1); err == nil {
		t.Errorf("Seek via SEEK_CUR should not be implemented; got %v", pos)
	}
}

func TestSeekEndZero(t *testing.T) {
	r, w := makeStreams()

	if _, err := w.Write([]byte("hello")); err != nil {
		t.Error(err)
	}

	_, err := r.Seek(0, 2)
	if err != nil {
		t.Error(err)
	}

	if _, err := w.Write([]byte("world")); err != nil {
		t.Error(err)
	}

	buf := make([]byte, 10)
	count, err := r.Read(buf)

	if err != nil {
		t.Error(err)
	}

	if count != 5 {
		t.Errorf("Got %v bytes in %v, expected 5", count, buf)
	}

	expected := []byte("world\x00\x00\x00\x00\x00")
	if bytes.Compare(buf, expected) != 0 {
		t.Errorf("Got unexpected content %v, expected %v", buf, expected)
	}
}

func TestSeekEndNonZero(t *testing.T) {
	r, w := makeStreams()

	if _, err := w.Write([]byte("hello")); err != nil {
		t.Error(err)
	}

	_, err := r.Seek(-2, 2)
	if err != nil {
		t.Error(err)
	}

	buf := make([]byte, 10)
	count, err := r.Read(buf)

	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Errorf("Got %v bytes in %v, expected 5", count, buf)
	}

	expected := []byte("lo\x00\x00\x00\x00\x00\x00\x00\x00")
	if bytes.Compare(buf, expected) != 0 {
		t.Errorf("Got unexpected content %v, expected %v", buf, expected)
	}
}

func TestSeekEndPastEOF(t *testing.T) {
	r, w := makeStreams()

	if _, err := w.Write([]byte("hello")); err != nil {
		t.Error(err)
	}

	_, err := r.Seek(2, 2)
	if err == nil {
		t.Errorf("Seeking past EOF should not be supported")
	}
}
