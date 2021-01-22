package subapi

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/dghubble/sling"

	"github.com/furiousassault/tc-cli/pkg/teamcity"
)

// Build represents TeamCity project build.
type Build struct {
	ID          int        `json:"id"`
	Number      string     `json:"number"`
	Status      string     `json:"status"`
	StatusText  string     `json:"statusText"`
	State       string     `json:"state"`
	BuildTypeID string     `json:"buildTypeId"`
	Progress    int        `json:"progress"`
	Properties  Parameters `json:"properties"`

	QueuedDate string `json:"queuedDate"`
	StartDate  string `json:"startDate"`
	FinishDate string `json:"finishDate"`
}

type BuildReference struct {
	ID     int    `json:"id,omitempty" xml:"id"`
	Number string `json:"number"`
	Status string `json:"status"`
	State  string `json:"state"`
}

type Builds struct {
	Count int               `json:"count,omitempty" xml:"count"`
	Items []*BuildReference `json:"build" xml:"build"`
}

type BuildService struct {
	sling         *sling.Sling
	httpClient    *http.Client
	requestsMaker *requestsMaker
}

func NewBuildService(base *sling.Sling, client *http.Client) *BuildService {
	sling := base.Path("builds/")
	return &BuildService{
		sling:         sling,
		httpClient:    client,
		requestsMaker: newRequestsMakerWithSling(client, sling),
	}
}

func (b *BuildService) GetBuildsByBuildConf(buildTypeID string, count int) (builds Builds, err error) {
	// TODO probably this locator is not properly queryescaped though it works
	locator := fmt.Sprintf("?locator=buildType:(%s)", LocatorID(buildTypeID))

	if count > 0 {
		locator = fmt.Sprintf("%s,count:%d", locator, count)
	}

	err = b.requestsMaker.get(locator, &builds, "builds")
	return
}

func (b *BuildService) GetBuild(buildTypeID string, number string) (build Build, err error) {
	// TODO probably this locator is not properly queryescaped though it works
	n := url.QueryEscape(fmt.Sprintf(",number:%s", number))
	buildTypeFull := url.QueryEscape(fmt.Sprintf("buildType:(%s)", buildTypeID))
	locator := fmt.Sprintf("%s%s", buildTypeFull, n)

	err = b.requestsMaker.get(locator, &build, "builds")
	if err != nil {
		return
	}

	queued, err := teamcity.ParseTCTimeFormat(build.QueuedDate)
	if err != nil {
		return
	}
	started, err := teamcity.ParseTCTimeFormat(build.StartDate)
	if err != nil {
		return
	}
	finished, err := teamcity.ParseTCTimeFormat(build.FinishDate)
	if err != nil {
		return
	}

	build.QueuedDate = queued.String()
	build.StartDate = started.String()
	build.FinishDate = finished.String()
	return
}
