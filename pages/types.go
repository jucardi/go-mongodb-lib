package pages

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Page encapsulates the essential information required to request a subset of a result set.
type Page struct {
	Page int      `json:"page"`           // Page number to be requested.
	Size int      `json:"size"`           // The page size of the subset.
	Sort []string `json:"sort,omitempty"` // The fields to use for a sorting algorithm. Use '-' at the beginning for reverse order. Eg "-name".
}

// Paginated result containing a subset of the result set. JSON keys were done to match the names used by Ten-X Java commons library.
type Paginated struct {
	*PaginatedBase
	Items interface{} `json:"content"` // The array of items in the result.
}

// PaginatedBase contains the base fields for the paginated wrapper. It was separated from Paginated so it can easily
// be used to created a deserialization struct when the array type is know. For example:
//
//	type PaginatedUsers struct {
//		*PaginatedBase
//		Users []*User `json:"content"`
//	}
//
type PaginatedBase struct {
	ItemsCount int `json:"count"`        // The total amount of elements in this subset.
	TotalPages int `json:"total_pages"`  // The total amount of pages by the provided page size.
	TotalCount int `json:"total_count"`  // The total amount of items in the query.
	Size       int `json:"size"`         // The page size
	Page       int `json:"current_page"` // The page number this subset represents.
}

// CreateFromContext creates a Page retrieving the query strings passed in a request from the gin.Context.
func CreateFromContext(c *gin.Context, defaultPage ...*Page) (ret *Page) {
	if len(defaultPage) > 0 {
		ret = defaultPage[0]
	}

	page, pageExists := getIntQuery(c, "page")
	size, sizeExists := getIntQuery(c, "size")

	if !pageExists || !sizeExists {
		return
	}

	field := c.QueryArray("sort_field")

	return &Page{
		Page: page,
		Size: size,
		Sort: field,
	}
}

// CreatePaginated creates the paginated object based on the given page and result.
func CreatePaginated(p *Page, array interface{}, count ...int) (*Paginated, error) {
	resultv := reflect.ValueOf(array)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice && resultv.Kind() != reflect.Array {
		return nil, fmt.Errorf("unable to create Paginated, 'array' arg must be a Slice or Array, %v", resultv.Kind())
	}

	arr := resultv.Elem()
	l := arr.Len()
	c := l

	if len(count) > 0 {
		c = count[0]
	}

	if p == nil {
		p = &Page{1, l, nil}
	}

	return &Paginated{
		Items: arr.Interface(),
		PaginatedBase: &PaginatedBase{
			ItemsCount: l,
			TotalPages: int(math.Ceil(float64(c) / float64(p.Size))),
			TotalCount: c,
			Size:       p.Size,
			Page:       p.Page,
		},
	}, nil
}

func getIntQuery(c *gin.Context, key string) (int, bool) {
	if str, e := c.GetQuery(key); !e {
		return 0, false
	} else {
		v, err := strconv.Atoi(str)
		return v, err == nil
	}
}
