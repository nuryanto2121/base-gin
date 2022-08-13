package models

// ParamList :
type ParamList struct {
	Page       int    `json:"page" form:"page"  valid:"Required"`
	PerPage    int    `json:"perpage" form:"perpage"  valid:"Required"`
	Search     string `json:"search" form:"search"`
	InitSearch string `json:"initsearch,omitempty" form:"initsearch"`
	SortField  string `json:"sortfield,omitempty" form:"sortfield"`
}

type ParamDynamicList struct {
	ParamList
	MenuUrl   string `json:"menu_url" valid:"Required"`
	LineNo    int    `json:"line_no,omitempty"`
	ParamView string `json:"param_view,omitempty"`
}
