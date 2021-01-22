package subapi

import (
	"net/http"

	"github.com/dghubble/sling"
)

type LogService struct {
	sling      *sling.Sling
	httpClient *http.Client
	requestsMaker *requestsMaker
}

func NewLogService(base *sling.Sling, client *http.Client) *LogService {
	sling := base.Path("downloadBuildLog.html")
	return &LogService{
		sling:      sling,
		httpClient: client,
		requestsMaker: newRequestsMakerWithSling(client, sling),
	}
}

type LogQueryParameters struct {
	BuildId string `url:"buildId"`
}

func (l *LogService) GetBuildLog(buildID string) (out []byte, err error) {
	return l.requestsMaker.getResponseBytes("",
		&LogQueryParameters{BuildId: buildID},
		"log",
	)
}
