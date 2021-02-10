package subapi

import (
	"fmt"

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

type TriggerResultJSON struct {
	ID    int    `json:"id,omitempty" xml:"id"`
	State string `json:"state,omitempty" xml:"state"`

	Triggered  TriggeredJSON `json:"triggered,omitempty" xml:"triggered"`
	Properties Parameters    `json:"properties,omitempty" xml:"properties"`
}

func (t TriggerResultJSON) TriggerResult() TriggerResult {
	return TriggerResult{
		BuildID:     fmt.Sprint(t.ID),
		BuildState:  t.State,
		TriggeredBy: t.Triggered.User.Name,
		Parameters:  t.Properties,
	}
}

type TriggeredJSON struct {
	User User   `json:"user,omitempty" xml:"user"`
	Date string `json:"date,omitempty" xml:"date"`
}

type User struct {
	Username string `json:"username,omitempty" xml:"username"`
	Name     string `json:"name,omitempty" xml:"name"`
}

type BuildQueueService struct {
	sling         *sling.Sling
	httpClient    sling.Doer
	requestsMaker *requestsMaker
}

func NewBuildQueueService(base *sling.Sling, client sling.Doer) *BuildQueueService {
	s := base.Path("buildQueue/")
	return &BuildQueueService{
		sling:         s,
		httpClient:    client,
		requestsMaker: newRequestsMakerWithSling(client, s),
	}
}

func (s *BuildQueueService) RunBuildByBuildConfID(buildconfID string) (result TriggerResult, err error) {
	data := triggerData{BuildType: buildTypeJSON{ID: buildconfID}}
	jsonStruct := &TriggerResultJSON{}

	err = s.requestsMaker.post("", data, &jsonStruct)
	return jsonStruct.TriggerResult(), err
}
