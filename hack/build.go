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
	// ENV variable ADO_ORG='' default https://dev.azure.com/msazure/
	organizationUrl := getEnv("ADO_ORG", "https://dev.azure.com/msazure/")

	// ENV variable AZURE_DEVOPS_EXT_PAT="" default None
	personalAccessToken := os.Getenv("AZURE_DEVOPS_EXT_PAT")
	if len(personalAccessToken) == 0 {
		log.Fatal("Personal Access Token required. | AZURE_DEVOPS_EXT_PAT=''")
	} else {
		// Create a connection to your organization
		connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)

		ctx := context.Background()

		// ENV variable ADO_PROJECT='' default AzureRedHatOpenShift
		project := getEnv("ADO_PROJECT", "AzureRedHatOpenShift")

		// ENV variable ADO_TIME='' default 60 minutes
		userTime := getEnv("ADO_TIME", "60")
		intUserTime, err := strconv.Atoi(userTime)
		if err != nil {
			log.Fatal(err)
		}

		// ENV variable ADO_DEFID='' default None
		userDefID := os.Getenv("ADO_DEFID")
		if len(userDefID) == 0 {
			log.Fatal("A comma-delimited list of Definitions IDs required. | ADO_DEFID=''")
		} else {
			intDefID, err := strconv.Atoi(userDefID)
			if err != nil {
				log.Fatal(err)
			}
			arryIntDefID := []int{intDefID}

			// Create a client to interact with build
			buildClient, err := build.NewClient(ctx, connection)
			if err != nil {
				log.Fatal(err)
			}

			// Get first page of builds for organization / project / definitions within timeframe
			buildResponse, err := buildClient.GetBuilds(ctx, build.GetBuildsArgs{Project: &project, MinTime: &azuredevops.Time{time.Now().Add(time.Duration(-intUserTime) * time.Minute)}, Definitions: &arryIntDefID})
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
			log.Print("No results returned, expand window in minutes if required. | ADO_TIME='' default 60 minutes")
		}
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
