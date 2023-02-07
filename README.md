This software export work items, work item comments, attached files to work items and screenshot in comments, descriptions or retrosteps fields. After software run, folder named `WorkItems` is created. `WorkItems` folder contains sub folders that contains attachments and screenshots for each work item. In addition, an Excel file containing the title, description, comment, etc. information of the Work Items is created under the `WorkItems` folder.

![image](https://user-images.githubusercontent.com/68951149/217345379-354e917e-b52a-451a-b8c6-816be4fe16b9.png)

For example, let's assume that for Work Item with ID 100061, there are 2 screenshots in 1st Comment, 1 screenshot in 2nd Comment, and 1 screenshot in retrosteps. Suppose there is 1 attachment file in addition. Screenshots and attachments are kept by naming as follows. “1_2” refers to the 2nd Screenshot in the 1st Comment.

![image](https://user-images.githubusercontent.com/68951149/217345348-e11ccf12-75b2-4bb7-85ae-00daab1c3d94.png)

These screenshots are referenced in the Excel file. There are 2 sheets in the created Excel file, Work Items and Work Item Comments.

The Work Items section contains the following information for each work item.

- ID
- Work Item Type
- Title
- Severity
- Created By
- Created Date
- Retro Steps

![image](https://user-images.githubusercontent.com/68951149/217345260-33ada9f1-fe1c-470c-87f6-7ac3f091eb78.png)

If a screenshot is needed for explanation, the screenshot is referenced as “look at the ….png”. In the Work Items comment section, the comments for each Work Item are kept in the following format.

![image](https://user-images.githubusercontent.com/68951149/217345448-ac20f405-b40c-4de4-b528-1c4ef5009288.png)

### Necessary Configurations To Run The Software

---

Before running the software, the following properties in the `application.yaml` file must be specified.

- organizationurl
- personalAccessToken - Personal Access Token should be specified for authentication.
- fileExportDir - should be specified temporary file operations
