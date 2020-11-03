package logger

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"
)

type writeBuffer struct {
	mu  sync.RWMutex
	buf *bytes.Buffer
}

func newWriterBuffer() *writeBuffer {
	return &writeBuffer{buf: bytes.NewBuffer(nil)}
}

func (wb *writeBuffer) Write(p []byte) (int, error) {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	return wb.buf.Write(p)
}

func (wb *writeBuffer) String() string {
	wb.mu.RLock()
	defer wb.mu.RUnlock()
	return wb.buf.String()
}

func (wb *writeBuffer) Reset() {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	wb.buf.Reset()
}

func (wb *writeBuffer) Close() error {
	return nil
}

func TestDebugLog(t *testing.T) {
	bw, _, tearDown := setupTest("debug")
	defer tearDown()

	Debugf("test debug log!")

	l := bw.String()

	if !strings.Contains(l, "[DEBUG]") {
		t.Fatalf("Debug flag not equal")
	}
	if !strings.Contains(l, "test debug log!") {
		t.Fatalf("Debug message not equal")
	}
}

func TestInfoLog(t *testing.T) {
	bw, _, tearDown := setupTest("info")
	defer tearDown()

	Infof("test info log!")

	l := bw.String()

	if !strings.Contains(l, "[INFO]") {
		t.Fatalf("Info flag not equal")
	}
	if !strings.Contains(l, "test info log!") {
		t.Fatalf("Info message not equal")
	}
}

func TestWarningLog(t *testing.T) {
	bw, _, tearDown := setupTest("warning")
	defer tearDown()

	Warnf("test warning log!")

	l := bw.String()

	if !strings.Contains(l, "[WARN]") {
		t.Fatalf("Warn flag not equal")
	}
	if !strings.Contains(l, "test warning log!") {
		t.Fatalf("Warn message not equal")
	}
}

func TestErrorLog(t *testing.T) {
	bw, _, tearDown := setupTest("error")
	defer tearDown()

	Errorf("test error log!")

	l := bw.String()

	if !strings.Contains(l, "[ERROR]") {
		t.Fatalf("Error flag not equal")
	}
	if !strings.Contains(l, "test error log!") {
		t.Fatalf("Error message not equal")
	}

	bw.Reset()

	Error(errors.New("some error string"))

	l = bw.String()
	if !strings.Contains(l, "some error string") {
		t.Fatalf("Error message not equal")
	}
}

func TestFatalLog(t *testing.T) {
	var exited bool
	exitHandler = func() {
		exited = true
	}

	bw, _, tearDown := setupTest("fatal")
	defer tearDown()

	Fatalf("test fatal log!")

	if !exited {
		t.Fatalf("Fatal should exit")
	}

	l := bw.String()
	if !strings.Contains(l, "[FATAL]") {
		t.Fatalf("Fatal flag not equal")
	}
	if !strings.Contains(l, "test fatal log!") {
		t.Fatalf("Fatal message not equal")
	}

	bw.Reset()
	exited = false

	Fatal(errors.New("some error string"))

	l = bw.String()
	if !strings.Contains(l, "some error string") {
		t.Fatalf("Fatal message not equal")
	}
}

func TestLogFile(t *testing.T) {
	bw, lf, tearDown := setupTest("debug")

	Debugf("test debug log!")

	if bw.String() != lf.String() {
		t.Fatalf("not equal")
	}

	// make sure file is closed
	tearDown()
}

func setupTest(level string) (*writeBuffer, *writeBuffer, func()) {
	output := newWriterBuffer()
	logFile := newWriterBuffer()
	l, _ := New(level, output, logFile)
	Set(l)
	return output, logFile, func() { Unset() }
}
