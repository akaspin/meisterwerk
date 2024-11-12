package app

import (
	"fmt"
	"strings"

	"github.com/akaspin/meisterwerk/api/gen/server/quotes"
)

// ItemsProcessor wraps method to process Items for quote creation.
type ItemsProcessor interface {
	ProcessItems([]quotes.Item) error
}

type ItemsProcessorFunc func([]quotes.Item) error

func (f ItemsProcessorFunc) ProcessItems(items []quotes.Item) error {
	return f(items)
}

var MockItemsProcessor = ItemsProcessorFunc(func(in []quotes.Item) error {
	processed := map[string]struct{}{}
	for i, v := range in {
		if strings.Contains(v.Id, "failed") {
			return fmt.Errorf("item '%s' not found", v.Id)
		}
		if _, ok := processed[v.Id]; ok {
			return fmt.Errorf("duplicate item '%s'", v.Id)
		}
		processed[v.Id] = struct{}{}

		if in[i].Price == 0 {
			in[i].Price = 100
		}
		if in[i].Tax == 0 {
			in[i].Tax = 0.1
		}
	}
	return nil
})
