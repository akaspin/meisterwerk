package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/akaspin/meisterwerk/api/gen/server/quotes"
	"github.com/akaspin/meisterwerk/app"
	"github.com/akaspin/meisterwerk/model"
	"github.com/akaspin/meisterwerk/storage"
	"gorm.io/gorm"
)

func Cause(err error) slog.Attr {
	return slog.Any("cause", err)
}

type QuotesAPI struct {
	db             *storage.Conn
	logger         *slog.Logger
	itemsProcessor app.ItemsProcessor
	quoteReporter  app.QuoteReporter
}

func NewQuotesAPI(logger *slog.Logger, db *storage.Conn, catalogCli app.ItemsProcessor,
	quoteReporter app.QuoteReporter) *QuotesAPI {
	return &QuotesAPI{
		logger:         logger,
		db:             db,
		itemsProcessor: catalogCli,
		quoteReporter:  quoteReporter,
	}
}

func (a *QuotesAPI) CreateQuote(_ context.Context, quote quotes.Quote) (quotes.ImplResponse, error) {
	err := a.itemsProcessor.ProcessItems(quote.Items)
	if err != nil {
		a.logger.Error("bad items in quote", Cause(err))
		return quotes.Response(http.StatusBadRequest, nil), err
	}
	id, err := app.QuoteToModel(quote).Create(a.db.GORM)
	if err != nil {
		a.logger.Error("could not create quote", Cause(err))
		return quotes.Response(http.StatusInternalServerError, nil), err
	}
	return quotes.Response(http.StatusCreated, quotes.CreateQuote201Response{
		Id: int64(id),
	}), nil
}

func (a *QuotesAPI) ListQuotes(_ context.Context, ids, customerIDs []int32, skip, limit int32,
	order string) (quotes.ImplResponse, error) {
	criteria := &model.Criteria{
		IDs:         toUint64(ids),
		CustomerIDs: toUint64(customerIDs),
		Limit:       int(limit),
		Skip:        int(skip),
	}
	if order != "desc" {
		criteria.Desc = true
	}
	m, err := criteria.Find(a.db.GORM)
	if err != nil {
		a.logger.Error("could not find quotes", Cause(err))
		return quotes.Response(http.StatusInternalServerError, nil), err
	}
	return quotes.Response(http.StatusOK, app.ModelSliceToQuote(m)), err
}

func (a *QuotesAPI) GetQuote(_ context.Context, id int32) (quotes.ImplResponse, error) {
	q := &model.Quote{ID: uint64(id)}
	err := q.Find(a.db.GORM)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return quotes.Response(http.StatusNotFound, nil), nil
		}
		a.logger.Error("could not retrieve quote", "id", id, Cause(err))
		return quotes.Response(http.StatusInternalServerError, nil), err
	}
	return quotes.Response(http.StatusOK, app.ModelToQuote(q)), nil
}

func (a *QuotesAPI) UpdateQuote(ctx context.Context, id int32, quote quotes.Quote) (quotes.ImplResponse, error) {
	q := app.QuoteToModel(quote)
	q.ID = uint64(id)
	err := q.Update(a.db.GORM)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return quotes.Response(http.StatusNotFound, nil), nil
		}
		a.logger.Error("could not update quote", "id", id, Cause(err))
		return quotes.Response(http.StatusInternalServerError, nil), err
	}
	err = a.quoteReporter.ReportQuotes(ctx, []quotes.Quote{quote})
	if err != nil {
		a.logger.Error("could not report quote to orders", "id", id, Cause(err))
	}
	return quotes.Response(http.StatusOK, nil), nil
}

func (a *QuotesAPI) DeleteQuote(_ context.Context, id int32) (quotes.ImplResponse, error) {
	err := (&model.Quote{ID: uint64(id)}).Delete(a.db.GORM)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return quotes.Response(404, nil), nil
		}
		a.logger.Error("could not delete quote", "id", id, Cause(err))
		return quotes.Response(http.StatusInternalServerError, nil), err
	}
	return quotes.Response(http.StatusOK, nil), nil
}

func (a *QuotesAPI) CreateQuotes(_ context.Context, qs []quotes.Quote) (quotes.ImplResponse, error) {
	for i := range qs {
		err := a.itemsProcessor.ProcessItems(qs[i].Items)
		if err != nil {
			a.logger.Error("bad items in quote", Cause(err))
			return quotes.Response(http.StatusBadRequest, nil), err
		}
	}
	m := app.QuotesSliceToModel(qs)
	ids, err := m.Create(a.db.GORM)
	if err != nil {
		a.logger.Error("could not create quotes", Cause(err))
		return quotes.Response(http.StatusInternalServerError, nil), err
	}
	var resp quotes.QuotesIds
	for _, id := range ids {
		resp.Ids = append(resp.Ids, int64(id))
	}
	return quotes.Response(http.StatusCreated, &resp), err
}

func (a *QuotesAPI) UpdateQuotes(ctx context.Context, qs []quotes.Quote) (quotes.ImplResponse, error) {
	m := app.QuotesSliceToModel(qs)
	err := m.Update(a.db.GORM)
	if err != nil {
		a.logger.Error("could not update quotes", Cause(err))
		return quotes.Response(http.StatusInternalServerError, nil), err
	}
	var ids []int64
	for _, q := range qs {
		if q.Status == app.QuoteStatusAccepted {
			ids = append(ids, q.Id)
		}
	}
	err = a.quoteReporter.ReportQuotes(ctx, qs)
	if err != nil {
		a.logger.Error("could not report accepted quotes to orders", "id", ids, Cause(err))
	}
	return quotes.Response(http.StatusOK, nil), nil
}

func (a *QuotesAPI) DeleteQuotes(_ context.Context, ids quotes.QuotesIds) (quotes.ImplResponse, error) {
	criteria := &model.Criteria{
		IDs: toUint64(ids.Ids),
	}
	err := criteria.Delete(a.db.GORM)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return quotes.Response(http.StatusNotFound, nil), nil
		}
		a.logger.Error("could not delete quotes", "ids", ids.Ids, Cause(err))
		return quotes.Response(http.StatusInternalServerError, nil), err
	}
	return quotes.Response(http.StatusOK, nil), nil
}

func toUint64[T ~[]S, S int32 | int64](is T) []uint64 {
	res := make([]uint64, len(is))
	for i, s := range is {
		res[i] = uint64(s)
	}
	return res
}
