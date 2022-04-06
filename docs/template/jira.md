# Create Template for JIRA

Ask to JIRA rest API to get which custom fields are usable.

https://developer.atlassian.com/server/jira/platform/jira-rest-api-examples/#creating-an-issue-examples

## Get information with API

TOKEN is access token but also you can request with basic auth.

With this query, we see service increment's issuetypeid.

```sh
curl -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" "https://jira.techno.ingenico.com/rest/api/2/issue/createmeta/LBO/issuetypes" | jq .
```

After that we need to check detail fields.

```sh
curl -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" "https://jira.techno.ingenico.com/rest/api/2/issue/createmeta/LBO/issuetypes/11707" | jq .
```

Check your issues

```sh
curl -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" "https://jira.techno.ingenico.com/rest/api/2/search?jql=assignee=eates" | jq .
```

Look an example issue

```sh
curl -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" "https://jira.techno.ingenico.com/rest/api/2/issue/LBO-72121" | jq .
```

Example: https://jira.techno.ingenico.com/browse/LBO-72121

Click export to XML format button to view datas but with looking api is better.

## Template

### Service increment

Template for create service increments

If field export value list give `{"value": "blabla"}` as value.

```json
{
  "fields": {
    "project":
    {
       "key": "LBO"
    },
    "summary": "{{.summary}}",
    "description": {{or (.description | quote) "null"}},
    "issuetype": {
      "name": "Service Increment"
    },
    "priority": {
      "name": "Minor"
    },
    "customfield_10006": "{{.epic}}",
    "customfield_11601": {
      "value": "FinOps - DeepCore"
    }
  }
}
```

Values

```yml
summary: Release - Validator to create TRSVALEX event when invalidating a transaction
epic: LBO-59087
squad: FinOps - DeepCore
```

Send this request to this link with POST method.

```
https://jira.techno.ingenico.com/rest/api/2/issue/
```

After that get this kind of result:

```json
{"id":"705832","key":"LBO-73023","self":"https://jira.techno.ingenico.com/rest/api/2/issue/705832"}
```

## Local JIRA for testing

For testing added own jira server. (using 8282 as port number)

```sh
docker run -v jiraVolume:/var/atlassian/application-data/jira --name="jira" -d -p 8282:8080 atlassian/jira-software
```

After that you need to enter a license key to use it.

When installation complete, check jira version and look at the REST-API documentation.

https://docs.atlassian.com/software/jira/docs/api/REST/8.20.1/

In the profile page, add a personal access token.

Use your token with bearer header

```sh
curl -H "Authorization: Bearer ${TOKEN}" http://localhost:8282/rest/api/2/issue/SCRM-10
```

Now add auth to chore with giving this header and `POST` method.

https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/
