package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dghubble/sling"

	"github.com/furiousassault/tc-cli/pkg/commands/token"
	"github.com/furiousassault/tc-cli/pkg/configuration"
	"github.com/furiousassault/tc-cli/pkg/teamcity/api/authorization"
	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

const (
	PathSuffixREST = "app/rest/"
	PathSuffixLog  = ""
)

func InitAPI(config configuration.Configuration, httpClient sling.Doer) *Client {
	if config.API.Authorization.Token != "" {
		return NewClientTokenAuth(
			config.API.URL,
			httpClient,
			authorization.NewAuthorizerToken(config.API.Authorization.Token),
		)
	}

	if config.API.Authorization.Username != "" && config.API.Authorization.Password != "" {
		return NewClientBasicAuth(
			config.API.URL,
			httpClient,
			authorization.NewAuthorizerHTTP(config.API.Authorization.Username, config.API.Authorization.Password),
		)
	}

	fmt.Println("No authorization is provided, trying to use guest auth")

	return NewClientGuestAuth(config.API.URL, httpClient, authorization.NewAuthorizerGuest())
}

// Client represents the base for connecting to TeamCity
type Client struct {
	address    string
	httpClient sling.Doer
	commonBase *sling.Sling

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

	slingBase.Base(fmt.Sprintf("%s/%s", strings.TrimSuffix(address, "/"), auth.ProvideURLAuthSuffix()))
	return newClientInstance(address, httpClient, slingBase)
}

func NewClientBasicAuth(address string, httpClient sling.Doer, auth authorization.AuthorizerHTTP) *Client {
	slingBase := sling.New().
		Set("Accept", "application/json").
		Set("Origin", address)

	slingBase.Base(fmt.Sprintf("%s/%s", strings.TrimSuffix(address, "/"), auth.ProvideURLAuthSuffix())).
		SetBasicAuth(auth.Username, auth.Password)

	return newClientInstance(address, httpClient, slingBase)
}

func NewClientTokenAuth(address string, httpClient sling.Doer, auth authorization.AuthorizerToken) *Client {
	slingBase := sling.New().
		Set("Accept", "application/json").
		Set("Origin", address)

	slingBase.Base(fmt.Sprintf("%s/%s", strings.TrimSuffix(address, "/"), auth.ProvideURLAuthSuffix())).
		Set("Authorization", fmt.Sprintf("Bearer %s", auth.Token))

	return newClientInstance(address, httpClient, slingBase)
}

func newClientInstance(address string, httpClient sling.Doer, sling *sling.Sling) *Client {
	slingRest := sling.New().Path(PathSuffixREST)
	slingLog := sling.New().Path(PathSuffixLog)

	return &Client{
		address:    address,
		httpClient: httpClient,
		commonBase: sling,
		Projects:   subapi.NewProjectService(slingRest.New(), httpClient),
		BuildTypes: subapi.NewBuildTypeService(slingRest.New(), httpClient),
		Builds:     subapi.NewBuildService(slingRest.New(), httpClient),
		BuildQueue: subapi.NewBuildQueueService(slingRest.New(), httpClient),
		Token:      subapi.NewTokenService(slingRest.New(), httpClient),
		Logs:       subapi.NewLogService(slingLog.New(), httpClient),
	}
}

func (c *Client) TokenServiceCurrent() token.API {
	return c.Token
}

// TokenServiceWithTokenAuth is a hack to achieve token service reinitialization during command execution.
func (c *Client) TokenServiceWithTokenAuth(token string) {
	slingBase := sling.New().
		Set("Accept", "application/json").
		Set("Origin", c.address)

	slingBase.Base(strings.TrimSuffix(c.address, "/")).
		Set("Authorization", fmt.Sprintf("Bearer %s", token))

	slingBase.Path(PathSuffixREST)
	c.Token = subapi.NewTokenService(slingBase, c.httpClient)
}

// Ping tests if the client is properly configured and can be used
func (c *Client) Ping() error {
	r, err := c.commonBase.Get("app/rest/server").Request()
	if err != nil {
		return fmt.Errorf("error constructing request for ping: %w", err)
	}

	response, err := c.httpClient.Do(r)
	if err != nil {
		return fmt.Errorf("error ping response: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return nil
	}

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return fmt.Errorf("ping error, status %d: %s", response.StatusCode, "unauthorized")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error unmarshalling ping response body: %w", err)
	}

	return fmt.Errorf("ping response status %s, body: %s", response.Status, body)
}
