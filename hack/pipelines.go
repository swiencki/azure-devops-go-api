package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
)

func main() {
	// ENV variable AZURE_DEVOPS_EXT_PAT="" default None
	personalAccessToken := os.Getenv("AZURE_DEVOPS_EXT_PAT")
	if len(personalAccessToken) == 0 {
		log.Fatal("Personal Access Token required. | AZURE_DEVOPS_EXT_PAT=''")
	} else {
		ctx := context.Background()

		// ENV variable ADO_ORG='' default https://dev.azure.com/msazure/
		organizationUrl := getEnv("ADO_ORG", "https://dev.azure.com/msazure/")

		// ENV variable ADO_PROJECT='' default AzureRedHatOpenShift
		project := getEnv("ADO_PROJECT", "AzureRedHatOpenShift")

		// ENV variable ADO_TIME='' default None
		userTime := os.Getenv("ADO_TIME")
		intUserTime, err := strconv.Atoi(userTime)
		if err != nil {
			log.Print("No ADO_TIME set.")
		}

		// ENV variable ADO_DEFID='' default None
		userDefID := os.Getenv("ADO_DEFID")
		intDefID, err := strconv.Atoi(userDefID)
		if err != nil {
			log.Print("No Definitions ID set.")
		}
		arrayIntDefID := []int{intDefID}

		// ENV variable ADO_BUILDID='' default None
		userBuildID := os.Getenv("ADO_BUILDID")
		intBuildID, err := strconv.Atoi(userBuildID)
		if err != nil {
			log.Print("No Build ID set.")
		}
		arrayIntBuildID := []int{intBuildID}

		// Create a connection to your organization
		connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)

		// Create a client to interact with build
		buildClient, err := build.NewClient(ctx, connection)
		if err != nil {
			log.Fatal(err)
		}

		// Sets Args for buildClient.GetBuilds
		// BuildID can have no other arguments in statement
		// TODO: Allow multiple ADO_DEFID's and ADO_BUILDID's in comma-delimited list format
		getBuildsArgs := build.GetBuildsArgs{Project: &project}
		if len(arrayIntDefID) == 1 && arrayIntDefID[0] != 0 {
			getBuildsArgs.MinTime = &azuredevops.Time{Time: time.Now().Add(time.Duration(-intUserTime) * time.Minute)}
			getBuildsArgs.Definitions = &arrayIntDefID
		} else {
			getBuildsArgs.BuildIds = &arrayIntBuildID
		}

		// Returns a list of builds based off of getBuildsArgs
		buildResponse, err := buildClient.GetBuilds(ctx, getBuildsArgs)
		if err != nil {
			log.Fatal(err)
		}

		// Returns LogLines based off user inputs
		getBuildLogs(ctx, buildResponse.Value, buildClient, project)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getBuildLogs(ctx context.Context, buildResponse []build.Build, buildClient build.Client, project string) (*[]build.BuildLog, error) {
	for len(buildResponse) != 0 {
		for _, BuildReference := range buildResponse {
			// Get first page of BuildLogs from organization / project / buildResponse / BuildReference.Id within timeframe`
			buildLogsResponse, err := buildClient.GetBuildLogs(ctx, build.GetBuildLogsArgs{Project: &project, BuildId: BuildReference.Id})
			if err != nil {
				return nil, err
			}
			for buildLogsResponse != nil {
				for _, BuildLogsReference := range *buildLogsResponse {
					// Get first page of LogLines from organization / project / buildLogsResponse / BuildReference.Id / BuildLogsReference.Id within timefram
					buildLogsLinesResponse, err := buildClient.GetBuildLogLines(ctx, build.GetBuildLogLinesArgs{Project: &project, BuildId: BuildReference.Id, LogId: BuildLogsReference.Id})
					if err != nil {
						return nil, err
					}
					for buildLogsLinesResponse != nil {
						for _, buildLogsLineReference := range *buildLogsLinesResponse {
							log.Print(buildLogsLineReference)
						}
					}
				}
			}
		}
	}
	log.Print("No results returned.")
	// Returns nil if no values returned in buildResponse
	return nil, nil
}

// Sets env var or sets defaults if nothing is specified
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
