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
	organizationUrl := ""
	personalAccessToken := ""

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

	//	TODO: This should be to specify which pipeline
	userDefID := os.Getenv("ADO_DEFID")
	if len(userDefID) == 0 {
		log.Fatal("A comma-delimited list of Definitions IDs required. | ADO_DEFID=''")
	} else {
		intDefID, err := strconv.Atoi(userDefID)
		if err != nil {
			log.Fatal(err)
		}
		arIntDefID := []int{intDefID}

		// Create a client to interact with build
		buildClient, err := build.NewClient(ctx, connection)
		if err != nil {
			log.Fatal(err)
		}

		// Get first page of builds for project / organization
		buildResponseValue, err := buildClient.GetBuilds(ctx, build.GetBuildsArgs{Project: &project, MinTime: &azuredevops.Time{time.Now().Add(time.Duration(-intUserTime) * time.Minute)}, Definitions: &arIntDefID})
		if err != nil {
			log.Fatal(err)
		}
		//	log.Print(buildResponseValue)
		for buildResponseValue != nil {
			for _, BuildReference := range (*buildResponseValue).Value {
				log.Print(&BuildReference.Id) //DEBUG: Prints buildId
				buildLogsResponseValue, err := buildClient.GetBuildLogs(ctx, build.GetBuildLogsArgs{Project: &project, BuildId: BuildReference.Id})
				if err != nil {
					log.Fatal(err)
				}
				for buildLogsResponseValue != nil {
					for _, BuildLogsReference := range *buildLogsResponseValue {
						log.Print(&BuildLogsReference.Id)
					}
				}
			}
		}
	}

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
