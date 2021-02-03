package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dghubble/sling"

	"github.com/furiousassault/tc-cli/pkg/configuration"
	"github.com/furiousassault/tc-cli/pkg/teamcity/api/authorization"
	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

const (
	PathSuffixREST = "app/rest/"
	PathSuffixLog  = ""
)

var api *Client

func InitAPI(config configuration.Configuration) error {
	// fmt.Println("configuration", config)
	httpClient := &http.Client{
		Timeout: config.API.HTTP.RequestTimeout,
	}

	if config.API.Authorization.Token != "" {
		newApi, err := NewClientTokenAuth(config.API.URL, httpClient, authorization.AuthToken(config.API.Authorization.Token))
		if err != nil {
			return err
		}

		api = newApi
		return nil
	}

	if config.API.Authorization.Username != "" && config.API.Authorization.Password != "" {
		newApi, err := NewClientBasicAuth(
			config.API.URL,
			httpClient,
			authorization.NewAuthorizerHTTP(config.API.Authorization.Username, config.API.Authorization.Password))

		if err != nil {
			return err
		}

		api = newApi
	}

	if api == nil {
		fmt.Println("No authorization is provided, trying to use guest auth")
	}

	// TODO unreachable now?
	api, err := NewClientGuestAuth(config.API.URL, httpClient, authorization.NewAuthorizerGuest())
	if err != nil {
		return err
	}

	return api.Ping()
}

func API() *Client {
	return api
}

// Client represents the base for connecting to TeamCity
type Client struct {
	address string
	baseURI string

	HTTPClient   *http.Client
	RetryTimeout time.Duration

	commonBase   *sling.Sling
	logFetchBase *sling.Sling

	Projects   *subapi.ProjectService
	BuildTypes *subapi.BuildTypeService
	Builds     *subapi.BuildService
	Logs       *subapi.LogService
	BuildQueue *subapi.BuildQueueService
	Token      *subapi.TokenService
}

func NewClientGuestAuth(address string, httpClient *http.Client, auth authorization.AuthorizerGuest) (*Client, error) {
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	slingBase := sling.New().
		Set("Accept", "application/json").
		Set("Origin", address)

	slingBase.Base(address + auth.ProvideURLAuthSuffix())
	return newClientInstance(address, httpClient, slingBase)
}

func NewClientBasicAuth(address string, httpClient *http.Client, auth authorization.AuthorizerHTTP) (*Client, error) {
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	slingBase := sling.New().
		Set("Accept", "application/json").
		Set("Origin", address)

	slingBase.Base(address+auth.ProvideURLAuthSuffix()).SetBasicAuth(auth.Username, auth.Password)
	return newClientInstance(address, httpClient, slingBase)
}

func NewClientTokenAuth(address string, httpClient *http.Client, auth authorization.AuthorizerToken) (*Client, error) {
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	slingBase := sling.New().
		Set("Accept", "application/json").
		Set("Origin", address)

	slingBase.Base(address+auth.ProvideURLAuthSuffix()).
		Set("Authorization", fmt.Sprintf("Bearer %s", auth.Token))
	return newClientInstance(address, httpClient, slingBase)
}

func newClientInstance(address string, httpClient *http.Client, sling *sling.Sling) (*Client, error) {
	slingRest := sling.New().Path(PathSuffixREST)
	slingLog := sling.New().Path(PathSuffixLog)

	return &Client{
		address:    address,
		HTTPClient: httpClient,
		commonBase: sling,
		Projects:   subapi.NewProjectService(slingRest.New(), httpClient),
		BuildTypes: subapi.NewBuildTypeService(slingRest.New(), httpClient),
		Builds:     subapi.NewBuildService(slingRest.New(), httpClient),
		BuildQueue: subapi.NewBuildQueueService(slingRest.New(), httpClient),
		Token:      subapi.NewTokenService(slingRest.New(), httpClient),

		Logs: subapi.NewLogService(slingLog.New(), httpClient),
	}, nil
}

// Ping tests if the client is properly configured and can be used
func (c *Client) Ping() error {
	response, err := c.commonBase.Get("server").ReceiveSuccess(nil)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 && response.StatusCode != 403 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("API error %s: %s", response.Status, body)
	}

	return nil
}
