package thesaurus

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Rican7/define/source"
)

const (
	// baseURLString is the base URL for all Oxford API interactions
	baseURLString = "https://od-api.oxforddictionaries.com/api/v2/"

	entriesURLString = baseURLString + "thesaurus/"

	httpRequestAcceptHeaderName = "Accept"
	httpRequestAppIDHeaderName  = "app_id"
	httpRequestAppKeyHeaderName = "app_key"

	jsonMIMEType = "application/json"
)

type client struct {
	httpClient *http.Client
	appID      string
	appKey     string
}

// apiURL is the URL instance used for Oxford API calls
var apiURL *url.URL

// New returns a new Oxford API dictionary source
func New(httpClient http.Client, appID, appKey string) *client {
	return &client{&httpClient, appID, appKey}
}

// Initialize the package
func init() {
	var err error

	apiURL, err = url.Parse(baseURLString)

	if nil != err {
		panic(err)
	}
}

// Define takes a word string and returns a dictionary source.Result
func (g *client) Define(word string) (*Results, error) {
	// Prepare our URL
	requestURL, err := url.Parse(entriesURLString + "en/" + word)

	if nil != err {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, apiURL.ResolveReference(requestURL).String(), nil)

	if nil != err {
		return nil, err
	}

	httpRequest.Header.Set(httpRequestAcceptHeaderName, jsonMIMEType)
	httpRequest.Header.Set(httpRequestAppIDHeaderName, g.appID)
	httpRequest.Header.Set(httpRequestAppKeyHeaderName, g.appKey)

	httpResponse, err := g.httpClient.Do(httpRequest)

	if nil != err {
		return nil, err
	}

	defer httpResponse.Body.Close()

	if http.StatusNotFound == httpResponse.StatusCode {
		return nil, &source.EmptyResultError{Word: word}
	}

	if http.StatusForbidden == httpResponse.StatusCode {
		return nil, &source.AuthenticationError{}
	}

	body, err := ioutil.ReadAll(httpResponse.Body)

	if nil != err {
		return nil, err
	}

	var result Results

	if err = json.Unmarshal(body, &result); nil != err {
		return nil, err
	}

	if len(result.Results) < 1 {
		return nil, &source.EmptyResultError{Word: word}
	}

	return &result, nil
}
