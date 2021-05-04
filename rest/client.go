package rest

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"k8s.io/klog/v2"
)

var jar = createJar()

func createJar() *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	return jar
}

type Client interface {
	GET(path string, expectedStatusCode int) ([]byte, error)
}

type client struct {
	config      *Config
	apiBaseURL  string
	credentials Credentials
	httpClient  *http.Client
}

func NewClient(config *Config, apiBaseURL string, credentials Credentials) Client {
	return &debugClient{client: &client{
		config:      config,
		credentials: credentials,
		apiBaseURL:  apiBaseURL,
		httpClient:  createHTTPClient(config),
	}}
}

func createHTTPClient(config *Config) *http.Client {
	var httpClient *http.Client
	if config.NoProxy {
		if config.Insecure {
			httpClient = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					Proxy:           http.ProxyURL(nil)}}
		} else {
			httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(nil)}}
		}
	} else {
		if config.Insecure {
			httpClient = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		} else {
			httpClient = &http.Client{}
		}
	}
	httpClient.Jar = jar
	return httpClient
}

func (c *client) getURL(path string) string {
	apiBaseURL := c.apiBaseURL
	if !strings.HasSuffix(apiBaseURL, "/") {
		apiBaseURL = apiBaseURL + "/"
	}
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	return apiBaseURL + path
}

func (c *client) GET(path string, expectedStatusCode int) ([]byte, error) {
	var err error
	var httpResponse *http.Response
	var request *http.Request

	url := c.getURL(path)
	if c.config.Verbose {
		klog.Info(fmt.Sprintf("GET %s", url))
	}
	if request, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return make([]byte, 0), err
	}
	if err = c.credentials.Authenticate(request); err != nil {
		return make([]byte, 0), err
	}

	if httpResponse, err = c.httpClient.Do(request); err != nil {
		return make([]byte, 0), err
	}
	return readHTTPResponse(httpResponse, http.MethodGet, url, expectedStatusCode, nil, c.config.Debug, c.config.Verbose)
}

func readHTTPResponse(httpResponse *http.Response, method string, url string, expectedStatusCode int, onResponse func(int) error, debug bool, verbose bool) ([]byte, error) {
	var err error
	var body []byte
	defer httpResponse.Body.Close()

	if verbose {
		log.Println(fmt.Sprintf("  %d %s", httpResponse.StatusCode, http.StatusText(httpResponse.StatusCode)))
	}

	if onResponse != nil {
		if err = onResponse(httpResponse.StatusCode); err != nil {
			return nil, err
		}
	}

	if httpResponse.StatusCode != expectedStatusCode {
		finalError := fmt.Errorf("%s (%s) %s", http.StatusText(httpResponse.StatusCode), method, url)
		if body, err = ioutil.ReadAll(httpResponse.Body); err != nil {
			return nil, finalError
		}
		if (body != nil) && len(body) > 0 {
			if debug {
				log.Println("  Response Body: " + string(body))
			}
			var errorEnvelope ErrorEnvelope
			if err := json.Unmarshal(body, &errorEnvelope); err != nil {
				return nil, finalError
			}
			if errorEnvelope.Error != nil {
				return nil, errorEnvelope.Error
			}
		}
		return body, finalError
	}
	if body, err = ioutil.ReadAll(httpResponse.Body); err != nil {
		return nil, err
	}
	if verbose && (body != nil) && len(body) > 0 {
		log.Println("  Response Body: " + string(body))
	}

	return body, nil
}
