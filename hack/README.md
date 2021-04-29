
# Azure DevOps Golang API

This repository contains Golang APIs for interacting with and managing Azure DevOps. These APIs power the Azure DevOps Extension for Azure CLI. To learn more about the Azure DevOps Extension for Azure CLI, visit the [Microsoft/azure-devops-cli-extension](https://github.com/Microsoft/azure-devops-cli-extension) repo and [azure-devops-go-api](https://github.com/Microsoft/azure-devops-go-api) repo.

## Install 

```
git clone <repo>
go mod init
go mod tidy
```

## Get started


To use the API, establish a connection using a [personal access token](https://docs.microsoft.com/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=vsts) and the URL to your Azure DevOps organization. Then get a client from the connection and make API calls.

```Golang
# Fill in with your personal access token, org URL and project
personal_access_token = 'YOURPAT'
organization_url = 'https://dev.azure.com/YOURORG'
project = 'YOURPROJECT'
```

## Usage
Run build.go
```bash
$ ADO_DEFID=<Pipeline ID> go run hack/build.go
```
