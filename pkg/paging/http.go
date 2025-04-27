package paging

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func ParseRequestWithKeysetPagination[T any](ctx echo.Context) (T, Paging, error) {
	var req T
	if err := ctx.Bind(&req); err != nil {
		return req, Paging{}, err
	}
	pagingRequest, err := ParseKeysetPagingRequest(ctx)
	if err != nil {
		return req, Paging{}, err
	}
	return req, pagingRequest, nil
}

func ParseKeysetPagingRequest(ctx echo.Context) (Paging, error) {
	cursor, _ := strconv.Atoi(ctx.QueryParam("cursor"))
	pageSize, _ := strconv.Atoi(ctx.QueryParam("pageSize"))
	sort := ctx.QueryParam("sort")
	res := Paging{
		Size:   pageSize,
		Cursor: cursor,
	}
	if res.Size == 0 {
		res.Size = DefaultPageSize
	}
	if res.Size > MaximumPageSize {
		res.Size = MaximumPageSize
	}
	if sortQuery := sort; sortQuery != "" {
		sort := strings.Split(sortQuery, ",")
		for _, str := range sort {
			if strings.HasPrefix(str, "-") {
				if len(str) == 1 {
					continue
				}
				res.Sort.Add(Order{Direction: DirectionDesc, ColumnName: str[1:]})
			} else if strings.HasPrefix(str, "+") {
				if len(str) == 1 {
					continue
				}
				res.Sort.Add(Order{Direction: DirectionAsc, ColumnName: str[1:]})
			} else {
				res.Sort.Add(Order{Direction: DirectionAsc, ColumnName: str})
			}
		}
	}
	return res, nil
}
