package volume

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/outofforest/logger"
	"github.com/outofforest/volume/lib/libhttp"
	"github.com/ridge/must"
	"go.uber.org/zap"
)

// Run starts API server
func Run(ctx context.Context, listener net.Listener) error {
	log := logger.Get(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/reduce", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Error", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		input := [][]string{}
		if err := json.Unmarshal(body, &input); err != nil {
			log.Error("Error", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hops := make([]Hop, 0, len(input))
		for _, i := range input {
			if len(i) != 2 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			hops = append(hops, Hop{Start: i[0], End: i[1]})
		}
		result, err := Reduce(hops)
		if err != nil {
			log.Error("Error", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(must.Bytes(json.Marshal([]string{result.Start, result.End}))); err != nil {
			log.Error("Error", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	return libhttp.Run(ctx, listener, mux)
}
