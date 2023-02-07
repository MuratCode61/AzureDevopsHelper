package excel

import (
	"azure_devops_helper/models"
	"fmt"
	"log"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

func WriteWorkItemsToExcel(workItemSummaries []models.WorkItemSummary, exportPath string) {
	excel := excelize.NewFile()
	createWorkItemsSheet(excel, workItemSummaries)
	createWorkItemsCommentsSheet(excel, workItemSummaries)
	if err := excel.SaveAs(fmt.Sprintf("%s%cWorkItems.xlsx", exportPath, filepath.Separator)); err != nil {
		log.Println(err)
	}
}

func createWorkItemsSheet(excel *excelize.File, workItemSummaries []models.WorkItemSummary) {
	excel.SetColWidth("Sheet1", "A", "A", 10)
	excel.SetColWidth("Sheet1", "B", "B", 20)
	excel.SetColWidth("Sheet1", "C", "C", 50)
	excel.SetColWidth("Sheet1", "D", "D", 10)
	excel.SetColWidth("Sheet1", "E", "E", 15)
	excel.SetColWidth("Sheet1", "F", "F", 15)
	excel.SetColWidth("Sheet1", "G", "G", 100)
	excel.SetColWidth("Sheet1", "H", "H", 100)

	headers := map[string]string{
		"A1": "ID",
		"B1": "Work Item Type",
		"C1": "Title",
		"D1": "Severity",
		"E1": "Created By",
		"F1": "Created Date",
		"G1": "Description",
		"H1": "Retro Steps",
	}

	headerStyle, err := excel.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 12, Bold: true},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#AAC57D"}}})
	if err != nil {
		log.Fatalln("Error while creating header style.", err)
	}
	excel.SetCellStyle("Sheet1", "A1", "H1", headerStyle)

	for k, v := range headers {
		excel.SetCellValue("Sheet1", k, v)
	}

	axis := 1
	for _, workItemSummary := range workItemSummaries {
		axis++
		excel.SetSheetRow("sheet1", fmt.Sprintf("A%d", axis), &[]interface{}{
			workItemSummary.Id,
			workItemSummary.WorkItemType,
			workItemSummary.Title,
			workItemSummary.Severity,
			workItemSummary.CreatedBy,
			workItemSummary.CreatedDate,
			workItemSummary.Description,
			workItemSummary.ReproSteps,
		})
	}
}

func createWorkItemsCommentsSheet(excel *excelize.File, workItemSummaries []models.WorkItemSummary) {
	sheetId := excel.NewSheet("Sheet2")
	excel.SetActiveSheet(sheetId)

	excel.SetColWidth("Sheet2", "A", "A", 20)
	excel.SetColWidth("Sheet2", "B", "B", 20)
	excel.SetColWidth("Sheet2", "C", "C", 100)

	commentHeaderStyle, err := excel.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 12, Bold: true},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#87cefa"}}})
	if err != nil {
		log.Fatalln("Error while creating comment header style.", err)
	}

	workItemSummaryHeaderStyle, err := excel.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 12, Bold: true},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#00ffff"}}})
	if err != nil {
		log.Fatalln("Error while creating workitem header style.", err)
	}

	axis := 0
	for _, workItemSummary := range workItemSummaries {
		axis++

		workItemSummaryHeader := fmt.Sprintf("%d - %s", workItemSummary.Id, workItemSummary.Title)
		excel.MergeCell("Sheet2", fmt.Sprintf("A%d", axis), fmt.Sprintf("C%d", axis))
		excel.SetCellStyle("Sheet2", fmt.Sprintf("A%d", axis), fmt.Sprintf("C%d", axis), workItemSummaryHeaderStyle)
		excel.SetSheetRow("Sheet2", fmt.Sprintf("A%d", axis), &[]interface{}{workItemSummaryHeader})

		axis++

		excel.SetCellStyle("Sheet2", fmt.Sprintf("A%d", axis), fmt.Sprintf("C%d", axis), commentHeaderStyle)
		excel.SetSheetRow("Sheet2", fmt.Sprintf("A%d", axis), &[]interface{}{"Created By", "Created Date", "Comment"})

		for _, workItemCommentSummary := range workItemSummary.WorkItemCommentSummaries {
			axis++

			excel.SetSheetRow("Sheet2", fmt.Sprintf("A%d", axis), &[]interface{}{
				workItemCommentSummary.CreatedBy,
				workItemCommentSummary.CreatedDate,
				workItemCommentSummary.Content.Text,
			})
		}

		axis++ // empty row between workitem comments
	}
}
