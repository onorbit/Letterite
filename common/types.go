package common

type Page struct {
	ID           int64  `json:"ID"`
	ParentPageID int64  `json:"parentPageID"`
	Subject      string `json:"subject"`
}

type PageSummary struct {
	Page
	ChildrenPageCount int `json:"childrenCount"`
	ContentCount      int `json:"contentCount"`
}
