package subapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

type requestsMaker struct {
	httpClient sling.Doer
	sling      *sling.Sling
}

func newRequestsMakerWithSling(httpClient sling.Doer, s *sling.Sling) *requestsMaker {
	return &requestsMaker{
		httpClient: httpClient,
		sling:      s,
	}
}

func (r *requestsMaker) getResponseBytes(
	path string, queryParams interface{}) (out []byte, err error) {
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

	if response.StatusCode == http.StatusOK {
		return bodyBytes, nil
	}

	if response.StatusCode == http.StatusNotFound {
		return []byte("Resource is not found"), nil
	}

	return nil, r.restError(bodyBytes, response.StatusCode, "GET")
}

func (r *requestsMaker) getJSON(path string, out interface{}) error {
	request, _ := r.sling.New().Get(path).Request()
	response, err := r.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		err = json.NewDecoder(response.Body).Decode(out)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return r.restError(b, response.StatusCode, "GET")
}

func (r *requestsMaker) post(path string, data interface{}, out interface{}) error {
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
	return r.restError(dt, response.StatusCode, "POST")
}

func (r *requestsMaker) delete(path string) error {
	return r.deleteByIDWithSling(r.sling, path)
}

func (r *requestsMaker) deleteByIDWithSling(sling *sling.Sling, resourceID string) error {
	request, _ := sling.New().Delete(resourceID).Request()
	response, err := r.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode == http.StatusNoContent {
		return nil
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		dt, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return r.restError(dt, response.StatusCode, "DELETE")
	}

	return nil
}

func (r *requestsMaker) restError(dt []byte, status int, op string) error {
	return errors.Wrapf(
		errAPI,
		"API error, status '%d' method '%s': %s", status, op, string(dt),
	)
}
