// Copyright 2015 HORISEN AG. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apiutil

import (
	"net/http"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// CollectionMeta holds CollectionMetaPagination data
type CollectionMeta struct {
	Pagination CollectionMetaPagination `json:"pagination"`
}

// CollectionMetaPagination holds pagination data
type CollectionMetaPagination struct {
	Total       int `json:"total"`
	TotalPages  int `json:"totalPages"`
	CurrentPage int `json:"currentPage"`
	PerPage     int `json:"perPage"`
	Count       int `json:"count"`
}

// DataCollection hold
type DataCollection struct {
	Data []interface{}  `json:"data"`
	Meta CollectionMeta `json:"meta"`
}

//DataCollection hold exchange data

// calculates total pages based od total records and page size
func totalPages(totRec, pageSize int) int {
	if pageSize == 0 {
		return 1
	}
	if totRec%pageSize != 0 {
		return totRec/pageSize + 1
	}
	return totRec / pageSize
}

func maxInt(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// AdjustListReturnedPage returns current page for a list of records with given requestedPage, totalRecords and number of records per page.
// In case requested page is greather than total number of pages , last page is returned.
func AdjustListReturnedPage(requestedPage, totalRecords, perPage int) int {
	return maxInt(1, minInt(requestedPage, totalPages(totalRecords, perPage)))
}

// NewDataCollection returns DataCollection object with pagination metadata
func NewDataCollection(l []interface{}, totRec int, p *Pagination) (*DataCollection, error) {
	if !p.ReturnAll() {
		return &DataCollection{
			Data: l,
			Meta: CollectionMeta{
				Pagination: CollectionMetaPagination{
					Total:       totRec,
					TotalPages:  totalPages(totRec, p.PageSize),
					Count:       len(l),
					PerPage:     p.PageSize,
					CurrentPage: p.PageNumber,
				},
			},
		}, nil
	}
	return &DataCollection{
		Data: l,
		Meta: CollectionMeta{
			Pagination: CollectionMetaPagination{
				Total:       totRec,
				TotalPages:  1,
				Count:       totRec,
				PerPage:     totRec,
				CurrentPage: 1,
			},
		},
	}, nil
}

// Pagination is common structure that should be included as anonymous class in PathRegistryHandler implementation.
// Pagination properties will be set automatically on HTTP request based on request's query parameters using mappings and defult values described in struct tag query.
// For example : ...?customerName=CoolSMS&paging=1,10&sort=-invoiceData,invoiceAmount.
type Pagination struct {
	All             string `query:"all" default:"no"`
	PageNumber      int    `query:"page_number" default:"1"`
	PageSize        int    `query:"page_size" default:"10"`
	Paging          string `query:"paging"`
	Sort            string `query:"sort"`
	SortOptionsList []*SortOptions
}

// Parse checks paging values for errors and sort format for error
func (p *Pagination) Parse() error {
	if p.Paging != "" {
		// expecting format `page_num,page_size`
		t := strings.Split(p.Paging, ",")
		if len(t) == 2 {
			pageNum, err := strconv.Atoi(t[0])
			if err != nil {
				return err
			}
			p.PageNumber = pageNum

			pageSize, err := strconv.Atoi(t[1])
			if err != nil {
				return err
			}

			p.PageSize = pageSize
		} else {
			return ErrBadParameter
		}
	}

	if p.Sort != "" {
		t := strings.Split(p.Sort, ",")
		for _, e := range t {
			if strings.HasPrefix(e, "-") {
				p.SortOptionsList = append(p.SortOptionsList, &SortOptions{
					Property:  strings.TrimPrefix(e, "-"),
					Direction: "DESC",
				})
			} else {
				p.SortOptionsList = append(p.SortOptionsList, &SortOptions{
					Property: e,
				})
			}
		}
	}

	if !p.ReturnAll() {
		if p.PageNumber < 1 {
			return ErrBadPageNumber
		}
		if p.PageSize < 1 {
			return ErrBadPageSize
		}
	}

	return nil
}

// ReturnAll pareses All property value and returns true or false
func (p *Pagination) ReturnAll() bool {
	switch p.All {
	case "1", "t", "T", "true", "TRUE", "True", "yes", "Yes", "YES", "y", "Y":
		return true
	case "0", "f", "F", "false", "FALSE", "False", "no", "No", "NO", "n", "N":
		return false
	default:
		return false
	}
}

// SortOptions holt sort parameter and direction. Direction can be DESC in  case of descending sorting and
// ASC or blank value in case of ascending sorting.
type SortOptions struct {
	Property  string
	Direction string
}

// URLQueryString is data structure for holding key, operator and data that form url query
type URLQueryString struct {
	Key      string
	Operator string
	Values   []string
}

// CreateQueryString is function that creates query string for URl
func CreateQueryString(urlQuery []*URLQueryString) string {
	var queryString string
	for _, value := range urlQuery {
		switch value.Operator {
		case "eq", "isnull", "isnotnull", "isempty", "isnotempty", "neq", "startswith", "contains", "endswith", "doesnotcontain":
			if value.Key == "language" || value.Key == "all" || value.Key == "page_size" || value.Key == "page_number" || value.Key == "q" {
				queryString += value.Key + "=" + value.Values[0] + "&"
			} else {
				queryString += value.Operator + "(" + value.Key + ")=" + value.Values[0] + "&"
			}
		case "lt", "lte", "gt", "gte":
			queryString += value.Operator + "(" + value.Key + ")=" + value.Values[0] + "&"
		case "between", "notbetween":
			queryString += value.Operator + "(" + value.Key + ")=" + value.Values[0] + "," + value.Values[1] + "&"
		case "in", "notin":
			valueString := ""
			for _, v := range value.Values {
				valueString += v + ","
			}
			valueString = strings.TrimRight(valueString, ",")
			queryString += value.Operator + "(" + value.Key + ")=" + valueString + "&"

		case "":
			queryString += value.Key + "=" + value.Values[0] + "&"
		}
	}
	queryString = strings.TrimRight(queryString, "&")
	return queryString
}

// ReplaceSubstr special characters with _
func ReplaceSubstr(input string) string {
	slc := []string{"\\", ":", "?", "]", "[", "*", "/"}

	res1 := input

	for _, value := range slc {
		res1 = strings.Replace(res1, value, "_", -1)
	}
	return res1

}

func ReplaceSubstrBackslash(input string) string {
	slc := []string{"/"}

	res1 := input

	for _, value := range slc {
		res1 = strings.Replace(res1, value, "_", -1)
	}
	return res1

}

// CustomSubstr
func CustomSubstr(input string, maxLength int) string {
	out := input
	if len(input) > maxLength {
		out = input[:maxLength]
	}
	fixUtf := func(r rune) rune {
		if r == utf8.RuneError {
			return -1
		}
		return r
	}

	return strings.Map(fixUtf, out)
}

// ValidateAndCleanPhoneNumber validates and cleans phone number for inserting to database
func ValidateAndCleanPhoneNumber(phoneNumber string) (string, *Error) {
	phone := strings.Replace(phoneNumber, " ", "", -1)

	if strings.HasPrefix(phone, "+") {
		if strings.HasPrefix(phone, "+0") {
			return "", &Error{HTTPCode: http.StatusBadRequest,
				Code:    "INVALID PHONE NUMBER",
				Message: "A given phone number cannot contain 0 after +"}
		}
		phone = strings.TrimLeft(phone, "+")

	} else if strings.HasPrefix(phone, "00") {
		if strings.HasPrefix(phone, "000") {
			return "", &Error{HTTPCode: http.StatusBadRequest,
				Code:    "INVALID PHONE NUMBER",
				Message: "A given number cannot contain 0 on third place"}
		}
		phone = strings.TrimLeft(phone, "00")
		// } else if strings.HasPrefix(phone, "99") {
		// 	if strings.HasPrefix(phone, "990") {
		// 		return "", &Error{HTTPCode: http.StatusBadRequest,
		// 			Code:    "INVALID PHONE NUMBER",
		// 			Message: "A given number cannot contain 0 on third place"}
		// 	}
		// 	phone = strings.TrimLeft(phone, "99")
	} else if strings.HasPrefix(phone, "0") {
		return "", &Error{HTTPCode: http.StatusBadRequest,
			Code:    "INVALID PHONE NUMBER",
			Message: "Domestic phone number is not supported"}
	}

	if len(phone) < 8 {
		return "", &Error{HTTPCode: http.StatusBadRequest,
			Code:    "INVALID PHONE NUMBER",
			Message: "The input is less than 8 characters long"}
	}

	if !isInt(phone) {
		return "", &Error{HTTPCode: http.StatusBadRequest,
			Code:    "INVALID PHONE NUMBER",
			Message: "Phone number can contain only digits"}
	}

	return phone, nil
}

// ValidateAndCleanPhoneNumber validates and cleans phone number od seve digits for inserting to database
func ValidateAndCleanPhoneNumberSevenDigits(phoneNumber string) (string, *Error) {
	phone := strings.Replace(phoneNumber, " ", "", -1)

	if strings.HasPrefix(phone, "+") {
		if strings.HasPrefix(phone, "+0") {
			return "", &Error{HTTPCode: http.StatusBadRequest,
				Code:    "INVALID PHONE NUMBER",
				Message: "A given phone number cannot contain 0 after +"}
		}
		phone = strings.TrimLeft(phone, "+")

	} else if strings.HasPrefix(phone, "00") {
		if strings.HasPrefix(phone, "000") {
			return "", &Error{HTTPCode: http.StatusBadRequest,
				Code:    "INVALID PHONE NUMBER",
				Message: "A given number cannot contain 0 on third place"}
		}
		phone = strings.TrimLeft(phone, "00")
		// } else if strings.HasPrefix(phone, "99") {
		// 	if strings.HasPrefix(phone, "990") {
		// 		return "", &Error{HTTPCode: http.StatusBadRequest,
		// 			Code:    "INVALID PHONE NUMBER",
		// 			Message: "A given number cannot contain 0 on third place"}
		// 	}
		// 	phone = strings.TrimLeft(phone, "99")
	} else if strings.HasPrefix(phone, "0") {
		return "", &Error{HTTPCode: http.StatusBadRequest,
			Code:    "INVALID PHONE NUMBER",
			Message: "Domestic phone number is not supported"}
	}

	if len(phone) < 7 {
		return "", &Error{HTTPCode: http.StatusBadRequest,
			Code:    "INVALID PHONE NUMBER",
			Message: "The input is less than 8 characters long"}
	}

	if !isInt(phone) {
		return "", &Error{HTTPCode: http.StatusBadRequest,
			Code:    "INVALID PHONE NUMBER",
			Message: "Phone number can contain only digits"}
	}

	return phone, nil
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func isUnsignedInteger(s string) bool {

	value, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	if value <= 0 {
		return false
	}
	return true

}

// PaginationHelper is common structure that should be included as anonymous class in PathRegistryHandler implementation.
// PaginationHelper properties will be set automatically on HTTP request based on request's query parameters using mappings and defult values described in struct tag query.
// For example : ...?customerName=CoolSMS&paging=1,10&sort=-invoiceData,invoiceAmount.
type PaginationHelper struct {
	All             string `query:"all" default:"no"`
	PageNumber      int    `query:"page_number" default:"1"`
	PageSize        string `query:"page_size"`
	Paging          string `query:"paging"`
	Sort            string `query:"sort"`
	PageSizeInteger int
	SortOptionsList []*SortOptions
}

// Parse checks paging values for errors and sort format for error
func (p *PaginationHelper) Parse() error {
	if p.Paging != "" {
		// expecting format `page_num,page_size`
		t := strings.Split(p.Paging, ",")
		if len(t) == 2 {
			//page_num
			pageNum, err := strconv.Atoi(t[0])
			if err != nil {
				return err
			}
			p.PageNumber = pageNum

			//page_size
			p.PageSize = t[1]
		} else {
			return ErrBadParameter
		}
	}

	if p.Sort != "" {
		t := strings.Split(p.Sort, ",")
		for _, e := range t {
			if strings.HasPrefix(e, "-") {
				p.SortOptionsList = append(p.SortOptionsList, &SortOptions{
					Property:  strings.TrimPrefix(e, "-"),
					Direction: "DESC",
				})
			} else {
				p.SortOptionsList = append(p.SortOptionsList, &SortOptions{
					Property: e,
				})
			}
		}
	}

	if !p.ReturnAll() {
		if p.PageNumber < 1 {
			return ErrBadPageNumber
		}

		if p.PageSize == "" || p.PageSize == "0" {
			p.All = "yes"
		} else if isUnsignedInteger(p.PageSize) {
			psNum, err := strconv.Atoi(p.PageSize)
			if err != nil {
				return ErrBadPageSize
			}
			p.PageSizeInteger = psNum
		} else {
			return ErrBadPageSize
		}

	}

	return nil
}

// ReturnAll pareses All property value and returns true or false
func (p *PaginationHelper) ReturnAll() bool {
	switch p.All {
	case "1", "t", "T", "true", "TRUE", "True", "yes", "Yes", "YES", "y", "Y":
		return true
	case "0", "f", "F", "false", "FALSE", "False", "no", "No", "NO", "n", "N":
		return false
	default:
		return false
	}
}

// NewDataCollectionHelper returns DataCollection object with pagination metadata
func NewDataCollectionHelper(l []interface{}, totRec int, p *PaginationHelper) (*DataCollection, error) {
	if !p.ReturnAll() {
		return &DataCollection{
			Data: l,
			Meta: CollectionMeta{
				Pagination: CollectionMetaPagination{
					Total:       totRec,
					TotalPages:  totalPages(totRec, p.PageSizeInteger),
					Count:       len(l),
					PerPage:     p.PageSizeInteger,
					CurrentPage: p.PageNumber,
				},
			},
		}, nil
	}
	return &DataCollection{
		Data: l,
		Meta: CollectionMeta{
			Pagination: CollectionMetaPagination{
				Total:       totRec,
				TotalPages:  1,
				Count:       totRec,
				PerPage:     totRec,
				CurrentPage: 1,
			},
		},
	}, nil
}

// CopyFilterMap copies a given filter map to a new filter map
func CopyFilterMap(inputFilter map[FilterKey][]string) map[FilterKey][]string {
	outputFilter := make(map[FilterKey][]string)
	for k, v := range inputFilter {
		outputFilter[k] = v
	}
	return outputFilter
}

// DeleteElementsFromFilterMapByProperty deletes elements that have given property in a key
func DeleteElementsFromFilterMapByProperty(property string, inputFilter map[FilterKey][]string) {
	for k := range inputFilter {
		if k.Property == property {
			delete(inputFilter, k)
		}
	}
	return
}
