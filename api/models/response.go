package models

type Response struct {
	StatusCode  int         `json:"status_code"`
	Description string      `json:"description,omitempty"`
	Count       int         `json:"count"`
	Data        interface{} `json:"data"`
}

type Pagination struct {
	Search
	Limit  int `form:"limit"`
	Page   int `form:"page"`
	Offset int `form:"-" json:"-"`
}

type Search struct {
	Query string `form:"q"`
}

func (p *Pagination) Fix() {
	if p.Limit <= 0 {
		p.Limit = 1000
	}

	if p.Page <= 0 {
		p.Page = 1
	}

	p.Offset = (p.Page - 1) * p.Limit
}
