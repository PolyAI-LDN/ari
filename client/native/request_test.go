package native

import (
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransportKeepsIdlePool(t *testing.T) {
	c := New(&Options{Application: "test"})

	transport, ok := c.httpClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.Equal(t, MaxIdleConnections, transport.MaxIdleConns)
	assert.Equal(t, MaxIdleConnections, transport.MaxIdleConnsPerHost)
	assert.Equal(t, IdleConnTimeout, transport.IdleConnTimeout)
}

func TestRequestsReuseConnections(t *testing.T) {
	var newConns atomic.Int32

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id":"x"}` + "\n"))
	}))
	srv.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			newConns.Add(1)
		}
	}
	srv.Start()
	defer srv.Close()

	c := New(&Options{Application: "test", URL: srv.URL})

	const parallel = 10
	const rounds = 20

	for i := 0; i < rounds; i++ {
		var wg sync.WaitGroup
		for j := 0; j < parallel; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var resp struct {
					ID string `json:"id"`
				}
				assert.NoError(t, c.get("/channels/x", &resp))
				assert.Equal(t, "x", resp.ID)
			}()
		}
		wg.Wait()
	}

	assert.LessOrEqual(t, int(newConns.Load()), parallel)
}
