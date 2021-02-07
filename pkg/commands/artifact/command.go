package artifact

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type artifactGetter interface {
	GetArtifact(buildID, path string) (artifactBinary []byte, err error)
}

func CreateCommandTreeArtifact(artifactDownloader artifactGetter, outputPathDefault string) *cobra.Command {
	cmdArtifact := &cobra.Command{
		Use:   "artifact <subcommand>",
		Short: "artifact subcommand tree",
	}
	cmdArtifactDownload := &cobra.Command{
		Use:   "download <buildID> <path>",
		Short: "Download artifact ",
		Args:  cobra.ExactArgs(2),
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
			return fmt.Errorf(
				"specified output path '%s' is an existing file. Use -f/--force to override",
				outputPath,
			)
		}
	}

	artifact, err := artifactGetter.GetArtifact(buildID, artifactPath)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(outputPath, artifact, 0666); err != nil {
		return err
	}

	fmt.Printf("Artifact has been downloaded to '%s'\n", outputPath)
	return nil
}

func artifactOutputPath(p, buildID, artPath string) (path string, err error) {
	absP, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}

	// not sure if it's a good approach
	if strings.HasSuffix(p, "/") || strings.HasSuffix(p, ".") || strings.HasSuffix(p, "..") {
		absP = filepath.Join(absP, artifactPathSuffix(buildID, artPath))
	}

	if err = os.MkdirAll(filepath.Dir(absP), 0775); err != nil {
		return "", err
	}

	exists, err := exists(absP)
	if err != nil {
		return "", err
	}

	if exists {
		return absP, errFileExists
	}

	return absP, nil
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
