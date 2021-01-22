package subapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

type requestsMaker struct {
	httpClient *http.Client
	sling      *sling.Sling
}

func newRequestsMakerWithSling(httpClient *http.Client, s *sling.Sling) *requestsMaker {
	return &requestsMaker{
		httpClient: httpClient,
		sling:      s,
	}
}

func (r *requestsMaker) getResponseBytes(path string, queryParams interface{}, resourceDescription string) (out []byte, err error) {
	request, _ := r.sling.New().Get(path).QueryStruct(queryParams).Request()
	response, err := r.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	if response.StatusCode == 200 {
		return bodyBytes, nil
	}

	if response.StatusCode == 404 {
		return []byte("Resource is not found"), nil
	}

	return nil, r.restError(bodyBytes, response.StatusCode, "GET", resourceDescription)
}

func (r *requestsMaker) get(path string, out interface{}, resourceDescription string) error {
	request, _ := r.sling.New().Get(path).Request()
	// fmt.Println(request.URL.Scheme)
	// fmt.Println(request.URL)
	response, err := r.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		err = json.NewDecoder(response.Body).Decode(out)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	dt, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return r.restError(dt, response.StatusCode, "GET", resourceDescription)
}

func (r *requestsMaker) post(path string, data interface{}, out interface{}, resourceDescription string) error {
	request, _ := r.sling.New().Post(path).BodyJSON(data).Request()
	response, err := r.httpClient.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 201 || response.StatusCode == 200 {
		json.NewDecoder(response.Body).Decode(out)
		return nil
	}
	dt, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return r.restError(dt, response.StatusCode, "POST", resourceDescription)
}

func (r *requestsMaker) delete(path string, resourceDescription string) error {
	return r.deleteByIDWithSling(r.sling, path, resourceDescription)
}

func (r *requestsMaker) deleteByIDWithSling(sling *sling.Sling, resourceID string, resourceDescription string) error {
	request, _ := sling.New().Delete(resourceID).Request()
	response, err := r.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode == 204 {
		return nil
	}

	if response.StatusCode != 200 && response.StatusCode != 204 {
		dt, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return r.restError(dt, response.StatusCode, "DELETE", resourceDescription)
	}

	return nil
}

func (r *requestsMaker) restError(dt []byte, status int, op string, res string) error {
	return fmt.Errorf("API error, status '%d' method '%s' operation - %s: %s", status, op, res, string(dt))
}
