package paginator

var (
	defLimit   = int64(100)
	defLastId  = ""
	defOffset  = int64(0)
	defOrderBy = string("createdAt")
	defSort    = string("asc")
)

type Query struct {
	Limit   int64
	LastId  string
	Offset  int64
	OrderBy string
	Sort    string
}

func Init(p *Query) *Query {
	page := p

	if p == nil {
		return &Query{
			Limit:   defLimit,
			OrderBy: defOrderBy,
			Sort:    defSort,
		}
	}

	if page.Limit == 0 {
		page.Limit = defLimit
	}

	if page.LastId == "" {
		page.LastId = defLastId
	}

	if page.Offset == 0 {
		page.Offset = defOffset
	}

	if page.OrderBy == "" {
		page.OrderBy = defOrderBy
	}

	if page.Sort == "" {
		page.Sort = defSort
	}

	return page
}
