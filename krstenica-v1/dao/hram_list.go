package dao

import (
	"krstenica/pkg/apiutil"
	"log"
	"strings"
)

var allowedAtributesInAttributeTypesFilters = []string{
	"hramId", "hramName", "status",
}
var allowedAtributesInAttributeTypesSort = []string{
	"hramId", "hramName", "status",
}

func transformAttributeTypeSortAttribute(p string) (string, error) {
	if !inList(p, allowedAtributesInAttributeTypesSort) {
		return "", apiutil.ErrBadPageNumber
	}

	p = Underscore(p)

	return p, nil
}

func validateAttributeTypeFilterAttr(p string, v []string) (string, error) {
	if !inList(p, allowedAtributesInAttributeTypesFilters) {
		return "", apiutil.ErrUnsupportedFilterProperty
	}

	// // validate type of values
	// fn, ok := filterAttributeTypeValueValidationMap[p]
	// if ok {
	// 	err := validateType(v, fn)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return "", apiutil.NewFilterValidationError(p, err.Error())
	// 	}
	// }

	p = Underscore(p)

	return p, nil
}

func (c *HramDaoPostgresSql) ListHram(all bool, page, count int,
	sort []*apiutil.SortOptions, filter map[apiutil.FilterKey][]string) ([]*HramDo, int, error) {
	c.Connect()
	defer c.Disconect()

	query := `SELECT COUNT(*) FROM public.hram`
	var total int
	err := c.db.QueryRow(query).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	deleteStatus := "deleted"
	var q string
	params := []interface{}{}

	where, whereParams, err := apiutil.FilterToSQLWhereWithCollation(filter, validateAttributeTypeFilterAttr, validateAttributeTypeFilterAttr)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}
	if where != "" {
		q += where
		params = append(params, whereParams...)
	}

	orderBy, err := apiutil.SortSQL(sort, transformAttributeTypeSortAttribute)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}
	if orderBy != "" {
		if !strings.Contains(orderBy, "hramId") {
			orderBy += ", hramId"
		}
	} else {
		orderBy = "hramId"
	}

	var hrams []*HramDo
	rows, err := c.db.Query("select hram_id, naziv_hrama, status, created_at from public.hram where status!=$1 ORDER BY $2", deleteStatus, orderBy)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	defer rows.Close()
	for rows.Next() {
		var hram HramDo
		err = rows.Scan(&hram.HramID, &hram.HramName, &hram.Status, &hram.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, 0, err
		}
		hrams = append(hrams, &hram)
	}
	if len(hrams) == 0 {
		log.Println("No records found for the given page and pageSize")
	}

	return hrams, total, nil

}
