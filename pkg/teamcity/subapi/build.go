package subapi

import (
	"fmt"
	"net/url"

	"github.com/dghubble/sling"

	"github.com/furiousassault/tc-cli/pkg/teamcity"
)

// BuildJSON represents TeamCity project build.
type BuildJSON struct {
	ID                  int        `json:"id"`
	Number              string     `json:"number"`
	Status              string     `json:"status"`
	StatusText          string     `json:"statusText"`
	State               string     `json:"state"`
	BuildTypeID         string     `json:"buildTypeId"`
	Progress            int        `json:"progress"`
	Properties          Parameters `json:"properties"`
	ResultingProperties Properties

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
	httpClient    sling.Doer
	requestsMaker *requestsMaker
}

func NewBuildService(base *sling.Sling, client sling.Doer) *BuildService {
	s := base.Path("builds/")
	return &BuildService{
		sling:         s,
		httpClient:    client,
		requestsMaker: newRequestsMakerWithSling(client, s),
	}
}

func (b *BuildService) GetBuildsByBuildConf(buildTypeID string, count int) (builds Builds, err error) {
	// probably this locator is not properly queryescaped though it works
	locator := fmt.Sprintf("?locator=buildType:(%s)", LocatorID(buildTypeID))

	if count > 0 {
		locator = fmt.Sprintf("%s,count:%d", locator, count)
	}

	err = b.requestsMaker.getResponseJSON(locator, &builds)
	return
}

func (b *BuildService) GetBuild(buildTypeID string, number string) (build BuildJSON, err error) {
	// probably this locator is not properly queryescaped though it works
	n := url.QueryEscape(fmt.Sprintf(",number:%s", number))
	buildTypeFull := url.QueryEscape(fmt.Sprintf("buildType:(%s)", buildTypeID))
	locator := fmt.Sprintf("%s%s", buildTypeFull, n)

	err = b.requestsMaker.getResponseJSON(locator, &build)
	if err != nil {
		return
	}

	err = transformFieldsTimeFormat(&build)
	return
}

func (b *BuildService) GetBuildResults(buildID string) (resultingProperties Properties, err error) {
	locator := fmt.Sprintf("%s/resulting-properties", LocatorID(buildID))

	err = b.requestsMaker.getResponseJSON(locator, &resultingProperties)
	return
}

func (b *BuildService) GetArtifact(buildID, path string) (artifactBinary []byte, err error) {
	locator := fmt.Sprintf("%s/artifacts/content/%s", LocatorID(buildID), path)

	return b.requestsMaker.getResponseBytes(locator, nil)
}

func transformFieldsTimeFormat(build *BuildJSON) error {
	queued, err := teamcity.ParseTCTimeFormat(build.QueuedDate)
	if err != nil {
		return err
	}
	started, err := teamcity.ParseTCTimeFormat(build.StartDate)
	if err != nil {
		return err
	}
	finished, err := teamcity.ParseTCTimeFormat(build.FinishDate)
	if err != nil {
		return err
	}

	build.QueuedDate = queued.String()
	build.StartDate = started.String()
	build.FinishDate = finished.String()
	return nil
}
