// Wrapper for seekable stream

package util

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/hpcloud/log"
)

type ReadSeekCloseWrapper struct {
	io.ReadCloser
	io.Seeker

	available *sync.Cond
	reader    io.ReadCloser
	buffer    bytes.Buffer
}

func WrapReadSeekClose(reader io.ReadCloser) *ReadSeekCloseWrapper {
	r := &ReadSeekCloseWrapper{
		available: sync.NewCond(&sync.Mutex{}),
		reader:    reader,
	}

	go r.Drain()

	return r
}

func (r *ReadSeekCloseWrapper) Drain() {
	var err, err2 error
	buf := make([]byte, 0x10000)
	for {
		var count int
		count, err = r.reader.Read(buf)
		if count > 0 {
			_, err2 = r.buffer.Write(buf[:count])
		}
		if err != nil && err != io.EOF && err != io.ErrClosedPipe {
			log.Errorf("Error reading stream %v: %v", r.reader, err)
		}
		if err2 != nil {
			log.Errorf("Error writing buffer: %v: %v", r.buffer, err2)
		}
		if err != nil || err2 != nil {
			break
		}
		if r.buffer.Len() > 0 {
			r.available.L.Lock()
			r.available.Broadcast()
			r.available.L.Unlock()
		}
	}
	log.Infof("Read complete (error %v/%v)", err, err2)
}

func (r *ReadSeekCloseWrapper) Read(p []byte) (n int, err error) {
	r.available.L.Lock()
	for r.buffer.Len() < 1 {
		r.available.Wait()
	}
	n, err = r.buffer.Read(p)
	if r.buffer.Len() > 0 {
		r.available.Broadcast()
	}
	r.available.L.Unlock()
	return n, err
}

func (r *ReadSeekCloseWrapper) Seek(offset int64, whence int) (int64, error) {
	const MaxInt = int64((^uint(0)) >> 1)
	log.Infof("Attempting to seek to whence %v offset %v", whence, offset)
	if whence != 2 {
		return 0, fmt.Errorf("Seeking to %v not implemented", whence)
	}
	if offset > 0 {
		return 0, fmt.Errorf("Seeking past EOF to %v not supported", offset)
	}
	if offset == 0 {
		// Optimize discarding everything (common case)
		r.buffer.Reset()
	} else {
		discard := int64(r.buffer.Len()) + offset
		for discard > 0 {
			chunkSize := discard % MaxInt
			r.buffer.Next(int(chunkSize))
			discard -= chunkSize
		}
	}
	return 0, nil
}

func (r *ReadSeekCloseWrapper) Close() error {
	return r.reader.Close()
}
