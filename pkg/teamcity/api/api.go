package api

import (
	"fmt"
	"io/ioutil"

	"github.com/dghubble/sling"

	"github.com/furiousassault/tc-cli/pkg/configuration"
	"github.com/furiousassault/tc-cli/pkg/teamcity/api/authorization"
	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

const (
	PathSuffixREST = "app/rest/"
	PathSuffixLog  = ""
)

func InitAPI(config configuration.Configuration, httpClient sling.Doer) (api *Client, err error) {
	if config.API.Authorization.Token != "" {
		api = NewClientTokenAuth(
			config.API.URL,
			httpClient,
			authorization.NewAuthorizerToken(config.API.Authorization.Token),
		)

		return api, api.Ping()
	}

	if config.API.Authorization.Username != "" && config.API.Authorization.Password != "" {
		api = NewClientBasicAuth(
			config.API.URL,
			httpClient,
			authorization.NewAuthorizerHTTP(config.API.Authorization.Username, config.API.Authorization.Password),
		)

		return api, api.Ping()
	}

	fmt.Println("No authorization is provided, trying to use guest auth")
	api = NewClientGuestAuth(config.API.URL, httpClient, authorization.NewAuthorizerGuest())

	return api, api.Ping()
}

// Client represents the base for connecting to TeamCity
type Client struct {
	address string
	baseURI string

	commonBase   *sling.Sling
	logFetchBase *sling.Sling

	Projects   *subapi.ProjectService
	BuildTypes *subapi.BuildTypeService
	Builds     *subapi.BuildService
	Logs       *subapi.LogService
	BuildQueue *subapi.BuildQueueService
	Token      *subapi.TokenService
}

func NewClientGuestAuth(address string, httpClient sling.Doer, auth authorization.AuthorizerGuest) *Client {
	slingBase := sling.New().
		Set("Accept", "application/json").
		Set("Origin", address)

	slingBase.Base(address + auth.ProvideURLAuthSuffix())
	return newClientInstance(address, httpClient, slingBase)
}

func NewClientBasicAuth(address string, httpClient sling.Doer, auth authorization.AuthorizerHTTP) *Client {
	slingBase := sling.New().
		Set("Accept", "application/json").
		Set("Origin", address)

	slingBase.Base(address+auth.ProvideURLAuthSuffix()).
		SetBasicAuth(auth.Username, auth.Password)

	return newClientInstance(address, httpClient, slingBase)
}

func NewClientTokenAuth(address string, httpClient sling.Doer, auth authorization.AuthorizerToken) *Client {
	slingBase := sling.New().
		Set("Accept", "application/json").
		Set("Origin", address)

	slingBase.Base(address+auth.ProvideURLAuthSuffix()).
		Set("Authorization", fmt.Sprintf("Bearer %s", auth.Token))

	return newClientInstance(address, httpClient, slingBase)
}

func newClientInstance(address string, httpClient sling.Doer, sling *sling.Sling) *Client {
	slingRest := sling.New().Path(PathSuffixREST)
	slingLog := sling.New().Path(PathSuffixLog)

	return &Client{
		address:    address,
		commonBase: sling,
		Projects:   subapi.NewProjectService(slingRest.New(), httpClient),
		BuildTypes: subapi.NewBuildTypeService(slingRest.New(), httpClient),
		Builds:     subapi.NewBuildService(slingRest.New(), httpClient),
		BuildQueue: subapi.NewBuildQueueService(slingRest.New(), httpClient),
		Token:      subapi.NewTokenService(slingRest.New(), httpClient),
		Logs: subapi.NewLogService(slingLog.New(), httpClient),
	}
}

// Ping tests if the client is properly configured and can be used
func (c *Client) Ping() error {
	response, err := c.commonBase.Get("app/rest/server").ReceiveSuccess(nil)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 && response.StatusCode != 403 {
		fmt.Println(response.StatusCode)
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("API error %s: %s", response.Status, body)
	}

	return nil
}
