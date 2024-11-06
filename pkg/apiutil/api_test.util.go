package apiutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
)

// PerformApiTest is function usefull for testing PathRegistry
func PerformApiTest(pr *PathRegistry,
	method, url string,
	reqObject interface{}, rspObject interface{},
	expectedError error) error {
	return performApiTest(pr, method, url, "", reqObject, rspObject, expectedError)
}

func performApiTest(pr *PathRegistry,
	method, url, token string,
	reqObject interface{}, rspObject interface{},
	expectedError error) error {

	req := createHTTPRequest("127.0.0.1", method, url, token, reqObject)

	w := httptest.NewRecorder()
	result, _, err := pr.Dispatch(w, req)

	if err != expectedError {
		if expectedError != nil {
			if err == nil {
				return fmt.Errorf("Error got none instead of %s", expectedError)
			}
			return fmt.Errorf("Error got %s instead of %s", err, expectedError)
		}
		return fmt.Errorf("Error got %s but expected no error", err)
	} else if expectedError != nil {
		return nil
	}
	if rspObject == nil {
		if result != nil {
			return fmt.Errorf("Found %#v but expected nothing", result)
		}
	} else {
		if result == nil {
			return fmt.Errorf("Found nothing but expected %#v", rspObject)
		}

		// Check type
		returnedType := reflect.TypeOf(result).Elem()
		expectedType := reflect.TypeOf(rspObject).Elem()
		if returnedType != expectedType {
			return fmt.Errorf("Expected type %s but found %s",
				expectedType.Name(), returnedType.Name())
		}

		vset := reflect.ValueOf(rspObject).Elem()
		if !vset.CanSet() {
			return fmt.Errorf("Cannot set %#v", vset)
		}

		vset.Set(reflect.ValueOf(result).Elem())

		if false {
			// Encode response
			b, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("Marshal error %s", err)
			}

			// Decode response into new rspObject
			err = json.Unmarshal(b, rspObject)
			if err != nil {
				return fmt.Errorf("Unmarshal error %s", err)
			}
		}
	}
	return nil
}

func createHTTPRequest(addr, method, uri, token string,
	jsonRequestData interface{}) *http.Request {

	var r *bytes.Reader
	if jsonRequestData != nil {
		b, err := json.Marshal(jsonRequestData)
		if err != nil {
			log.Fatalf("Cannot encode request data %s\n", err)
		}
		log.Printf("req= [%s]\n", string(b))
		r = bytes.NewReader(b)
	} else {
		r = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, "http://"+addr+uri, r)
	if err != nil {
		log.Fatalf("Cannot create request")
	}
	// fool ip address check to pass whitelist validation
	req.RemoteAddr = addr

	err = req.ParseForm()
	if err != nil {
		log.Fatalf("Cannot parse form %s\n", err)
	}

	if token == "" {
		authToken := os.Getenv("AUTHAPI_TOKEN_TEST")
		if authToken != "" {
			req.Header.Add("Authorization", "Bearer "+authToken)
		}
	} else {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	authHost := os.Getenv("TANE_HOST_TEST")
	if authHost != "" {
		req.Host = authHost
	} else {
		req.Host = "api.tane.dev"
	}

	req.Header.Set("Host", req.Host)
	req.Header.Add("X-Forwarded-For", "127.0.0.1")

	return req
}
