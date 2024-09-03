package httplib

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"go-ocr/utils"

	"github.com/gin-gonic/gin"
)

const (
	defaultSize = 10
)

type Query struct {
	SortOrder string `json:"sortOrder,omitempty"`
	OrderBy   string `json:"orderBy,omitempty"`
	Size      int    `json:"size,omitempty"`
	Page      int    `json:"page,omitempty"`
}

func (q *Query) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = n

	return nil
}

func (q *Query) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Page = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

func (q *Query) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

func (q *Query) SetSortOrder(sortOrderByQuery string) {
	arrRules := []string{"asc", "desc"}
	if utils.Contains(arrRules, strings.ToLower(sortOrderByQuery)) {
		if strings.ToLower(sortOrderByQuery) == "asc" {
			sortOrderByQuery = "asc"
		} else if strings.ToLower(sortOrderByQuery) == "desc" {
			sortOrderByQuery = "desc"
		} else {
			sortOrderByQuery = "asc"
		}
	} else {
		sortOrderByQuery = "desc"
	}
	q.SortOrder = sortOrderByQuery
}

func (q *Query) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

func (q *Query) GetLimit() int {
	return q.Size
}

func (q *Query) GetOrderBy() string {
	return q.OrderBy
}

func (q *Query) GetSortOrder() string {
	return q.SortOrder
}

func (q *Query) GetPage() int {
	return q.Page
}

func (q *Query) GetSize() int {
	return q.Size
}

func (q *Query) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s", q.GetPage(), q.GetSize(), q.GetOrderBy())
}

func GetPaginationFromCtx(c *gin.Context) (*Query, error) {
	q := &Query{}
	if err := q.SetPage(c.Query("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(c.Query("size")); err != nil {
		return nil, err
	}
	q.SetOrderBy(c.Query("orderBy"))
	q.SetSortOrder(c.Query("sortOrder"))

	return q, nil
}

func GetTotalPages(totalCount int, pageSize int) int {
	d := float64(totalCount) / float64(pageSize)
	return int(math.Ceil(d))
}
