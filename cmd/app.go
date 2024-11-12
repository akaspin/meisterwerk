package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/akaspin/meisterwerk/api"
	"github.com/akaspin/meisterwerk/api/gen/client/orders"
	"github.com/akaspin/meisterwerk/api/gen/server/quotes"
	"github.com/akaspin/meisterwerk/app"
	"github.com/akaspin/meisterwerk/storage"
	"github.com/spf13/pflag"
)

func Run(args []string) error {
	var dbConfig storage.DBConfig
	var port int
	var ordersHost string

	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	(&dbConfig).FS(fs, "meisterwerk")
	fs.IntVar(&port, "port", 8080, "bind port")
	fs.StringVar(&ordersHost, "orders-host", "orders:8080", "orders service host")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	conn, connErr := storage.Connect(&dbConfig)
	if connErr != nil {
		return connErr
	}
	defer conn.Close()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	ordersAPIConf := orders.NewConfiguration()
	ordersAPIConf.Host = ordersHost
	quoteReporter := app.MockQuoteReporter(orders.NewAPIClient(ordersAPIConf))

	service := api.NewQuotesAPI(logger, conn, app.MockItemsProcessor, quoteReporter)
	h := quotes.NewRouter(quotes.NewDefaultAPIController(service))
	api.WithHealthcheck(h)

	return api.ServeHTTP(logger, &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           h,
		ReadHeaderTimeout: time.Second * 5,
	})
}
