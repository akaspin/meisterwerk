package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/akaspin/meisterwerk/api/gen/client/orders"
	"github.com/akaspin/meisterwerk/api/gen/client/quotes"
	"github.com/akaspin/meisterwerk/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeQuotes(cID int32, name string, n int) []quotes.Quote {
	var res []quotes.Quote
	for i := 0; i < n; i++ {
		q := quotes.Quote{
			CustomerId:  int64(cID),
			Description: fmt.Sprintf("api-quote-%s-%d", name, i),
			Status:      quotes.PtrString("pending"),
		}
		for j := 0; j < 3; j++ {
			q.Items = append(q.Items, quotes.Item{
				Id:      fmt.Sprintf("api-item-%s-%d-%d", name, i, j),
				Segment: "product",
				Price:   quotes.PtrFloat32(100),
				Tax:     quotes.PtrFloat32(0.1),
			})
		}
		res = append(res, q)
	}
	return res
}

//nolint:bodyclose // tests
func TestQuotesAPI(t *testing.T) {
	ctx := context.Background()
	cID := storage.TestPackNextID()
	cli := quotes.NewAPIClient(quotes.NewConfiguration())

	ordersCnf := orders.NewConfiguration()
	ordersCnf.Host = "localhost:8090"
	ordersCli := orders.NewAPIClient(ordersCnf)

	t.Run("create one", func(t *testing.T) {
		q := quotes.Quote{
			CustomerId:  int64(cID),
			Description: "quote-1",
			Items: []quotes.Item{
				{
					Id:      "item-1",
					Segment: "product",
				},
			},
		}
		res, resp, err := cli.DefaultAPI.CreateQuote(ctx).Quote(q).Execute()
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		t.Run("get one", func(t *testing.T) {
			res2, resp2, err2 := cli.DefaultAPI.GetQuote(ctx, int32(*res.Id)).Execute()
			require.NoError(t, err2)
			assert.Equal(t, http.StatusOK, resp2.StatusCode)
			assert.Len(t, res2.Items, 1)
			assert.EqualValues(t, 100, *res2.Items[0].Price)
		})
		t.Run("item not found", func(t *testing.T) {
			q2 := quotes.Quote{
				CustomerId:  int64(cID),
				Description: "quote-1",
				Items: []quotes.Item{
					{
						Id:      "item-failed-1",
						Segment: "product",
					},
				},
			}
			_, _, nfErr := cli.DefaultAPI.CreateQuote(ctx).Quote(q2).Execute()
			require.EqualError(t, nfErr, "400 Bad Request")
		})
		t.Run("duplicate item", func(t *testing.T) {
			q2 := quotes.Quote{
				CustomerId:  int64(cID),
				Description: "quote-1",
				Items: []quotes.Item{
					{
						Id:      "item-dup-1",
						Segment: "product",
					},
					{
						Id:      "item-dup-1",
						Segment: "product",
					},
				},
			}
			_, _, nfErr := cli.DefaultAPI.CreateQuote(ctx).Quote(q2).Execute()
			require.EqualError(t, nfErr, "400 Bad Request")
		})
	})
	t.Run("get not found", func(t *testing.T) {
		_, resp, err := cli.DefaultAPI.GetQuote(ctx, 22000000).Execute()
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("list quotes", func(t *testing.T) {
		cID2 := storage.TestPackNextID()
		var ids []int32
		for i := 0; i < 5; i++ {
			q := quotes.Quote{
				CustomerId:  int64(cID2),
				Description: fmt.Sprintf("quote %d", i),
				Items: []quotes.Item{
					{
						Id:      "product-1",
						Segment: "product",
						Price:   quotes.PtrFloat32(100),
						Tax:     quotes.PtrFloat32(0.1),
					},
				},
			}
			res, resp, err := cli.DefaultAPI.CreateQuote(ctx).Quote(q).Execute()
			require.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
			id := *res.Id
			ids = append(ids, int32(id))
		}
		t.Run("by ids", func(t *testing.T) {
			res, resp, err := cli.DefaultAPI.ListQuotes(ctx).Id(ids).Execute()
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Len(t, res, 5)
		})
	})

	t.Run("update quote", func(t *testing.T) {
		q := quotes.Quote{
			CustomerId:  int64(cID),
			Description: "quote-update-id",
			Items: []quotes.Item{
				{
					Id:      "product-1",
					Segment: "product",
					Price:   quotes.PtrFloat32(100),
					Tax:     quotes.PtrFloat32(0.1),
				},
			},
		}
		res, resp, err := cli.DefaultAPI.CreateQuote(ctx).Quote(q).Execute()
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		q2 := quotes.Quote{
			Id:          res.Id,
			CustomerId:  int64(cID),
			Description: "quote-update-id-2",
			Status:      quotes.PtrString("accepted"),
			Items: []quotes.Item{
				{
					Id:      "product-1",
					Segment: "product",
					Price:   quotes.PtrFloat32(200),
					Tax:     quotes.PtrFloat32(0.1),
				},
			},
		}
		resp2, err2 := cli.DefaultAPI.UpdateQuote(ctx, int32(*res.Id)).Quote(q2).Execute()
		require.NoError(t, err2)
		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		res3, resp3, err3 := cli.DefaultAPI.GetQuote(ctx, int32(*res.Id)).Execute()
		require.NoError(t, err3)
		assert.Equal(t, http.StatusOK, resp3.StatusCode)
		assert.EqualValues(t, "quote-update-id-2", res3.Description) // rest is covered in model tests
	})
	t.Run("delete quote", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			q := quotes.Quote{
				CustomerId:  int64(cID),
				Description: "quote-delete",
			}
			res, resp, err := cli.DefaultAPI.CreateQuote(ctx).Quote(q).Execute()
			require.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
			resp2, err2 := cli.DefaultAPI.DeleteQuote(ctx, int32(*res.Id)).Execute()
			require.NoError(t, err2)
			assert.Equal(t, http.StatusOK, resp2.StatusCode)
		})
		t.Run("not found", func(t *testing.T) {
			_, err2 := cli.DefaultAPI.DeleteQuote(ctx, 9999).Execute()
			require.EqualError(t, err2, "404 Not Found")
		})
	})

	t.Run("bulk create", func(t *testing.T) {
		qs := makeQuotes(cID, "bulk-create", 3)
		res, resp, err := cli.DefaultAPI.CreateQuotes(ctx).Quote(qs).Execute()
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Len(t, res.Ids, 3)
	})
	t.Run("bulk update", func(t *testing.T) {
		res, resp, err := cli.DefaultAPI.CreateQuotes(ctx).Quote(makeQuotes(cID, "bulk-update", 3)).Execute()
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var ids32 []int32
		for _, id := range res.Ids {
			ids32 = append(ids32, int32(id))
		}
		qs2, resp, err := cli.DefaultAPI.ListQuotes(ctx).Id(ids32).Execute()
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		for i := range qs2 {
			qs2[i].Status = quotes.PtrString("accepted")
		}

		resp2, err2 := cli.DefaultAPI.UpdateQuotes(ctx).Quote(qs2).Execute()
		require.NoError(t, err2)
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		qs4, resp, err := cli.DefaultAPI.ListQuotes(ctx).Id(ids32).Execute()
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		for _, q := range qs4 {
			assert.Equal(t, "accepted", *q.Status)
		}
	})
	t.Run("bulk delete", func(t *testing.T) {
		res, resp, err := cli.DefaultAPI.CreateQuotes(ctx).Quote(makeQuotes(storage.TestPackNextID(),
			"bulk-delete", 3)).Execute()
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		resp2, err2 := cli.DefaultAPI.DeleteQuotes(ctx).QuotesIDs(*res).Execute()
		require.NoError(t, err2)
		assert.Equal(t, http.StatusOK, resp2.StatusCode)
		t.Run("not found", func(t *testing.T) {
			_, err3 := cli.DefaultAPI.DeleteQuotes(ctx).QuotesIDs(*res).Execute()
			assert.EqualError(t, err3, "404 Not Found")
		})
	})
	t.Run("accept quotes", func(t *testing.T) {
		q := makeQuotes(cID, "accept", 1)[0]
		res, _, err := cli.DefaultAPI.CreateQuote(ctx).Quote(q).Execute()
		require.NoError(t, err)
		q.Id = res.Id
		q.Status = quotes.PtrString("accepted")
		_, err = cli.DefaultAPI.UpdateQuote(ctx, int32(*q.Id)).Quote(q).Execute()
		require.NoError(t, err)

		oRes, _, oErr := ordersCli.DefaultAPI.GetQuotes(ctx).Execute()
		require.NoError(t, oErr)
		assert.Contains(t, oRes, *q.Id)
	})
}
