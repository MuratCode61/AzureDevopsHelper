package azuredevopsclient

import (
	"azure_devops_helper/config"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
)

var httpClient *http.Client = &http.Client{}

var azureDevopsConnection *azuredevops.Connection

func getHttpClient() *http.Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return httpClient
}

func GetAzureDevopsConnection() *azuredevops.Connection {
	if azureDevopsConnection == nil {
		appConfig := config.GetAppConfig()
		azureDevopsConnection = azuredevops.NewPatConnection(appConfig.OrganizationUrl, appConfig.AuthConfig.PersonalAccessToken)
	}
	return azureDevopsConnection
}

func DownloadFile(filePath string, fileUrl string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error has been occurred while creating file.", err)
	}
	defer file.Close()

	httpRequest, err := http.NewRequest(http.MethodGet, fileUrl, nil)
	if err != nil {
		log.Println("Error has been occurred while creating http request.", err)
	}

	httpRequest.SetBasicAuth("", config.GetAppConfig().AuthConfig.PersonalAccessToken)
	response, err := getHttpClient().Do(httpRequest)
	if err != nil {
		log.Println("Error has been occurred while downloading file.", err)
	}

	if response.StatusCode == 200 {
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Println("Error has been occurred while reading response", err)
		}
	} else {
		log.Println("Downloading File is failed: Response Status Code:", response.StatusCode, "fileUrl:", fileUrl)
	}

	defer response.Body.Close()
}
