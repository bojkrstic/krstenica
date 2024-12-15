package api

import (
	"krstenica/pkg/apiutil"
	"net/http"
)

type HramList struct {
	apiutil.Pagination
	Filters map[apiutil.FilterKey][]string `query:"data-filters"`
}

func (ac *HramList) Handle(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	//get hram from db, to create return json
	lists, totalCount, err := db.ListHram(
		ac.Pagination.ReturnAll(),
		ac.PageNumber,
		ac.PageSize,
		ac.Pagination.SortOptionsList,
		ac.Filters)
	if err != nil {
		return nil, err
	}
	res := make([]interface{}, len(lists))
	for i, list := range lists {
		res[i] = makeResultSysApplication(list)
	}

	return apiutil.NewDataCollection(res, totalCount, &ac.Pagination)
}
