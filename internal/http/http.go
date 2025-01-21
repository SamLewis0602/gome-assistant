// http is used to interact with the home assistant
// REST API. Currently only used to retrieve state for
// a single entity_id
package http

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

type HttpClient struct {
	url   string
	token string
}

func NewHttpClient(url *url.URL, token string) *HttpClient {
	// Shallow copy the URL to avoid modifying the original
	u := *url
	u.Path = "/api"
	if u.Scheme == "ws" {
		u.Scheme = "http"
	}
	if u.Scheme == "wss" {
		u.Scheme = "https"
	}

	return &HttpClient{
		url:   u.String(),
		token: token,
	}
}

func (c *HttpClient) GetState(entityId string) ([]byte, error) {
	resp, err := get(c.url+"/states/"+entityId, c.token)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *HttpClient) States() ([]byte, error) {
	resp, err := get(c.url+"/states", c.token)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func get(url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("Error creating HTTP request: " + err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error on response.\n[ERROR] -" + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Error while reading the response bytes:" + err.Error())
	}

	return body, nil
}

// func post(url string, token string, data any) ([]byte, error) {
// 	postBody, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req, err := http.NewRequest("GET", url, bytes.NewBuffer(postBody))
// 	if err != nil {
// 		return nil, errors.New("Error building post request: " + err.Error())
// 	}

// 	req.Header.Add("Authorization", "Bearer "+token)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, errors.New("Error in post response: " + err.Error())
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == 401 {
// 		panic("ERROR: Auth token is invalid. Please double check it or create a new token in your Home Assistant profile")
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return body, nil
// }
