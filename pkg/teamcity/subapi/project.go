package subapi

import (
	"fmt"

	"github.com/dghubble/sling"
)

// Project is the model for project entity.
type Project struct {
	Archived        *bool               `json:"archived,omitempty" xml:"archived"`
	Description     string              `json:"description,omitempty" xml:"description"`
	Href            string              `json:"href,omitempty" xml:"href"`
	ID              string              `json:"id,omitempty" xml:"id"`
	Name            string              `json:"name,omitempty" xml:"name"`
	Parameters      *Parameters         `json:"parameters,omitempty"`
	ParentProject   *ProjectReference   `json:"parentProject,omitempty"`
	ParentProjectID string              `json:"parentProjectId,omitempty" xml:"parentProjectId"`
	WebURL          string              `json:"webUrl,omitempty" xml:"webUrl"`
	BuildTypes      BuildTypeReferences `json:"buildTypes,omitempty" xml:"buildTypes"`
	ChildProjects   ProjectsReferences  `json:"projects,omitempty" xml:"projects"`
}

// ProjectsReferences represents list of projects
type ProjectsReferences struct {
	Count int                 `json:"count,omitempty" xml:"count"`
	Items []*ProjectReference `json:"project,omitempty"`
}

// ProjectReference contains reduced information on a project which is returned when list is being asked.
type ProjectReference struct {
	ID          string `json:"id,omitempty" xml:"id"`
	Name        string `json:"name,omitempty" xml:"name"`
	Description string `json:"description,omitempty" xml:"description"`
	Href        string `json:"href,omitempty" xml:"href"`
	WebURL      string `json:"webUrl,omitempty" xml:"webUrl"`
}

// ProjectService is a service for requesting projects.
type ProjectService struct {
	sling         *sling.Sling
	httpClient    sling.Doer
	requestsMaker *requestsMaker
}

func NewProjectService(base *sling.Sling, client sling.Doer) *ProjectService {
	sling := base.Path("projects/")
	return &ProjectService{
		sling:         sling,
		httpClient:    client,
		requestsMaker: newRequestsMakerWithSling(client, sling),
	}
}

// GetList retrieves a projects list.
func (s *ProjectService) GetList() (refs ProjectsReferences, err error) {
	err = s.requestsMaker.getJSON("", &refs)
	return
}

// GetBuildTypesList retrieves build types list of project specified by id
func (s *ProjectService) GetBuildTypesList(projectID string) (refs BuildTypeReferences, err error) {
	err = s.requestsMaker.getJSON(fmt.Sprintf("%s/buildTypes", projectID), &refs)
	return
}
