package apiutil

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// FilterKey contains property and operator.
// filter is a map[FilterKey][]string
type FilterKey struct {
	Property string
	Operator string
}

func setFieldsExt(objVal *reflect.Value, pathMap map[string]string,
	r *http.Request, reqValues *url.Values, filters *reflect.Value) error {

	if reqValues == nil {
		r.ParseForm()
		reqVal := r.Form
		reqValues = &reqVal
	}

	objTyp := objVal.Type()

	if objTyp.Kind() == reflect.Struct {
		for i := 0; i < objVal.NumField(); i++ {
			fld := objTyp.Field(i)
			fldVal := objVal.Field(i)

			if fldVal.Type().Kind() == reflect.Struct {
				err := setFieldsExt(&fldVal, pathMap, r, reqValues, filters)
				if err != nil {
					return err
				}
			}

			setDefault := true

			tag := fld.Tag.Get("path")
			if tag != "" {
				if matchVal, ok := pathMap[tag]; ok {
					if matchVal != "" {
						if err := setFldValue(&fldVal, matchVal); err != nil {
							// since variable is in path, error here means path is not matched
							// eg. /hello/$id when id is int should be matched when URL is  /hello/1 but not matched wen URL is /hello/abc
							return ErrCannotMatchPath
						}
						setDefault = false
					}
				}
			}

			tag = fld.Tag.Get("query")
			if tag != "" {
				checkSlash := true
				allowValueBeginsWithSlash := fld.Tag.Get("allowValueBeginsWithSlash")
				if allowValueBeginsWithSlash == "yes" {
					checkSlash = false
				}
				// query filters is special tag ... this tag represents map of filters and its values and operators
				// there should be only one property of this kind per struct...
				// if there are more , first of a kind will contain all filters
				if tag != "data-filters" {
					val := r.FormValue(tag)
					if val != "" {
						if checkSlash {
							if strings.HasPrefix(val, "/") {
								// return ErrParameterValueStartsWithSlash
								return ErrCannotMatchPath
							}
						}
						setFldValue(&fldVal, val)
						setDefault = false
						// remove from values
						reqValues.Del(tag)
						if filters != nil {
							// if previously added to filters map , remove it
							k := FilterKey{
								Property: tag,
								Operator: "eq",
							}
							filters.SetMapIndex(reflect.ValueOf(k), reflect.Value{})
						}
					}
				} else {
					// put request values to filters-map
					if fldVal.Kind() == reflect.Map {
						m := make(map[FilterKey][]string)
						if fldVal.IsNil() {
							filters = &fldVal
							filters.Set(reflect.ValueOf(m))
						}
						// filter is is in format operator(attribute)=value or attribute=value (when operator is eq)
						for k, v := range *reqValues {
							opatt := strings.Split(k, "(")
							if len(opatt) == 1 {
								if checkSlash {
									for _, e := range v {
										if strings.HasPrefix(e, "/") {
											// return ErrParameterValueStartsWithSlash
											return ErrCannotMatchPath
										}
									}
								}
								setFilterMapValue(m, opatt[0], "eq", v)
							}
							if len(opatt) == 2 {
								attr := strings.TrimSuffix(opatt[1], ")")
								setFilterMapValue(m, attr, opatt[0], v)
								if checkSlash {
									for _, e := range v {
										if strings.HasPrefix(e, "/") {
											// return ErrParameterValueStartsWithSlash
											return ErrCannotMatchPath
										}
									}
								}
							}
						}

						filters = &fldVal
						filters.Set(reflect.ValueOf(m))
					} else {
						val := r.FormValue(tag)
						if val != "" {
							if checkSlash {
								if strings.HasPrefix(val, "/") {
									// return ErrParameterValueStartsWithSlash
									return ErrCannotMatchPath
								}
							}
							setFldValue(&fldVal, val)
							setDefault = false
							// remove from values
							reqValues.Del(tag)
						}
					}
				}
			}

			if setDefault {
				tag = fld.Tag.Get("default")
				if tag != "" {
					setFldValue(&fldVal, tag)
				}
			}
		}
	}
	return nil
}

func setFields(objVal *reflect.Value, pathMap map[string]string,
	r *http.Request) error {

	objTyp := objVal.Type()

	if objTyp.Kind() == reflect.Struct {
		for i := 0; i < objVal.NumField(); i++ {
			fld := objTyp.Field(i)
			fldVal := objVal.Field(i)

			if fldVal.Type().Kind() == reflect.Struct {
				err := setFields(&fldVal, pathMap, r)
				if err != nil {
					return err
				}
			}

			setDefault := true

			tag := fld.Tag.Get("path")
			if tag != "" {
				if matchVal, ok := pathMap[tag]; ok {
					if matchVal != "" {
						setFldValue(&fldVal, matchVal)
						setDefault = false
					}
				}
			}

			tag = fld.Tag.Get("query")
			if tag != "" {
				val := r.FormValue(tag)
				if val != "" {
					setFldValue(&fldVal, val)
					setDefault = false
				}
			}

			if setDefault {
				tag = fld.Tag.Get("default")
				if tag != "" {
					setFldValue(&fldVal, tag)
				}
			}
		}
	}
	return nil
}

func setFilterMapValue(m map[FilterKey][]string, attr, oper string, val []string) {
	k := FilterKey{attr, oper}
	if _, ok := m[k]; ok {
		m[k] = append(m[k], val...)
		return
	}
	m[k] = val
}

func setFldValue(v *reflect.Value, strVal string) error {
	switch v.Kind() {
	case reflect.Uint:
		iv, err := strconv.Atoi(strVal)
		if err != nil {
			return err
		}
		v.SetUint(uint64(iv))
	case reflect.String:
		v.SetString(strVal)
	case reflect.Int:
		iv, err := strconv.Atoi(strVal)
		if err != nil {
			return err
		}
		v.SetInt(int64(iv))
	case reflect.Int64:
		iv, err := strconv.Atoi(strVal)
		if err != nil {
			return err
		}
		v.SetInt(int64(iv))
	case reflect.Ptr:
		val := reflect.New(v.Type().Elem())
		p := val.Elem()
		setFldValue(&p, strVal)
		v.Set(val)
	default:
		return ErrBadParameter
	}
	return nil
}

// GetRequestBody decodes JSON in body to obj
func GetRequestBody(r *http.Request, obj interface{}) error {
	respBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBodyBytes, obj)
	if err != nil {
		return err
	}
	return nil
}
