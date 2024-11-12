package main

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/akaspin/meisterwerk/api"
	"github.com/akaspin/meisterwerk/api/gen/server/orders"
	"github.com/spf13/pflag"
)

type OrdersAPI struct {
	accepted map[int64]struct{}
	mu       sync.RWMutex
}

func NewOrdersAPI() *OrdersAPI {
	return &OrdersAPI{
		accepted: map[int64]struct{}{},
	}
}

func (o *OrdersAPI) ReportQuotes(ctx context.Context, ids []int64) (orders.ImplResponse, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, id := range ids {
		o.accepted[id] = struct{}{}
	}
	return orders.Response(http.StatusCreated, nil), nil
}

func (o *OrdersAPI) GetQuotes(ctx context.Context) (orders.ImplResponse, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	ids := slices.Collect(maps.Keys(o.accepted))

	return orders.Response(http.StatusOK, ids), nil
}

func main() {
	var port int
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.IntVar(&port, "port", 8080, "bind port")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		logger.Error("bad args", "error", err.Error())
		os.Exit(2)
	}

	service := NewOrdersAPI()
	h := orders.NewRouter(orders.NewDefaultAPIController(service))
	api.WithHealthcheck(h)

	err = api.ServeHTTP(logger, &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           h,
		ReadHeaderTimeout: time.Second * 5,
	})
	if err != nil {
		logger.Error("bad exit", "error", err)
		os.Exit(2)
	}
}
