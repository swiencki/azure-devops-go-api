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
	personalAccessToken := os.Getenv("AZURE_DEVOPS_EXT_PAT") // ENV variable AZURE_DEVOPS_EXT_PAT="" default None
	if len(personalAccessToken) == 0 {
		log.Fatal("Personal Access Token required. | AZURE_DEVOPS_EXT_PAT=''")
	} else {
		// ENV variable ADO_ORG='' default https://dev.azure.com/msazure/
		organizationUrl := getEnv("ADO_ORG", "https://dev.azure.com/msazure/")

		ctx := context.Background()

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
		arryIntDefID := []int{intDefID}

		// ENV variable ADO_BUILDID='' default None
		userBuildID := os.Getenv("ADO_BUILDID")
		intBuildID, err := strconv.Atoi(userBuildID)
		if err != nil {
			log.Print("No Build ID set.")
		}
		arryIntBuildID := []int{intBuildID}

		// Create a connection to your organization
		connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
		// Create a client to interact with build
		buildClient, err := build.NewClient(ctx, connection)
		if err != nil {
			log.Fatal(err)
		}

		getBuildArgs := build.GetBuildsArgs{Project: &project}
		if len(arryIntBuildID) == 1 {
			getBuildArgs.BuildIds = &arryIntBuildID
		} else {
			getBuildArgs.MinTime = &azuredevops.Time{Time: time.Now().Add(time.Duration(-intUserTime) * time.Minute)}
			getBuildArgs.Definitions = &arryIntDefID
		}

		buildResponse, err := buildClient.GetBuilds(ctx, getBuildArgs)
		if err != nil {
			log.Fatal(err)
		}

		for len(buildResponse.Value) != 0 {
			for _, BuildReference := range (*buildResponse).Value {
				//Get first page of buildlogs from organization / project / definitions / buildid within timeframe
				buildLogsResponse, err := buildClient.GetBuildLogs(ctx, build.GetBuildLogsArgs{Project: &project, BuildId: BuildReference.Id})
				if err != nil {
					log.Fatal(err)
				}
				for buildLogsResponse != nil {
					for _, BuildLogsReference := range *buildLogsResponse {
						//Get first page of loglines from organization / project / definitions / buildid / logid within timeframe
						buildLogsLinesResponse, err := buildClient.GetBuildLogLines(ctx, build.GetBuildLogLinesArgs{Project: &project, BuildId: BuildReference.Id, LogId: BuildLogsReference.Id})
						if err != nil {
							log.Fatal(err)
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
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
