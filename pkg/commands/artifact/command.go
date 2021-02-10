package artifact

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	pkgErrors "github.com/pkg/errors"

	"github.com/spf13/cobra"
)

const commandDownloadArgsNum = 2

type artifactGetter interface {
	GetArtifact(buildID, path string) (artifactBinary []byte, err error)
}

// CreateCommandTreeArtifact creates artifact command subtree.
func CreateCommandTreeArtifact(artifactDownloader artifactGetter, outputPathDefault string) *cobra.Command {
	cmdArtifact := &cobra.Command{
		Use:   "artifact <subcommand>",
		Short: "artifact subcommand tree",
	}
	cmdArtifactDownload := &cobra.Command{
		Use:   "download <buildID> <path>",
		Short: "Download artifact ",
		Args:  cobra.ExactArgs(commandDownloadArgsNum),
	}
	outputPathPointer := cmdArtifactDownload.Flags().StringP(
		"outputPath",
		"o",
		outputPathDefault,
		"treated as path to directory and the buildID/artifact_path suffix is applied; "+
			"otherwise, treated like file to write artifact content into.",
	)
	forceFlagPointer := cmdArtifactDownload.Flags().BoolP(
		"force",
		"f",
		false,
		"override output file if it exists",
	)
	cmdArtifactDownload.RunE = createHandlerArtifactDownload(artifactDownloader, outputPathPointer, forceFlagPointer)
	cmdArtifact.AddCommand(cmdArtifactDownload)

	return cmdArtifact
}

func createHandlerArtifactDownload(buildGetter artifactGetter,
	outputPath *string, forceFlag *bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		buildID := args[0]
		artifactPath := args[1]

		return artifactDownload(buildGetter, buildID, artifactPath, *outputPath, *forceFlag)
	}
}

func artifactDownload(artifactGetter artifactGetter, buildID, artifactPath, outputPath string, forceFlag bool) error {
	outputPath, err := artifactOutputPath(outputPath, buildID, artifactPath)
	if err != nil {
		if !errors.Is(err, errFileExists) {
			return err
		}

		if !forceFlag {
			return pkgErrors.Wrapf(errFileExists,
				"specified output path '%s' is an existing file. Use -f/--force to override",
				outputPath,
			)
		}
	}

	artifact, err := artifactGetter.GetArtifact(buildID, artifactPath)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(outputPath, artifact, 0600); err != nil {
		return err
	}

	fmt.Printf("Artifact has been downloaded to '%s'\n", outputPath)
	return nil
}

func artifactOutputPath(pathProvided, buildID, artPath string) (path string, err error) {
	absPath, err := filepath.Abs(pathProvided)
	if err != nil {
		return "", err
	}

	// not sure if it's a good approach
	if strings.HasSuffix(pathProvided, "/") ||
		strings.HasSuffix(pathProvided, ".") ||
		strings.HasSuffix(pathProvided, "..") {
		absPath = filepath.Join(absPath, artifactPathSuffix(buildID, artPath))
	}

	if err = os.MkdirAll(filepath.Dir(absPath), 0775); err != nil {
		return "", err
	}

	exists, err := exists(absPath)
	if err != nil {
		return "", err
	}

	if exists {
		return absPath, errFileExists
	}

	return absPath, nil
}

func exists(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}

		return false, nil
	}

	return true, nil
}

func artifactPathSuffix(buildID, artPath string) string {
	return filepath.Join(buildID, artPath)
}
