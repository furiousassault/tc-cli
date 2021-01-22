package subapi

import (
	"net/http"

	"github.com/dghubble/sling"
)

type BuildQueueService struct {
	sling         *sling.Sling
	httpClient    *http.Client
	requestsMaker *requestsMaker
}

func NewBuildQueueService(base *sling.Sling, client *http.Client) *BuildQueueService {
	sling := base.Path("buildQueue")
	return &BuildQueueService{
		sling:         sling,
		httpClient:    client,
		requestsMaker: newRequestsMakerWithSling(client, sling),
	}
}

func (b *BuildQueueService) RunBuildByBuildConfID(buildconfID string) error {
	return nil
}
