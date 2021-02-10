package describe

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

type buildGetterMock struct {
	build      subapi.BuildJSON
	properties subapi.Properties

	errGetBuild        error
	errGetBuildResults error
}

func newBuildGetterMock(
	build subapi.BuildJSON, properties subapi.Properties, errGetBuild, errGetResults error) *buildGetterMock {
	return &buildGetterMock{
		build:              build,
		properties:         properties,
		errGetBuild:        errGetBuild,
		errGetBuildResults: errGetResults,
	}
}

func (m *buildGetterMock) GetBuild(_ string, _ string) (build subapi.BuildJSON, err error) {
	return m.build, m.errGetBuild
}

func (m *buildGetterMock) GetBuildResults(_ string) (resultingProperties subapi.Properties, err error) {
	return m.properties, m.errGetBuildResults
}

type buildDescriptionWriterMock struct {
	build subapi.BuildJSON
}

func newBuildDescriptionWriterMock() *buildDescriptionWriterMock {
	return &buildDescriptionWriterMock{}
}

func (m *buildDescriptionWriterMock) WriteBuildDescription(build subapi.BuildJSON) {
	m.build = build
}

func (m *buildDescriptionWriterMock) getBuild() subapi.BuildJSON {
	return m.build
}

var (
	buildTest               = subapi.BuildJSON{ID: 1}
	propertiesTest          = subapi.Properties{Items: []*subapi.Property{{Name: "test", Value: "test"}}}
	buildWithPropertiesTest = subapi.BuildJSON{
		ID:                  1,
		ResultingProperties: subapi.Properties{Items: []*subapi.Property{{Name: "test", Value: "test"}}},
	}

	errBuildTest      = errors.New("error build")
	errPropertiesTest = errors.New("error properties")
)

type testCaseDescribeBuild struct {
	name          string
	build         subapi.BuildJSON
	properties    subapi.Properties
	buildExpected subapi.BuildJSON
	shortFlag     bool

	errBuild      error
	errProperties error
	errExpected   error
}

func TestDescribeBuild(t *testing.T) {
	cases := []testCaseDescribeBuild{
		{
			name:          "Error while getting build",
			errBuild:      errBuildTest,
			errProperties: errPropertiesTest,
			errExpected:   errBuildTest,
		},
		{
			name:          "Error while getting build properties, short flag is disabled",
			errProperties: errPropertiesTest,
			errExpected:   errPropertiesTest,
		},
		{
			name:          "Positive get build, short flag is present, properties are not requested",
			build:         buildTest,
			buildExpected: buildTest,
			errProperties: errPropertiesTest,
			shortFlag:     true,
		},
		{
			name:          "Positive get build, positive get properties",
			build:         buildTest,
			properties:    propertiesTest,
			buildExpected: buildWithPropertiesTest,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			writerMock := newBuildDescriptionWriterMock()

			err := describeBuild(
				newBuildGetterMock(testCase.build, testCase.properties, testCase.errBuild, testCase.errProperties),
				writerMock, "", "", testCase.shortFlag,
			)

			if testCase.errExpected != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.errExpected, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.buildExpected, writerMock.getBuild())
		})
	}
}
