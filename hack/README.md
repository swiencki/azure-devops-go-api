# Azure DevOps Golang API scripts
This repository contains Golang API scripts for interacting with and managing Azure DevOps.

To learn more about the Azure DevOps Extension for Azure CLI, visit the [Microsoft/azure-devops-cli-extension](https://github.com/Microsoft/azure-devops-cli-extension) repo and [azure-devops-go-api](https://github.com/Microsoft/azure-devops-go-api) repo.
## Get started
To use the Golang API scripts establish a connection using a [personal access token](https://docs.microsoft.com/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=vsts).

## pipelines.go 

### Usage
Run pipelines.go
```bash
$ [env_variable1 env_variable2] go run hack/pipelines.go
```
Example
```bash
$  AZURE_DEVOPS_EXT_PAT=foo ADO_DEFID=123 go run hack/pipelines.go
```
### Env
List of enviromental variables
```Golang
ADO_ORG='' // Organization Url *string | default https://dev.azure.com/msazure/
ADO_PROJECT='' // Project Name *string | default AzureRedHatOpenShift
AZURE_DEVOPS_EXT_PAT='' // Personal Access Token *string | default None | Required
ADO_DEFID='' // Comma-delimited list of Definition IDs *int | default None | Required
```

## boards.go 

### Usage
Run boards.go
```bash
$ [env_variable1 env_variable2] go run hack/boards.go
```
Example
```bash
$  AZURE_DEVOPS_EXT_PAT=foo go run hack/boards.go
```
### Env
List of enviromental variables:
```Golang
ADO_ORG='' // Organization Url *string | default https://dev.azure.com/msazure/
```