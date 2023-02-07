package main

import (
	"azure_devops_helper/models"
	"azure_devops_helper/workitems"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
)

func main() {

	http.HandleFunc("/exportWorkItems", func(w http.ResponseWriter, r *http.Request) {
		var exportWorkItemsRequest models.ExportWorkItemsRequest
		err := json.NewDecoder(r.Body).Decode(&exportWorkItemsRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		currentTime := time.Now()
		workItemsFolderName := fmt.Sprintf("WorkItems%02d%02d%d", currentTime.Day(), currentTime.Month(), currentTime.Year())

		workItemsZipPath := workitems.GetAndExportWorkItems(exportWorkItemsRequest, workItemsFolderName)

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fmt.Sprintf("%s.zip", workItemsFolderName)))
		http.ServeFile(w, r, workItemsZipPath)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
