package mgo

import (
	"fmt"

	"github.com/jucardi/go-mongodb-lib/pages"
)

// IQueryPageExtension encapsulates the new extended functions to the original IQuery
type IQueryPageExtension interface {
	// Page adds to the query the information required to fetch the requested page of objects.
	Page(p ...*pages.Page) IQuery

	// WrapPage attempts to obtain the items in the requested page and wraps the result in *pages.Paginated
	WrapPage(result interface{}, p ...*pages.Page) (*pages.Paginated, error)
}

func (q *query) Page(page ...*pages.Page) IQuery {
	return pageHandler(q, page...)
}

func (q *query) WrapPage(result interface{}, page ...*pages.Page) (*pages.Paginated, error) {
	return wrapPageHandler(q, result, page...)
}

func pageHandler(q IQuery, page ...*pages.Page) IQuery {
	if len(page) < 1 || page[0] == nil {
		return q
	}

	p := page[0]

	if len(p.Sort) > 0 {
		q.Sort(p.Sort...)
	}
	return q.Skip((p.Page - 1) * p.Size).Limit(p.Size)
}

func wrapPageHandler(q IQuery, result interface{}, page ...*pages.Page) (*pages.Paginated, error) {
	if len(page) < 1 || page[0] == nil {
		return wrap(q, result, nil)
	}

	n, err := q.Count()

	if err != nil {
		return nil, fmt.Errorf("unable to obtain a count of elements, %v", err)
	}

	p := page[0]
	q.Page(page...)

	return wrap(q, result, p, n)
}

func wrap(q IQuery, result interface{}, p *pages.Page, n ...int) (*pages.Paginated, error) {
	if err := q.All(result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal page, %v", err)
	}

	return pages.CreatePaginated(p, result, n...)
}
