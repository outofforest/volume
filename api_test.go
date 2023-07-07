package volume

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/outofforest/logger"
	"github.com/ridge/must"
	"github.com/ridge/parallel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/outofforest/volume/lib/libnet"
)

func env(t *testing.T, testFn func(url string)) {
	ctx, cancel := context.WithCancel(logger.WithLogger(context.Background(), logger.New(logger.DefaultConfig)))
	t.Cleanup(cancel)

	l := libnet.ListenOnRandomPort()
	defer func() {
		_ = l.Close()
	}()

	group := parallel.NewGroup(ctx)
	group.Spawn("server", parallel.Fail, func(ctx context.Context) error {
		return Run(ctx, l)
	})

	testFn("http://" + l.Addr().String() + "/reduce")
}

func TestValidRequest(t *testing.T) {
	env(t, func(url string) {
		req := must.HTTPRequest(http.NewRequest(http.MethodPost, url, bytes.NewReader(must.Bytes(json.Marshal([][]string{
			{"BBB", "CCC"},
			{"AAA", "BBB"},
			{"CCC", "DDD"},
		})))))
		resp := must.HTTPResponse(http.DefaultClient.Do(req))

		require.Equal(t, http.StatusOK, resp.StatusCode)
		body := must.Bytes(ioutil.ReadAll(resp.Body))

		result := []string{}
		must.OK(json.Unmarshal(body, &result))

		assert.Equal(t, []string{"AAA", "DDD"}, result)
	})
}

func TestInvalidRequest1(t *testing.T) {
	env(t, func(url string) {
		req := must.HTTPRequest(http.NewRequest(http.MethodPost, url, bytes.NewReader(must.Bytes(json.Marshal([][]string{
			{"BBB"},
			{"AAA", "BBB"},
		})))))
		resp := must.HTTPResponse(http.DefaultClient.Do(req))

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestInvalidRequest2(t *testing.T) {
	env(t, func(url string) {
		req := must.HTTPRequest(http.NewRequest(http.MethodPost, url, bytes.NewReader(must.Bytes(json.Marshal([][]string{
			{"BBB", "CCC", "DDD"},
			{"AAA", "BBB"},
		})))))
		resp := must.HTTPResponse(http.DefaultClient.Do(req))

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestInvalidRequest3(t *testing.T) {
	env(t, func(url string) {
		req := must.HTTPRequest(http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("invaliddata"))))
		resp := must.HTTPResponse(http.DefaultClient.Do(req))

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestNotSolvable(t *testing.T) {
	env(t, func(url string) {
		req := must.HTTPRequest(http.NewRequest(http.MethodPost, url, bytes.NewReader(must.Bytes(json.Marshal([][]string{
			{"BBB", "CCC"},
			{"AAA", "BBB"},
			{"DDD", "EEE"},
		})))))
		resp := must.HTTPResponse(http.DefaultClient.Do(req))

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
