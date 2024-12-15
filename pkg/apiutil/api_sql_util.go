// Copyright 2016 HORISEN AG. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apiutil

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

/*

Possible operators:

	Value Operators
	isnull - Is Null
	isnotnull - Is not Null
	isempty - Is Empty
	isnotempty - Is not Empty
	isnullorempty - Is Null or Empty
	isnotnullorempty - Is not Null or Empty
	String Operators
	eq - Is Equal To
	neq - Not Equals To
	startswith - Starts With
	contains - Contains
	endswith - Ends With
	doesnotcontain - Does Not Contain
	Numeric Operators
	eq - Is Equal To
	neq - Not Equals To
	lt - Less Than
	lte - Less Than or Equal
	gte - Greater Than or Equal
	gt - Greater Than
	Range Operators
	in - In array of values
	notin - Not In array of values
	between - Between values
	notbetween - Not Between values
*/

// SortPropertyTransformer is function responsible for transforming attribute name to database table attribute name.
// It should take attribute name from api filter and return attribute name in database.
// For example : bizPartnerId -> biz_partner_id
// This function should be written for each endpoint individually(list endpoints) as it is highly dependant on context of its execution
// like database table name, query etc...
// in some cases it is not simple camelCase to underscorecase transformation but also database attribute name
// will contain schema and table prefix ...
// When function return "",nil it means that given property is valid but not supported in current context
type SortPropertyTransformer func(attributeName string) (string, error)

// FilterPropertyValidator is function responsible for transforming attribute name to database table attribute name and to validate attributa value
// It should take attribute name from api filter and return attribute name from database.
// For example : bizPartnerId -> biz_partner_id
// This function should be written for each endpoint individually(list endpoints) as it is highly dependant on context of its execution
// like database table name, query etc...
// in some cases it is not simple camelCase to underscorecase transformation but also resulting database attribute name
// will contain schema and table prefix ...
// Attribute value should be validated for correctness such as type, value, format etc...
// This requires deep knowledge about possible attribute value , type, format and it  mostly
// depends on context.
// when function return "", nil  it means thar current attribute is not supported in current context and should be ignored
type FilterPropertyValidator func(attributeName string, attributeValue []string) (string, error)

// defaultTransform is AttributeTransformer implementation which only returns the same attribute
func defaultTransform(v string) (string, error) {
	return v, nil
}

// defaultTransform is AttributeTransformer implementation which only returns the same attribute
func defaultValidate(v string, val []string) (string, error) {
	return v, nil
}

// ErrUnsupportedFilterOperator is returned when unknown operator is detected in filter
var ErrUnsupportedFilterOperator = NewError(http.StatusBadRequest,
	"REQUEST_ERROR", "unsupported filter operator")

// ErrUnsupportedFilterProperty is returned when filter property with given name is not supported
var ErrUnsupportedFilterProperty = NewError(http.StatusBadRequest,
	"REQUEST_ERROR", "unsupported filter property")

// ErrUnsupportedSortProperty is returned when sort property with given name is not supported
var ErrUnsupportedSortProperty = NewError(http.StatusBadRequest,
	"REQUEST_ERROR", "unsupported sort property")

// NewInvalidFilterPropertyValueError should be used by AttributeValidator when attribute value is not valid
func NewInvalidFilterPropertyValueError(propertyName, message string) *Error {
	return NewError(http.StatusBadRequest,
		"REQUEST_ERROR", fmt.Sprintf("query parameter: '%s' error: %s", propertyName, message))
}

// FilterValidationError implements Error interface and shod be returned by FilterPropertyValidator function
// in case filter property value is bad
type FilterValidationError struct {
	Item    string
	Message string
}

// NewFilterValidationError is FilterValidationError constructor
func NewFilterValidationError(property, errorMessage string) *FilterValidationError {
	return &FilterValidationError{
		Item:    property,
		Message: errorMessage,
	}
}

func (f *FilterValidationError) Error() string {
	return fmt.Sprintf("filter property '%s' error : %s", f.Item, f.Message)
}

func stringValueToInterface(s string) []interface{} {
	return []interface{}{s}
}

func stringSliceToInterface(s []string) []interface{} {
	par := []interface{}{}
	for _, e := range s {
		par = append(par, e)
	}
	return par
}

func filterValueParser(op string, value []string) ([]string, error) {
	switch op {
	case "isnull", "isnotnull", "isempty", "isnotempty", "isnullorempty", "isnotnullorempty":
		return nil, nil
	case "eq", "neq", "lt", "lte", "gt", "gte", "startswith", "endswith", "contains", "doesnotcontain":
		if len(value) != 1 {
			return nil, ErrBadParameter
		}
		return value[:1], nil
	case "in", "notin":
		if len(value) == 0 {
			return nil, ErrBadParameter
		}
		par := []string{}
		for _, e := range value {
			sl := strings.Split(e, ",")
			for _, e := range sl {
				par = append(par, e)
			}
		}
		return par, nil
	case "between", "notbetween":
		if len(value) != 1 {
			return nil, ErrBadParameter
		}
		vals := strings.Split(value[0], ",")
		if len(vals) != 2 {
			return nil, ErrBadParameter
		}
		return []string{vals[0], vals[1]}, nil
	default:
		return nil, ErrUnsupportedFilterOperator
	}
}

func filterToSQL(attr, op string, value []string, fn FilterPropertyValidator) (string, []interface{}, error) {

	val, err := filterValueParser(op, value)
	if err != nil {
		if err == ErrUnsupportedFilterProperty {
			return "", nil, nil
		}
		return "", nil, err
	}

	//log.Println(op, value, val)

	attribute, err := fn(attr, val)
	if err != nil {
		// when attribute us unsupported, skip , ignore unsupported attributes
		if err == ErrUnsupportedFilterProperty {
			return "", nil, nil
		}
		return "", nil, err
	}

	if attribute == "" {
		return "", nil, nil
	}
	switch op {
	case "isnull":
		return fmt.Sprintf("%s IS NULL", attribute), nil, nil
	case "isnotnull":
		return fmt.Sprintf("%s IS NOT NULL", attribute), nil, nil
	case "isempty":
		return fmt.Sprintf("%s = ''", attribute), nil, nil
	case "isnotempty":
		return fmt.Sprintf("%s <> ''", attribute), nil, nil
	case "isnullorempty":
		return fmt.Sprintf("(%s IS NULL OR %s = '') ", attribute, attribute), nil, nil
	case "isnotnullorempty":
		return fmt.Sprintf("(%s IS NOT NULL OR %s = '') ", attribute, attribute), nil, nil
	case "eq":
		return fmt.Sprintf("%s = ?", attribute), stringValueToInterface(val[0]), nil
	case "neq":
		return fmt.Sprintf("%s <> ?", attribute), stringValueToInterface(val[0]), nil
	case "startswith":
		return fmt.Sprintf("%s LIKE ?", attribute), stringValueToInterface(LikeStringValidation(val[0]) + "%"), nil
	case "contains":
		st := stringValueToInterface("%" + LikeStringValidation(val[0]) + "%")
		return fmt.Sprintf("%s LIKE ?", attribute), st, nil
	case "endswith":
		return fmt.Sprintf("%s LIKE ?", attribute), stringValueToInterface("%" + LikeStringValidation(val[0])), nil
	case "doesnotcontain":
		return fmt.Sprintf("%s NOT LIKE ?", attribute), stringValueToInterface("%" + LikeStringValidation(val[0]) + "%"), nil
	case "lt":
		return fmt.Sprintf("%s < ?", attribute), stringValueToInterface(val[0]), nil
	case "lte":
		return fmt.Sprintf("%s <= ?", attribute), stringValueToInterface(val[0]), nil
	case "gt":
		return fmt.Sprintf("%s > ?", attribute), stringValueToInterface(val[0]), nil
	case "gte":
		return fmt.Sprintf("%s >= ?", attribute), stringValueToInterface(val[0]), nil
	case "in":
		return fmt.Sprintf("%s IN  (", attribute) + strings.TrimSuffix(strings.Repeat("?,", len(val)), ",") + ")", stringSliceToInterface(val), nil
	case "notin":
		return fmt.Sprintf("%s NOT IN  (", attribute) + strings.TrimSuffix(strings.Repeat("?,", len(val)), ",") + ")", stringSliceToInterface(val), nil
	case "between":
		return fmt.Sprintf("%s BETWEEN ? AND ?", attribute), stringSliceToInterface(val), nil
	case "notbetween":
		return fmt.Sprintf("%s NOT BETWEEN ? AND ?", attribute), stringSliceToInterface(val), nil
	default:
		return "", nil, nil
	}
}

// FilterToSQLWhere returns SQL where clause with parameter values slice suitable for embedding in SQL query.
func FilterToSQLWhere(filters map[FilterKey][]string, fn FilterPropertyValidator) (string, []interface{}, error) {
	var where []string
	params := []interface{}{}
	log.Printf("filters %v", filters)
	for k, val := range filters {
		log.Printf("val  %v %v", k, val)
		w, p, err := filterToSQL(k.Property, k.Operator, val, fn)
		if err != nil {
			return "", nil, err
		}
		if w != "" {
			where = append(where, w)
			if p != nil {
				params = append(params, p...)
			}
		}
	}

	return strings.Join(where, " AND "), params, nil
}

func sortToSQL(attr *SortOptions, fn SortPropertyTransformer) (string, error) {
	property, err := fn(attr.Property)
	if err != nil {
		// when attribute us unsupported, ignore it
		if err == ErrUnsupportedSortProperty {
			return "", nil
		}
		return "", err
	}

	if property == "" {
		return "", nil
	}

	return property + " " + attr.Direction, nil

}

// SortSQL returns SQL Order BY clause based on list of sort options suitable to query embeding
func SortSQL(sorting []*SortOptions, fn SortPropertyTransformer) (string, error) {
	var sorts []string
	//params := []interface{}{}

	for _, val := range sorting {
		s, err := sortToSQL(val, fn)
		if err != nil {
			return "", err
		}
		if s != "" {
			sorts = append(sorts, s)
		}
	}

	return strings.Join(sorts, ","), nil
}

// CheckStringValid3BytesUTF8 DB uses utf-8 with max 3 bytes, this function will check if string has
// utf-8 character with more then 3 bytes, and replace it with '?'
func CheckStringValid3BytesUTF8(str string) string {
	rstring := []rune(str)
	for _, r := range rstring {
		s := string(r)
		b := []byte(s)
		if len(b) > 3 {
			str = strings.Replace(str, s, "?", -1)
		}
	}
	return str
}

// FilterToSQLWhereWithCollation returns SQL where clause with parameter values slice suitable for embedding in SQL query.
// with collation adding if passed to
func FilterToSQLWhereWithCollation(filters map[FilterKey][]string, fn, fnCollationCheck FilterPropertyValidator) (string, []interface{}, error) {
	var where []string
	params := []interface{}{}
	log.Printf("filters %v", filters)
	for k, val := range filters {
		log.Printf("val  %v %v", k, val)
		w, p, err := filterToSQLWithCollationCheck(k.Property, k.Operator, val, fn, fnCollationCheck)
		if err != nil {
			return "", nil, err
		}
		if w != "" {
			where = append(where, w)
			if p != nil {
				params = append(params, p...)
			}
		}
	}

	return strings.Join(where, " AND "), params, nil
}

func filterToSQLWithCollationCheck(attr, op string, value []string, fn, fnCollationCheck FilterPropertyValidator) (string, []interface{}, error) {

	val, err := filterValueParser(op, value)
	if err != nil {
		if err == ErrUnsupportedFilterProperty {
			return "", nil, nil
		}
		return "", nil, err
	}
	//log.Println(op, value, val)

	attribute, err := fn(attr, val)
	if err != nil {
		// when attribute us unsupported, skip , ignore unsupported attributes
		if err == ErrUnsupportedFilterProperty {

			if fnCollationCheck == nil {
				return "", nil, nil
			}
			attributeCollCheck, err := fnCollationCheck(attr, val)
			if err != nil {
				if err == ErrUnsupportedFilterProperty {
					return "", nil, nil
				}
				return "", nil, err
			}
			if attributeCollCheck == "" {
				return "", nil, nil
			}

			switch op {
			case "isnull":
				return fmt.Sprintf("%s IS NULL", attributeCollCheck), nil, nil
			case "isnotnull":
				return fmt.Sprintf("%s IS NOT NULL", attributeCollCheck), nil, nil
			case "isempty":
				return fmt.Sprintf("%s = ''", attributeCollCheck), nil, nil
			case "isnotempty":
				return fmt.Sprintf("%s <> ''", attributeCollCheck), nil, nil
			case "isnullorempty":
				return fmt.Sprintf("(%s IS NULL OR %s = '') ", attributeCollCheck, attributeCollCheck), nil, nil
			case "isnotnullorempty":
				return fmt.Sprintf("(%s IS NOT NULL OR %s = '') ", attributeCollCheck, attributeCollCheck), nil, nil
			case "eq":
				return fmt.Sprintf("%s = CONVERT( ? USING LATIN1) COLLATE latin1_general_ci", attributeCollCheck), stringValueToInterface(val[0]), nil
			case "neq":
				return fmt.Sprintf("%s <> CONVERT(? USING LATIN1) COLLATE latin1_general_ci", attributeCollCheck), stringValueToInterface(val[0]), nil
			case "startswith":
				return fmt.Sprintf("%s LIKE CONVERT( ? USING LATIN1) COLLATE latin1_general_ci", attributeCollCheck), stringValueToInterface(val[0] + "%"), nil
			case "contains":
				return fmt.Sprintf("%s LIKE CONVERT( ? USING LATIN1) COLLATE latin1_general_ci", attributeCollCheck), stringValueToInterface("%" + val[0] + "%"), nil
			case "endswith":
				return fmt.Sprintf("%s LIKE CONVERT( ? USING LATIN1) COLLATE latin1_general_ci", attributeCollCheck), stringValueToInterface("%" + val[0]), nil
			case "doesnotcontain":
				return fmt.Sprintf("%s NOT LIKE CONVERT( ? USING LATIN1) COLLATE latin1_general_ci", attributeCollCheck), stringValueToInterface("%" + val[0] + "%"), nil
			case "lt":
				return fmt.Sprintf("%s < CONVERT( ? USING LATIN1) COLLATE latin1_general_ci", attributeCollCheck), stringValueToInterface(val[0]), nil
			case "lte":
				return fmt.Sprintf("%s <= CONVERT( ? USING LATIN1) COLLATE latin1_general_ci", attributeCollCheck), stringValueToInterface(val[0]), nil
			case "gt":
				return fmt.Sprintf("%s > CONVERT( ? USING LATIN1) COLLATE latin1_general_ci", attributeCollCheck), stringValueToInterface(val[0]), nil
			case "gte":
				return fmt.Sprintf("%s >= CONVERT( ? USING LATIN1) COLLATE latin1_general_ci", attributeCollCheck), stringValueToInterface(val[0]), nil
			case "in":
				return fmt.Sprintf("%s IN  (", attributeCollCheck) +
					strings.TrimSuffix(strings.Repeat("CONVERT( ? USING LATIN1) COLLATE latin1_general_ci, ", len(val)), ", ") + ")", stringSliceToInterface(val), nil
			case "notin":
				return fmt.Sprintf("%s NOT IN  (", attributeCollCheck) +
					strings.TrimSuffix(strings.Repeat("CONVERT( ? USING LATIN1) COLLATE latin1_general_ci, ", len(val)), ", ") + ")", stringSliceToInterface(val), nil
			case "between":
				return fmt.Sprintf("%s BETWEEN CONVERT( ? USING LATIN1) COLLATE latin1_general_ci AND CONVERT( ? USING LATIN1) COLLATE latin1_general_ci",
					attributeCollCheck), stringSliceToInterface(val), nil
			case "notbetween":
				return fmt.Sprintf("%s NOT BETWEEN CONVERT( ? USING LATIN1) COLLATE latin1_general_ci AND CONVERT( ? USING LATIN1) COLLATE latin1_general_ci",
					attributeCollCheck), stringSliceToInterface(val), nil
			default:
				return "", nil, nil
			}

		}
		return "", nil, err
	}

	if attribute == "" {
		return "", nil, nil
	}
	switch op {
	case "isnull":
		return fmt.Sprintf("%s IS NULL", attribute), nil, nil
	case "isnotnull":
		return fmt.Sprintf("%s IS NOT NULL", attribute), nil, nil
	case "isempty":
		return fmt.Sprintf("%s = ''", attribute), nil, nil
	case "isnotempty":
		return fmt.Sprintf("%s <> ''", attribute), nil, nil
	case "isnullorempty":
		return fmt.Sprintf("(%s IS NULL OR %s = '') ", attribute, attribute), nil, nil
	case "isnotnullorempty":
		return fmt.Sprintf("(%s IS NOT NULL OR %s = '') ", attribute, attribute), nil, nil
	case "eq":
		return fmt.Sprintf("%s = ?", attribute), stringValueToInterface(val[0]), nil
	case "neq":
		return fmt.Sprintf("%s <> ?", attribute), stringValueToInterface(val[0]), nil
	case "startswith":
		return fmt.Sprintf("%s LIKE ?", attribute), stringValueToInterface(val[0] + "%"), nil
	case "contains":
		return fmt.Sprintf("%s LIKE ?", attribute), stringValueToInterface("%" + val[0] + "%"), nil
	case "endswith":
		return fmt.Sprintf("%s LIKE ?", attribute), stringValueToInterface("%" + val[0]), nil
	case "doesnotcontain":
		return fmt.Sprintf("%s NOT LIKE ?", attribute), stringValueToInterface("%" + val[0] + "%"), nil
	case "lt":
		return fmt.Sprintf("%s < ?", attribute), stringValueToInterface(val[0]), nil
	case "lte":
		return fmt.Sprintf("%s <= ?", attribute), stringValueToInterface(val[0]), nil
	case "gt":
		return fmt.Sprintf("%s > ?", attribute), stringValueToInterface(val[0]), nil
	case "gte":
		return fmt.Sprintf("%s >= ?", attribute), stringValueToInterface(val[0]), nil
	case "in":
		return fmt.Sprintf("%s IN  (", attribute) + strings.TrimSuffix(strings.Repeat("?,", len(val)), ",") + ")", stringSliceToInterface(val), nil
	case "notin":
		return fmt.Sprintf("%s NOT IN  (", attribute) + strings.TrimSuffix(strings.Repeat("?,", len(val)), ",") + ")", stringSliceToInterface(val), nil
	case "between":
		return fmt.Sprintf("%s BETWEEN ? AND ?", attribute), stringSliceToInterface(val), nil
	case "notbetween":
		return fmt.Sprintf("%s NOT BETWEEN ? AND ?", attribute), stringSliceToInterface(val), nil
	default:
		return "", nil, nil
	}
}

// LikeStringValidation validates the string used for LIKE mysql function, replacing \,_,% with appropriate values
func LikeStringValidation(str string) string {
	//str = strings.ReplaceAll(str, "\\", "\\\\\\\\")
	//	str = strings.ReplaceAll(str, "%", "\\%")
	//str = strings.ReplaceAll(str, "_", "\\_")

	return str
}
