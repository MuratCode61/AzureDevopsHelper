package workitems

import (
	"azure_devops_helper/azuredevopsclient"
	"azure_devops_helper/config"
	"azure_devops_helper/excel"
	"azure_devops_helper/htmlparser"
	"azure_devops_helper/models"
	"azure_devops_helper/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
)

func GetAndExportWorkItems(exportWorkItemsRequest models.ExportWorkItemsRequest, workItemsFolderName string) string {
	workItemSummaries := GetWorkItems(exportWorkItemsRequest.QueryId)
	exportDir := filepath.FromSlash(config.GetAppConfig().FileExportDir)
	workItemsExportPath := fmt.Sprintf("%s%c%s", exportDir, os.PathSeparator, workItemsFolderName)
	os.Mkdir(workItemsExportPath, os.ModePerm)
	excel.WriteWorkItemsToExcel(workItemSummaries, workItemsExportPath)
	downloadWorkItemFiles(workItemSummaries, workItemsExportPath)
	return utils.ZipFolder(exportDir, workItemsExportPath)
}

func GetWorkItems(queryId string) []models.WorkItemSummary {
	azureDevopsConnection := azuredevopsclient.GetAzureDevopsConnection()

	ctx := context.Background()
	workitemtrackingConnection, err := workitemtracking.NewClient(ctx, azureDevopsConnection)
	if err != nil {
		log.Fatal("Error while getting work item tracking connection.", err)
	}

	workItemIds := getWorkItemsIdsByQueryId(workitemtrackingConnection, queryId)

	workItems, err := workitemtrackingConnection.GetWorkItems(ctx, workitemtracking.GetWorkItemsArgs{Ids: &workItemIds, Expand: &workitemtracking.WorkItemExpandValues.All})
	if err != nil {
		log.Fatal("Error while getting work item details: ", err)
	}

	var workItemSummaries []models.WorkItemSummary
	var projectId *string

	// TODO: concurrent processing for workitems
	if workItems != nil {
		for _, workItem := range *workItems {
			if projectId == nil {
				projectId = getProjectIdFromWorkItemUrl(*workItem.Url)
			}

			workItemSummary := convertWorkItemToWorkItemSummary(workitemtrackingConnection, workItem, projectId)
			workItemSummaries = append(workItemSummaries, workItemSummary)
		}
	}

	return workItemSummaries
}

func getWorkItemsIdsByQueryId(workitemtrackingConnection workitemtracking.Client, queryId string) []int {
	uuidVal := uuid.MustParse(queryId)
	workItemQueryResult, err := workitemtrackingConnection.QueryById(context.Background(), workitemtracking.QueryByIdArgs{Id: &uuidVal})
	if err != nil {
		log.Fatal("Error while getting work items by query: ", err)
	}

	var workItemIds []int

	if workItemQueryResult != nil {
		for _, workItemReference := range *workItemQueryResult.WorkItems {
			workItemIds = append(workItemIds, *workItemReference.Id)
		}
	}

	return workItemIds
}

func convertWorkItemToWorkItemSummary(workitemtrackingConnection workitemtracking.Client, workItem workitemtracking.WorkItem, projectId *string) models.WorkItemSummary {
	createdBy := workitemtracking.IdentityReference{}
	convertMapToStruct((*workItem.Fields)["System.CreatedBy"], &createdBy)

	assignedTo := workitemtracking.IdentityReference{}
	convertMapToStruct((*workItem.Fields)["System.AssignedTo"], &assignedTo)

	description := (*workItem.Fields)["System.Description"]
	var descriptionHtmlContent models.HtmlContent
	if description != nil {
		imageFileNamePrefix := fmt.Sprintf("%d_description", *workItem.Id)
		descriptionHtmlContent = htmlparser.ParseHtmlText(description.(string), imageFileNamePrefix)
	}

	retroSteps := (*workItem.Fields)["Microsoft.VSTS.TCM.ReproSteps"]
	var retroStepsHtmlContent models.HtmlContent
	if retroSteps != nil {
		imageFileNamePrefix := fmt.Sprintf("%d_retroSteps", *workItem.Id)
		retroStepsHtmlContent = htmlparser.ParseHtmlText(retroSteps.(string), imageFileNamePrefix)
	}

	workItemCommentSummaries := GetWorkItemComments(workitemtrackingConnection, projectId, workItem.Id)
	workItemRelationSummaries := GetWorkItemRelations(workItem)

	workItemSummary := models.WorkItemSummary{
		Id:                        *workItem.Id,
		WorkItemType:              ((*workItem.Fields)["System.WorkItemType"]).(string),
		Title:                     ((*workItem.Fields)["System.Title"]).(string),
		State:                     ((*workItem.Fields)["System.State"]).(string),
		CreatedDate:               ((*workItem.Fields)["System.CreatedDate"]).(string),
		CreatedBy:                 *createdBy.DisplayName,
		AssignedTo:                *assignedTo.DisplayName,
		IterationPath:             ((*workItem.Fields)["System.IterationPath"]).(string),
		Description:               descriptionHtmlContent,
		ReproSteps:                retroStepsHtmlContent,
		WorkItemCommentSummaries:  workItemCommentSummaries,
		WorkItemRelationSummaries: workItemRelationSummaries,
	}

	severity := (*workItem.Fields)["Microsoft.VSTS.Common.Severity"]
	if severity != nil {
		workItemSummary.Severity = (severity).(string)
	}

	priority := (*workItem.Fields)["Microsoft.VSTS.Common.Priority"]
	if priority != nil {
		workItemSummary.Priority = (priority).(float64)
	}

	return workItemSummary
}

func GetWorkItemRelations(workItem workitemtracking.WorkItem) []models.WorkItemRelationSummary {
	var workItemRelationSummaries []models.WorkItemRelationSummary
	if workItem.Relations != nil {
		for _, workItemRelation := range *workItem.Relations {
			if strings.EqualFold(*workItemRelation.Rel, "AttachedFile") {
				relationName := (*workItemRelation.Attributes)["name"]
				workItemRelationSummaries = append(workItemRelationSummaries, models.WorkItemRelationSummary{
					Name: relationName.(string),
					Url:  *workItemRelation.Url,
				})
			}
		}
	}
	return workItemRelationSummaries
}

func GetWorkItemComments(workitemtrackingConnection workitemtracking.Client, projectId *string, workItemId *int) []models.WorkItemCommentSummary {
	commentList, err := workitemtrackingConnection.GetComments(context.Background(), workitemtracking.GetCommentsArgs{Project: projectId, WorkItemId: workItemId})
	if err != nil {
		log.Println("Error has been occurred while retriving comments for workItem: ", *workItemId)
	}

	var workItemCommentSummaries []models.WorkItemCommentSummary

	if commentList != nil && *commentList.Count > 0 {
		for _, comment := range *commentList.Comments {
			imageFileNamePrefix := fmt.Sprintf("%d_comments", *workItemId)
			htmlContent := htmlparser.ParseHtmlText(*comment.Text, imageFileNamePrefix)
			if strings.TrimSpace(htmlContent.Text) != "" {
				workItemCommentSummary := models.WorkItemCommentSummary{
					CreatedBy:   *comment.CreatedBy.DisplayName,
					Content:     htmlContent,
					CreatedDate: comment.CreatedDate.Time,
				}
				workItemCommentSummaries = append(workItemCommentSummaries, workItemCommentSummary)
			}
		}
	}

	return workItemCommentSummaries
}

func getProjectIdFromWorkItemUrl(workItemUrl string) *string {
	organizationUrl := config.GetAppConfig().OrganizationUrl
	workItemUrlOrgUrlTrimmed := strings.TrimPrefix(workItemUrl, fmt.Sprintf("%s%s", organizationUrl, "/"))
	projectId := workItemUrlOrgUrlTrimmed[0:strings.Index(workItemUrlOrgUrlTrimmed, "/_apis")]
	return &projectId
}

func convertMapToStruct(mapObj, structObj interface{}) {
	mapAsBytes, err := json.Marshal(mapObj)
	if err != nil {
		log.Fatalln("Error has been ocurred while marshalling map to json")
	}
	err = json.Unmarshal(mapAsBytes, &structObj)
	if err != nil {
		log.Fatalln("Error has been ocurred while unmarshalling byte array to struct")
	}
}

/*
download images and attachments for workitems into folders named as workitemid
TODO: concurrent download for work items
*/
func downloadWorkItemFiles(workItemSummaries []models.WorkItemSummary, workItemsExportPath string) {

	for _, workItemSummary := range workItemSummaries {

		workItemExportPath := fmt.Sprintf("%s%c%d", workItemsExportPath, filepath.Separator, workItemSummary.Id)

		for _, image := range workItemSummary.Description.Images {
			downloadFile(workItemExportPath, image.Name, image.Url)
		}

		for _, image := range workItemSummary.ReproSteps.Images {
			downloadFile(workItemExportPath, image.Name, image.Url)
		}

		for _, workItemRelationSummary := range workItemSummary.WorkItemRelationSummaries {
			downloadFile(workItemExportPath, workItemRelationSummary.Name, workItemRelationSummary.Url)
		}

		for _, workItemCommentSummary := range workItemSummary.WorkItemCommentSummaries {
			for _, image := range workItemCommentSummary.Content.Images {
				downloadFile(workItemExportPath, image.Name, image.Url)
			}
		}
	}
}

func downloadFile(workItemExportPath string, fileName string, fileUrl string) {
	// check work item export path is exist. if not exist, create it.
	if _, err := os.Stat(workItemExportPath); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(workItemExportPath, os.ModePerm)
		}
	}
	azuredevopsclient.DownloadFile(fmt.Sprintf("%s%c%s", workItemExportPath, filepath.Separator, fileName), fileUrl)
}
