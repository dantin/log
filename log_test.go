package logger

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
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

	require.True(t, strings.Contains(l, "[DEBUG]"))
	require.True(t, strings.Contains(l, "test debug log!"))
}

func TestInfoLog(t *testing.T) {
	bw, _, tearDown := setupTest("info")
	defer tearDown()

	Infof("test info log!")

	l := bw.String()

	require.True(t, strings.Contains(l, "[INFO]"))
	require.True(t, strings.Contains(l, "test info log!"))
}

func TestWarningLog(t *testing.T) {
	bw, _, tearDown := setupTest("warning")
	defer tearDown()

	Warnf("test warning log!")

	l := bw.String()

	require.True(t, strings.Contains(l, "[WARN]"))
	require.True(t, strings.Contains(l, "test warning log!"))
}

func TestErrorLog(t *testing.T) {
	bw, _, tearDown := setupTest("error")
	defer tearDown()

	Errorf("test error log!")

	l := bw.String()

	require.True(t, strings.Contains(l, "[ERROR]"))
	require.True(t, strings.Contains(l, "test error log!"))

	bw.Reset()

	Error(errors.New("some error string"))

	l = bw.String()
	require.True(t, strings.Contains(l, "some error string"))
}

func TestFatalLog(t *testing.T) {
	var exited bool
	exitHandler = func() {
		exited = true
	}

	bw, _, tearDown := setupTest("fatal")
	defer tearDown()

	Fatalf("test fatal log!")

	require.True(t, exited)

	l := bw.String()
	require.True(t, strings.Contains(l, "[FATAL]"))
	require.True(t, strings.Contains(l, "test fatal log!"))

	bw.Reset()
	exited = false

	Fatal(errors.New("some error string"))

	l = bw.String()
	require.True(t, strings.Contains(l, "some error string"))
}

func TestLogFile(t *testing.T) {
	bw, lf, tearDown := setupTest("debug")

	Debugf("test debug log!")

	require.Equal(t, bw.String(), lf.String())

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
