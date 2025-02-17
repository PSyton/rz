package rz

import (
	"io"
	"sync"
)

// LevelWriter defines as interface a writer may implement in order
// to receive level information with payload.
type LevelWriter interface {
	io.Writer
	WriteEvent(e *Event) (err error)
}

type levelWriterAdapter struct {
	io.Writer
}

func (lw levelWriterAdapter) WriteEvent(e *Event) (err error) {
	if e.formatter != nil {
		return e.formatter(lw, e)
	}

	n, err := lw.Write(e.buf)
	if err == nil && n < len(e.buf) {
		err = io.ErrShortWrite
	}
	return
}

type syncWriter struct {
	mu sync.Mutex
	lw LevelWriter
}

// SyncWriter wraps w so that each call to Write is synchronized with a mutex.
// This syncer can be the call to writer's Write method is not thread safe.
// Note that os.File Write operation is using write() syscall which is supposed
// to be thread-safe on POSIX systems. So there is no need to use this with
// os.File on such systems as rz guaranties to issue a single Write call
// per log event.
func SyncWriter(w io.Writer) io.Writer {
	if lw, ok := w.(LevelWriter); ok {
		return &syncWriter{lw: lw}
	}
	return &syncWriter{lw: levelWriterAdapter{w}}
}

// Write implements the io.Writer interface.
func (s *syncWriter) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lw.Write(p)
}

// WriteLevel implements the LevelWriter interface.
func (s *syncWriter) WriteEvent(e *Event) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lw.WriteEvent(e)
}

type multiLevelWriter struct {
	writers []LevelWriter
}

func (t multiLevelWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

func (t multiLevelWriter) WriteEvent(e *Event) (err error) {
	for _, w := range t.writers {
		err = w.WriteEvent(e)
		if err != nil {
			return
		}
	}
	return nil
}

// MultiLevelWriter creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command. If some writers
// implement LevelWriter, their WriteLevel method will be used instead of Write.
func MultiLevelWriter(writers ...io.Writer) LevelWriter {
	lwriters := make([]LevelWriter, 0, len(writers))
	for _, w := range writers {
		if lw, ok := w.(LevelWriter); ok {
			lwriters = append(lwriters, lw)
		} else {
			lwriters = append(lwriters, levelWriterAdapter{w})
		}
	}
	return multiLevelWriter{lwriters}
}
