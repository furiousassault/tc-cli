package subapi

import (
	"github.com/dghubble/sling"
)

type LogService struct {
	sling         *sling.Sling
	httpClient    sling.Doer
	requestsMaker *requestsMaker
}

func NewLogService(base *sling.Sling, client sling.Doer) *LogService {
	s := base.Path("downloadBuildLog.html")
	return &LogService{
		sling:         s,
		httpClient:    client,
		requestsMaker: newRequestsMakerWithSling(client, s),
	}
}

type logQueryParameters struct {
	BuildID string `url:"buildId"`
}

func (l *LogService) GetBuildLog(buildID string) (out []byte, err error) {
	return l.requestsMaker.getResponseBytes("", &logQueryParameters{BuildID: buildID},
	)
}
