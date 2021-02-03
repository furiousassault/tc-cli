package subapi

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type triggerData struct {
	BuildType buildTypeJSON `json:"buildType,omitempty" xml:"buildType"`
}

type TriggerResult struct {
	BuildID     string
	BuildState  string
	TriggeredBy string
	Parameters  Parameters
}

type TriggerResultJson struct {
	ID    int    `json:"id,omitempty" xml:"id"`
	State string `json:"state,omitempty" xml:"state"`

	Triggered  TriggeredJson `json:"triggered,omitempty" xml:"triggered"`
	Properties Parameters    `json:"properties,omitempty" xml:"properties"`
}

func (t TriggerResultJson) TriggerResult() TriggerResult {
	return TriggerResult{
		BuildID:     fmt.Sprint(t.ID),
		BuildState:  t.State,
		TriggeredBy: t.Triggered.User.Name,
		Parameters:  t.Properties,
	}
}

type TriggeredJson struct {
	User User   `json:"user,omitempty" xml:"user"`
	Date string `json:"date,omitempty" xml:"date"`
}

type User struct {
	Username string `json:"username,omitempty" xml:"username"`
	Name     string `json:"name,omitempty" xml:"name"`
}

type BuildQueueService struct {
	sling         *sling.Sling
	httpClient    *http.Client
	requestsMaker *requestsMaker
}

func NewBuildQueueService(base *sling.Sling, client *http.Client) *BuildQueueService {
	sling := base.Path("buildQueue/")
	return &BuildQueueService{
		sling:         sling,
		httpClient:    client,
		requestsMaker: newRequestsMakerWithSling(client, sling),
	}
}

func (s *BuildQueueService) RunBuildByBuildConfID(buildconfID string) (result TriggerResult, err error) {
	data := triggerData{BuildType: buildTypeJSON{ID: buildconfID}}
	jsonStruct := &TriggerResultJson{}

	err = s.requestsMaker.post("", data, &jsonStruct, "runBuildConfiguration")
	return jsonStruct.TriggerResult(), err
}
