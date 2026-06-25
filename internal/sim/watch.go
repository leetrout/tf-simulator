package sim

import (
	"os"
	"path/filepath"
	"time"
)

// stateFilePath returns the terraform.tfstate path for a working directory.
func stateFilePath(workDir string) string {
	return filepath.Join(workDir, "terraform.tfstate")
}

// watcher polls a file's modification time (and size) and invokes onChange when it
// changes. A poll is dependency-free and adequate for a local teaching tool.
type watcher struct {
	path     string
	onChange func()
	interval time.Duration
	stopCh   chan struct{}
}

func newWatcher(path string, onChange func()) *watcher {
	return &watcher{
		path:     path,
		onChange: onChange,
		interval: 500 * time.Millisecond,
		stopCh:   make(chan struct{}),
	}
}

// fingerprint returns a value that changes whenever the file is written. Returns
// ("", false) when the file is absent.
func (w *watcher) fingerprint() (string, bool) {
	fi, err := os.Stat(w.path)
	if err != nil {
		return "", false
	}
	return fi.ModTime().Format(time.RFC3339Nano) + ":" + itoa(fi.Size()), true
}

func (w *watcher) start() {
	last, _ := w.fingerprint()
	ticker := time.NewTicker(w.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-w.stopCh:
				return
			case <-ticker.C:
				cur, ok := w.fingerprint()
				if !ok {
					last = ""
					continue
				}
				if cur != last {
					last = cur
					w.onChange()
				}
			}
		}
	}()
}

func (w *watcher) stop() {
	select {
	case <-w.stopCh:
	default:
		close(w.stopCh)
	}
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
