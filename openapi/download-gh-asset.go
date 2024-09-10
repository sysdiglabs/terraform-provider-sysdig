package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v56/github"
)

func main() {
	fmt.Println("Running download-gh-asset.go in fail-silently mode. Any error won't cause an os.Exit. They will be simply logged.")

	// GitHub API Key is required in order to perform GitHub operations
	authToken := os.Getenv("GITHUB_API_KEY")
	if authToken == "" {
		fmt.Println("missing env variable: GITHUB_API_KEY")
		return
	}

	f := flags{}

	// Parse flags
	flag.StringVar(&f.tag, "tag", "", "Github Release Tag")
	flag.StringVar(&f.assetName, "assetName", "", "api-spec asset name")
	flag.StringVar(&f.outputFile, "outputFile", "", "location where to write the openAPI YAML. e.g.: ./api/my-service.yaml")

	flag.Parse()

	// Validate flags
	if err := f.validate(); err != nil {
		fmt.Println(err)
		return
	}

	// Create a new GitHub Client using the env GITHUB_API_KEY
	client := github.NewClient(nil).WithAuthToken(authToken)

	// Get asset from GitHub by finding a release by tag.
	// Attempt to find the asset by name within the release and download it
	asset, err := downloadAsset(client, f.tag, f.assetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer asset.Close()

	// Finally, write the output of the asset yaml to a file
	if err = writeToFile(asset, f.outputFile); err != nil {
		fmt.Println(err)
		return
	}
}

type flags struct {
	tag        string
	assetName  string
	outputFile string
}

func (f flags) validate() error {
	msg := "unexpected nil flag: %q\n"

	if f.tag == "" {
		return fmt.Errorf(msg, "tag")
	}
	if f.assetName == "" {
		return fmt.Errorf(msg, "assetName")
	}
	if f.outputFile == "" {
		return fmt.Errorf(msg, "outputFile")
	}

	return nil
}

func writeToFile(content io.ReadCloser, location string) error {
	outFile, err := os.Create(location)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, content)

	return err
}

func downloadAsset(client *github.Client, tag string, assetName string) (io.ReadCloser, error) {
	ctx := context.Background()
	owner := "draios"
	repo := "api-spec"

	release, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while fetching release: %v\n", err)
	}
	if release == nil {
		return nil, fmt.Errorf("unexpected nil release")
	}

	releaseAsset, ok := findReleaseAssetByName(release, assetName)
	if !ok {
		return nil, fmt.Errorf("unable to find asset by name %q in release\n", assetName)
	}
	if releaseAsset == nil {
		return nil, fmt.Errorf("unexpected nil releaseAsset")
	}

	c := http.Client{}
	asset, _, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repo, *releaseAsset.ID, &c)
	if err != nil {
		return nil, fmt.Errorf("unable to get asset by id %q: %v\n", *releaseAsset.ID, err)
	}

	return asset, nil
}

func findReleaseAssetByName(release *github.RepositoryRelease, assetName string) (*github.ReleaseAsset, bool) {
	for _, asset := range release.Assets {
		if *asset.Name == assetName {
			return asset, true
		}
	}

	return nil, false
}
