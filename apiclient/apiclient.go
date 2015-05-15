package apiclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	httpClient *http.Client
	BaseURL    *url.URL
}

func NewClient(baseURL string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	client := &Client{httpClient: http.DefaultClient, BaseURL: parsedURL}
	return client, nil
}

func (client *Client) Call(method, endpoint string, data url.Values) (*http.Response, error) {
	request, err := client.NewRequest(method, endpoint, data)
	if err != nil {
		return nil, err
	}

	// This seems heavy handed but as a rule we are closing the connection after
	// GetBody below. This ensures that we are communicating our intentions in
	// the HTTP request.
	request.Close = true

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *Client) NewRequest(method, endpoint string, data url.Values) (*http.Request, error) {
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	fullUrl := client.BaseURL.ResolveReference(endpointUrl)

	request, err := http.NewRequest(method, fullUrl.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	return request, nil
}

func (client *Client) GetBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
