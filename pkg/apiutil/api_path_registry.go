package apiutil

import (
	"encoding/json"
	"fmt"
	"krstenica/pkg/config"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// HTTPMethod contains HTTP method value
type HTTPMethod string

// possible HTTPMethod values
const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	PATCH   HTTPMethod = "PATCH"
	OPTIONS HTTPMethod = "OPTIONS"
)

// PathRegistryHandler handles routing for one route
// Handle method handles one API call.
// Returns: JSON object (may be null), and API error
type PathRegistryHandler interface {
	Handle(w http.ResponseWriter, r *http.Request) (interface{}, error)
}

// PathRegistryItem is configuration of one route
type pathRegistryItem struct {
	pattern string
	method  HTTPMethod
	re      *regexp.Regexp
	handler PathRegistryHandler
}

// PathRegistry is registry of API endpoints
type PathRegistry struct {
	handlers []*pathRegistryItem
	Config   config.APISrvDynConfigurator
	subPath  string
}

// NewPathRegistry creates Path registry using dynamic configuration
func NewPathRegistry(configFilePath string) (*PathRegistry, *config.Config) {
	c := config.NewAPISrvDynConf(configFilePath)

	cfg, _ := c.GetConf()

	pr := &PathRegistry{
		Config: c,
	}

	return pr, cfg
}

// MapWithOptions adds new route with varPattern and including handler options
func (pr *PathRegistry) MapWithOptions(varPattern string, method HTTPMethod,
	handler PathRegistryHandler) {

	re := regexp.MustCompile(`\\$[a-zA-Z9-9]+`)
	pattern := re.ReplaceAllStringFunc(varPattern, func(varname string) string {
		return fmt.Sprintf("(?P<%s>[^\\/]+)", varname[1:])

	}) + "/?"

	pr.add(pattern, method, handler)
}

// Add adds new route with regexp
func (pr *PathRegistry) add(pattern string, method HTTPMethod,
	handler PathRegistryHandler) {

	re := regexp.MustCompile("^" + pattern + "$")
	pri := &pathRegistryItem{
		pattern: pattern,
		method:  method,
		re:      re,
		handler: handler,
	}
	pr.handlers = append(pr.handlers, pri)
}

// Map adds new route with varPattern
// func (pr *PathRegistry) Map(varPattern string, method HTTPMethod,
// 	handler PathRegistryHandler) {

// 	re := regexp.MustCompile(`\\$[a-zA-Z9-9]+`)
// 	pattern := re.ReplaceAllStringFunc(varPattern, func(varname string) string {
// 		return fmt.Sprintf("(?P<%s>[^\\/]+)", varname[1:])

// 	}) //+ "/?"

// 	fmt.Printf("Mapping pattern: %s\n", pattern)

// 	pr.add(pattern, method, handler)
// }

func (pr *PathRegistry) Map(varPattern string, method HTTPMethod, handler PathRegistryHandler) {
	// re := regexp.MustCompile(`\\$[a-zA-Z9-9]+`)
	// re := regexp.MustCompile("\\$[a-zA-Z9-9]+")  - i ovo je ispravan deo ide sa ""
	re := regexp.MustCompile(`\$[a-zA-Z0-9]+`) // Ispravan regex za prepoznavanje dinamičkih parametara
	pattern := re.ReplaceAllStringFunc(varPattern, func(varname string) string {
		return fmt.Sprintf("(?P<%s>[^\\/]+)", varname[1:])
	}) + "/?"

	log.Printf("Mapping pattern: %s\n", pattern) // Dodajte log za mapiranje

	pr.add(pattern, method, handler)
	log.Printf("Added pattern: %s with method: %s\n", pattern, method) // Dodajte log za dodavanje
}

// NewDefaultHTTPSrv creates new http server based on configured parameters
func (pr *PathRegistry) NewDefaultHTTPSrv(configFilePath string) (*http.Server, error) {
	cfg, err := pr.Config.GetConf()
	if err != nil {
		return nil, err
	}

	mux := pr.ServeMux()

	readTimeout := 10
	writeTimeout := 10

	if cfg.HTTPSServiceTimeout != nil {
		if cfg.HTTPSServiceTimeout.ReadTimeout != 0 {
			readTimeout = cfg.HTTPSServiceTimeout.ReadTimeout
			log.Println("setting read timeout to :", readTimeout, "seconds")
		}

		if cfg.HTTPSServiceTimeout.WriteTimeout != 0 {
			writeTimeout = cfg.HTTPSServiceTimeout.WriteTimeout
			log.Println("setting write timeout to :", writeTimeout, "seconds")
		}
	}

	httpsrv := &http.Server{
		Addr:           cfg.Listen,
		Handler:        mux,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return httpsrv, nil
}

// ServeMux is a HTTP server multiplexer
func (pr *PathRegistry) ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pr.handleHTTPRequest(w, r)
	})

	return mux
}

func (pr *PathRegistry) handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	cfg, err := pr.Config.GetConf()
	if err != nil {
		ResponseSendError(w, err, http.StatusInternalServerError, "*")
		return
	}
	jsonResult, httpCode, err := pr.Dispatch(w, r)
	if err != nil {
		if httpCode == 0 {
			httpCode = http.StatusInternalServerError
		}
		ResponseSendError(w, err, httpCode, cfg.AccessControlAllowOrigin)
		return
	}
	if httpCode == 0 {
		httpCode = http.StatusOK
	}
	if jsonResult != nil {
		SendJSONResponse(w, httpCode, jsonResult, cfg.AccessControlAllowOrigin)
		return
	}

	// the following code has no effect if handler already written headers, status and body
	// if handler return nil and has not written body data or set status , this code is in effect..
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	w.Header().Add("Content-length", "0")
	w.Header().Add("Access-Control-Allow-Origin", cfg.AccessControlAllowOrigin)
	w.WriteHeader(httpCode)
}

func timeTrack(start time.Time, method, name string) {
	elapsed := time.Since(start)
	log.Printf("%s  %s took %s\n", method, name, elapsed)
}

// Dispatch dispatches HTTP request using path registry
func (pr *PathRegistry) Dispatch(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	defer timeTrack(time.Now(), r.Method, r.URL.String())
	cfg, err := pr.Config.GetConf()
	if err != nil {

		return nil, http.StatusInternalServerError, err
	}

	method := HTTPMethod(strings.ToUpper(r.Method))

	if method == OPTIONS {
		return pr.handleOptionsMethod(w, r, cfg.AccessControlAllowOrigin)
	}

	handler, err := pr.findHandler(r)
	if err != nil {
		if err == ErrCannotMatchPath {
			return nil, http.StatusNotFound, err
		}

		return nil, http.StatusInternalServerError, err
	}

	if handler != nil {

		jsonResponse, apiErr := handler.Handle(w, r)
		if apiErr != nil {
			e, ok := apiErr.(*Error)
			if ok {
				return nil, e.HTTPCode, apiErr
			}
			return nil, http.StatusInternalServerError, apiErr
		}
		httpCode := http.StatusOK
		if method == POST {
			httpCode = http.StatusCreated
		}
		if jsonResponse == nil {
			httpCode = http.StatusNoContent
		}
		return jsonResponse, httpCode, nil
	}
	return nil, 0, nil
}

// ResponseSendError send error information encoded in JSON in response body
func ResponseSendError(w http.ResponseWriter, errCause error, httpcode int,
	accessControlAllowOrigin string) {

	var apiErr *ErrorObjWO
	if newAPIErr, ok := errCause.(*Error); ok {
		apiErr = CreateErrorObjWO(newAPIErr)
		// Override API error with the one provided by errcode
		httpcode = newAPIErr.HTTPCode
	} else {
		apiErr = &ErrorObjWO{
			Error: ErrorWO{
				Code:    "SYSTEM_ERROR",
				Message: errCause.Error(),
			},
		}
	}

	js, err := json.Marshal(&apiErr)
	if err != nil {
		log.Printf("Cannot send %s %+v\n", err, apiErr)
		return
	}
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	w.Header().Add("Content-length", strconv.Itoa(len(js)))
	w.Header().Add("Access-Control-Allow-Origin",
		accessControlAllowOrigin)
	w.WriteHeader(httpcode)
	w.Write(js)
}

// SendJSONResponse sends JSON response defined in jsonResp
func SendJSONResponse(w http.ResponseWriter, httpStatusCode int, jsonResp interface{}, accessControlAllowOrigin string) {
	content, err := json.MarshalIndent(jsonResp, "", "\t")
	if err != nil {
		ResponseSendError(w, ErrFailedEncoding, 500, accessControlAllowOrigin)
		return
	}

	contentLen := len(content)
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	w.Header().Add("Content-length", strconv.Itoa(contentLen))
	w.Header().Add("Access-Control-Allow-Origin", accessControlAllowOrigin)
	w.WriteHeader(httpStatusCode)
	_, err = w.Write(content)
	if err != nil {
		log.Printf("Cannot write response %+v\n", err)
	}
}

func (pr *PathRegistry) findHandler(r *http.Request) (PathRegistryHandler, error) {
	method := HTTPMethod(strings.ToUpper(r.Method))

	if pr.handlers != nil {
		for _, pri := range pr.handlers {
			if method == pri.method {
				// ending path with / is not allowed
				if strings.HasSuffix(r.URL.Path, "/") {
					return nil, ErrCannotMatchPath
				}
				// log.Println(pri.pattern, r.URL.Path)
				handler, err := parseURL(r, pri.re, pri.handler)
				if err != nil {
					if err != ErrCannotMatchPath {
						return nil, err
					}
				} else if handler != nil {
					return handler, nil
				}
			}
		}
	}
	return nil, ErrCannotMatchPath
}

func (pr *PathRegistry) handleOptionsMethod(w http.ResponseWriter, r *http.Request, accessControlAllowOrigin string) (interface{}, int, error) {
	var methods []string
	for _, pri := range pr.handlers {
		_, err := parseURL(r, pri.re, pri.handler)
		if err != nil {
			if err != ErrCannotMatchPath {
				log.Printf("Error %+v\n", err)
			}
		} else {
			supportedMethod := string(pri.method)
			methods = append(methods, supportedMethod)
		}
	}
	if methods == nil {
		return nil, http.StatusNotFound, ErrCannotMatchPath
	}
	w.Header().Add("Access-Control-Allow-Methods", strings.Join(methods, ","))
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Headers", "Authorization")
	w.Header().Del("Content-type")
	return nil, 0, nil
}

// ParseURL parses URL from request
func parseURL(r *http.Request, re *regexp.Regexp, obj interface{}) (PathRegistryHandler, error) {

	path := r.URL.Path

	matches := re.FindStringSubmatch(path)
	if matches == nil {
		return nil, ErrCannotMatchPath
	}
	names := re.SubexpNames()

	pathMap := make(map[string]string)
	for i, name := range names {
		pathMap[name] = matches[i]

		if name != "" && matches[i] == "" {
			return nil, ErrCannotMatchPath
		}
	}

	// instantiate new handler object per request using reflection
	handlerVal := reflect.ValueOf(obj).Elem()
	newHandlerVal := reflect.New(handlerVal.Type())
	handler := newHandlerVal.Interface().(PathRegistryHandler)

	objVal := reflect.ValueOf(handler).Elem()
	// return handler and bind HTTP query parameters and path parameters to handler properties
	return handler, setFieldsExt(&objVal, pathMap, r, nil, nil)
}

// GetPrefix returns api path prefix
func (pr *PathRegistry) GetPrefix() (string, error) {
	cfg, err := pr.Config.GetConf()
	if err != nil {
		return "", err
	}
	return cfg.URIPrefix + pr.subPath, nil
}
