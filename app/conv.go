package app

import (
	"github.com/akaspin/meisterwerk/api/gen/server/quotes"
	"github.com/akaspin/meisterwerk/model"
)

// QuoteToModel converts Quote to model
func QuoteToModel(q quotes.Quote) *model.Quote {
	res := &model.Quote{
		ID:          uint64(q.Id),
		CustomerID:  uint64(q.CustomerId),
		Description: q.Description,
		Status:      q.Status,
	}
	for _, item := range q.Items {
		res.Items = append(res.Items, &model.Item{
			QuoteID: uint64(q.Id),
			ItemID:  item.Id,
			Segment: item.Segment,
			Price:   item.Price,
			Tax:     item.Tax,
		})
	}
	return res
}

// QuotesSliceToModel converts slice of Quotes to model
func QuotesSliceToModel(qs []quotes.Quote) model.Quotes {
	res := make([]*model.Quote, len(qs))
	for i, q := range qs {
		res[i] = QuoteToModel(q)
	}
	return res
}

// ModelToQuote converts model to quote
func ModelToQuote(q *model.Quote) quotes.Quote {
	res := quotes.Quote{
		Id:          int64(q.ID),
		CustomerId:  int64(q.CustomerID),
		Description: q.Description,
		Status:      q.Status,
	}
	for _, item := range q.Items {
		res.Items = append(res.Items, quotes.Item{
			Id:      item.ItemID,
			Segment: item.Segment,
			Price:   item.Price,
			Tax:     item.Tax,
		})
	}
	return res
}

// ModelSliceToQuote converts model to slice of quotes
func ModelSliceToQuote(m model.Quotes) []quotes.Quote {
	var res []quotes.Quote
	for _, q := range m {
		res = append(res, ModelToQuote(q))
	}
	return res
}
