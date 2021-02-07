package subapi

import (
	"encoding/json"

	"github.com/dghubble/sling"
)

type buildTypeJSON struct {
	Description       string      `json:"description,omitempty" xml:"description"`
	Href              string      `json:"href,omitempty" xml:"href"`
	ID                string      `json:"id,omitempty" xml:"id"`
	InternalID        string      `json:"internalId,omitempty" xml:"internalId"`
	Locator           string      `json:"locator,omitempty" xml:"locator"`
	Name              string      `json:"name,omitempty" xml:"name"`
	Parameters        *Parameters `json:"parameters,omitempty"`
	Paused            *bool       `json:"paused,omitempty" xml:"paused"`
	Project           *Project    `json:"project,omitempty"`
	ProjectID         string      `json:"projectId,omitempty" xml:"projectId"`
	ProjectInternalID string      `json:"projectInternalId,omitempty" xml:"projectInternalId"`
	ProjectName       string      `json:"projectName,omitempty" xml:"projectName"`
	TemplateFlag      *bool       `json:"templateFlag,omitempty" xml:"templateFlag"`
	Type              string      `json:"type,omitempty" xml:"type"`
	UUID              string      `json:"uuid,omitempty" xml:"uuid"`
	Settings          *Properties `json:"settings,omitempty"`
	Templates         *Templates  `json:"templates,omitempty"`
}

// Templates represents a collection of BuildTypeReference.
type Templates struct {
	Count int32                 `json:"count,omitempty" xml:"count"`
	Items []*BuildTypeReference `json:"buildType"`
}

// BuildType represents a build configuration or a build configuration template
type BuildType struct {
	ProjectID   string
	ID          string
	Name        string
	Description string
	Disabled    bool
	IsTemplate  bool
	// steps are a bit complex and useless for now
	Templates *Templates

	Parameters    *Parameters
	buildTypeJSON *buildTypeJSON
}

// UnmarshalJSON implements JSON deserialization for TriggerSchedule
func (b *BuildType) UnmarshalJSON(data []byte) error {
	var aux buildTypeJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if err := b.read(&aux); err != nil {
		return err
	}

	return nil
}

func (b *BuildType) read(dt *buildTypeJSON) error {
	var isTemplate bool
	if dt.TemplateFlag != nil {
		isTemplate = *dt.TemplateFlag
		b.IsTemplate = isTemplate
	}

	b.ID = dt.ID
	b.Name = dt.Name
	b.Description = dt.Description
	b.ProjectID = dt.ProjectID
	b.Parameters = dt.Parameters
	b.Templates = dt.Templates
	return nil
}

// BuildTypeReference represents a brief of a BuildType
type BuildTypeReference struct {
	ID        string `json:"id,omitempty" xml:"id"`
	Name      string `json:"name,omitempty" xml:"name"`
	ProjectID string `json:"projectId,omitempty" xml:"projectId"`
}

// BuildTypeReferences represents a collection of *BuildTypeReference
type BuildTypeReferences struct {
	Count int32                 `json:"count,omitempty" xml:"count"`
	Items []*BuildTypeReference `json:"buildType"`
}

// BuildTypeService has operations for handling build configurations and templates
type BuildTypeService struct {
	sling         *sling.Sling
	httpClient    sling.Doer
	requestsMaker *requestsMaker
}

func NewBuildTypeService(base *sling.Sling, httpClient sling.Doer) *BuildTypeService {
	s := base.Path("buildTypes/")
	return &BuildTypeService{
		httpClient:    httpClient,
		sling:         s,
		requestsMaker: newRequestsMakerWithSling(httpClient, s),
	}
}
