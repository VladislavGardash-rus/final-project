package middleware

import (
	"context"
	"encoding/json"
	"github.com/gardashvs/final-project/internal/logger"
	"net/http"
	"strconv"
	"time"
)

type errorMessage struct {
	Method string `json:"method"`
	Api    string `json:"api"`
	Error  string `json:"error"`
}

type HandlerFunc func(ctx context.Context, r *http.Request) (interface{}, error)

func Serve(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := new(http.Response)
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)

		data, err := h(ctx, r)
		if err != nil {
			logger.UseLogger().Error(err)

			response.StatusCode = http.StatusBadRequest
			r.Response = response

			errorMessage := new(errorMessage)
			errorMessage.Method = r.Method
			errorMessage.Api = r.RequestURI
			errorMessage.Error = err.Error()
			b, _ := json.Marshal(errorMessage)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(b)

			return
		}

		if file, ok := data.([]byte); ok {
			w.Header().Add("Content-Type", http.DetectContentType(file))
			w.Header().Add("Content-Length", strconv.Itoa(len(file)))

			response.StatusCode = http.StatusOK
			r.Response = response
			w.WriteHeader(http.StatusOK)
			w.Write(file)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		response.StatusCode = http.StatusOK
		r.Response = response
		b, _ := json.Marshal(data)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
