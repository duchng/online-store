package paging

type Direction string

const (
	DirectionAsc    Direction = "ASC"
	DirectionDesc   Direction = "DESC"
	DefaultPageSize int       = 20
	MaximumPageSize int       = 200
)

type Order struct {
	Direction  Direction
	ColumnName string
}

type Orders []Order

func (oo *Orders) Contain(columnName string) bool {
	for _, o := range *oo {
		if o.ColumnName == columnName {
			return true
		}
	}
	return false
}

// Add merge given orders with skipping duplicated items
func (oo *Orders) Add(orders ...Order) {
	for i := range orders {
		order := &orders[i]
		if !oo.Contain(order.ColumnName) {
			*oo = append(*oo, *order)
		}
	}
}

func (oo *Orders) Strings() []string {
	res := make([]string, 0, len(*oo))
	for _, o := range *oo {
		res = append(res, o.ColumnName+" "+string(o.Direction))
	}
	return res
}

// Paging request
type Paging struct {
	Sort   Orders
	Size   int
	Cursor int
}

func (p *Paging) Orders() Orders {
	return p.Sort
}

type MetaData struct {
	PageSize    int  `json:"pageSize"`
	HasNext     bool `json:"hasNext"`
	HasPrevious bool `json:"hasPrevious"`
}

type Page[T any] struct {
	Data     []T      `json:"data"`
	Metadata MetaData `json:"metadata"`
}
