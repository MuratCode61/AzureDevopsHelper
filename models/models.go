package models

import "time"

type WorkItemSummary struct {
	Id                        int
	WorkItemType              string
	Title                     string
	State                     string
	CreatedDate               string
	Severity                  string
	CreatedBy                 string
	AssignedTo                string
	IterationPath             string
	Description               HtmlContent
	ReproSteps                HtmlContent
	Priority                  float64
	WorkItemCommentSummaries  []WorkItemCommentSummary
	WorkItemRelationSummaries []WorkItemRelationSummary
}

type WorkItemCommentSummary struct {
	Content     HtmlContent
	CreatedDate time.Time
	CreatedBy   string
}

type WorkItemRelationSummary struct {
	Name string
	Url  string
}

type HtmlContent struct {
	Text   string
	Images []Image
}

type Image struct {
	Name string
	Url  string
}

type ExportWorkItemsRequest struct {
	QueryId string
}

type AuthConfig struct {
	PersonalAccessToken string `mapstructure:"personalAccessToken"`
}

type AppConfig struct {
	AuthConfig      AuthConfig `mapstructure:"auth"`
	OrganizationUrl string     `mapstructure:"organizationurl"`
	FileExportDir   string     `mapstructure:"fileExportDir"`
}
