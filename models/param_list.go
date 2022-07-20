package models

// ParamList :
type ParamList struct {
	Page       int    `json:"page" form:"page"  valid:"Required"`
	PerPage    int    `json:"perpage" form:"perpage"  valid:"Required"`
	Search     string `json:"search,omitempty"`
	InitSearch string `json:"init_search,omitempty"`
	SortField  string `json:"sort_field,omitempty"`
}

type ParamDynamicList struct {
	ParamList
	MenuUrl   string `json:"menu_url" valid:"Required"`
	LineNo    int    `json:"line_no,omitempty"`
	ParamView string `json:"param_view,omitempty"`
}
