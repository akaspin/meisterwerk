package app

import (
	"context"

	"github.com/akaspin/meisterwerk/api/gen/client/orders"
	"github.com/akaspin/meisterwerk/api/gen/server/quotes"
)

const (
	QuoteStatusPending  = "pending"
	QuoteStatusAccepted = "accepted"
	QuoteStatusRejected = "rejected"
)

type QuoteReporter interface {
	ReportQuotes(context.Context, []quotes.Quote) error
}

type QuoteReporterFunc func(context.Context, []quotes.Quote) error

func (f QuoteReporterFunc) ReportQuotes(ctx context.Context, qs []quotes.Quote) error {
	return f(ctx, qs)
}

func MockQuoteReporter(ordersAPI *orders.APIClient) QuoteReporter {
	return QuoteReporterFunc(func(ctx context.Context, qs []quotes.Quote) error {
		var ids []int64
		for _, q := range qs {
			if q.Status == QuoteStatusAccepted {
				ids = append(ids, q.Id)
			}
		}
		_, err := ordersAPI.DefaultAPI.ReportQuotes(ctx).RequestBody(ids).Execute()
		return err
	})
}
