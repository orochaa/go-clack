package prompts_test

import (
	"os"
	"sync"
	"time"
)

type MockDirEntry struct {
	name  string
	isDir bool
}

func (e MockDirEntry) Name() string {
	return e.name
}

func (e MockDirEntry) IsDir() bool {
	return e.isDir
}

func (e MockDirEntry) Type() os.FileMode {
	if e.isDir {
		return os.ModeDir
	}
	return 0
}

func (e MockDirEntry) Info() (os.FileInfo, error) {
	return nil, nil
}

type MockFileSystem struct{}

func (fs MockFileSystem) Getwd() (string, error) {
	return "/clack", nil
}

func (fs MockFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return []os.DirEntry{
		MockDirEntry{name: "dir", isDir: true},
		MockDirEntry{name: "file", isDir: false},
	}, nil
}

func (fs MockFileSystem) UserHomeDir() (string, error) {
	return "/home/clack", nil
}

type MockTimer struct {
	mu          sync.Mutex
	waiters     []chan struct{}
	autoResolve bool
}

func (t *MockTimer) Sleep(duration time.Duration) {
	if !t.autoResolve {
		waiter := make(chan struct{})
		t.mu.Lock()
		t.waiters = append(t.waiters, waiter)
		t.mu.Unlock()
		<-waiter
	}
}

func (m *MockTimer) ResolveAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, waiter := range m.waiters {
		close(waiter)
	}
	m.waiters = []chan struct{}(nil)
}

type MockWriter struct {
	mu   sync.Mutex
	Data []string
}

func (w *MockWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	w.Data = append(w.Data, string(data))
	w.mu.Unlock()
	return 0, nil
}
