package hosting

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/mitz-it/go-toolkit/logger"
)

type MockWorker struct{}

func (m *MockWorker) Run(ctx context.Context, exit func()) {
	defer exit()
	// Simulate work
	time.Sleep(500 * time.Millisecond)
	logger.Info(ctx).Msg("Worker is running")
}

func TestHost(t *testing.T) {
	var suts = map[string]struct {
		arrange func() *Host
		act     func(t *testing.T, h *Host) *bytes.Buffer
		assert  func(t *testing.T, output *bytes.Buffer)
		cleanup func(t *testing.T)
	}{
		"host with router": {
			arrange: func() *Host {
				router := http.NewServeMux()
				router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("Hello, World!"))
				})
				return &Host{
					Router: router,
					Addr:   ":8081",
				}
			},
			act: func(t *testing.T, h *Host) *bytes.Buffer {
				output := &bytes.Buffer{}
				go h.Start()
				time.Sleep(1 * time.Second)

				resp, err := http.Get("http://localhost:8081/")
				if err != nil {
					t.Fatalf("Failed to send GET request: %v", err)
				}
				defer resp.Body.Close()

				output.ReadFrom(resp.Body)
				return output
			},
			assert: func(t *testing.T, output *bytes.Buffer) {
				expected := "Hello, World!"
				if output.String() != expected {
					t.Errorf("expected %q, but got %q", expected, output.String())
				}
			},
			cleanup: func(t *testing.T) {
				p, err := os.FindProcess(os.Getpid())
				if err != nil {
					t.Fatalf("Failed to find process: %v", err)
				}

				if err := p.Signal(syscall.SIGINT); err != nil {
					t.Fatalf("Failed to send SIGINT: %v", err)
				}

				time.Sleep(1 * time.Second)
			},
		},
		"host with workers only": {
			arrange: func() *Host {
				worker := &MockWorker{}
				return &Host{
					Workers: []Worker{worker},
				}
			},
			act: func(t *testing.T, h *Host) *bytes.Buffer {
				output := &bytes.Buffer{}
				go h.Start()
				time.Sleep(1 * time.Second)
				output.WriteString("Worker executed")
				return output
			},
			assert: func(t *testing.T, output *bytes.Buffer) {
				expected := "Worker executed"
				if output.String() != expected {
					t.Errorf("expected %q, but got %q", expected, output.String())
				}
			},
			cleanup: func(t *testing.T) {},
		},
		"host with no router or workers": {
			arrange: func() *Host {
				return &Host{}
			},
			act: func(t *testing.T, h *Host) *bytes.Buffer {
				output := &bytes.Buffer{}
				go h.Start()
				time.Sleep(1 * time.Second)
				output.WriteString("Host not started due to no router or workers")
				return output
			},
			assert: func(t *testing.T, output *bytes.Buffer) {
				expected := "Host not started due to no router or workers"
				if output.String() != expected {
					t.Errorf("expected %q, but got %q", expected, output.String())
				}
			},
			cleanup: nil,
		},
	}

	for name, sut := range suts {
		t.Run(name, func(t *testing.T) {
			h := sut.arrange()

			o := sut.act(t, h)

			sut.assert(t, o)

			if sut.cleanup != nil {
				sut.cleanup(t)
			}
		})
	}
}
