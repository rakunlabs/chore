# Issue links

Template for update links

```json
{
  "type": {
    "name": "Relates"
  },
  "inwardIssue": {
    "key": "{{.key}}"
  },
  "outwardIssue": {
    "key": "{{.issue}}"
  }
}
```

Values, key is generated when creating service increment

```yml
key: LBO-73023
issue: LBO-71558
```

Send request to this link as POST

```
https://jira.techno.ingenico.com/rest/api/2/issueLink
```

## Issue link in different areas

__You should create link on both of servers.__

https://developer.atlassian.com/server/jira/platform/jira-rest-api-for-remote-issue-links/

Check remote link

__JIRA__

```sh
curl -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json"  https://jira.techno.ingenico.com/rest/api/2/issue/LBO-72594/remotelink | jq .
```

__MYITSM__

```sh
curl -H "Authorization: Basic ${TOKEN}" -H "Content-Type: application/json" https://myitsm.services.ingenico.com/rest/api/2/issue/ENCHG-36960/remotelink | jq .
```

GlobalID should be appId and issueId, if not UI cannot load it.

__NOTE__ appId cannot see very weird, exists in admin->systeminfo but I cannot see it as normal user

GlobalID useful for change or delete it.

Example

JIRA -> MYITSM

```json
{
  "globalId": "appId=55f82c15-c230-365f-b5ec-d9d6971f5d1f&issueId=196838",
  "application": {
    "type": "com.atlassian.jira",
    "name": "MyITSM"
  },
  "relationship": "relates to",
  "object": {
    "url": "https://myitsm.services.ingenico.com/browse/ENCHG-36960",
    "title": "ENCHG-36960",
    "icon": {},
    "status": {
      "icon": {}
    }
  }
}
```

Template

```json
{
  "globalId": "appId=55f82c15-c230-365f-b5ec-d9d6971f5d1f&issueId={{.myitsm.id}}",
  "application": {
    "type": "com.atlassian.jira",
    "name": "MyITSM"
  },
  "relationship": "relates to",
  "object": {
    "url": "https://myitsm.services.ingenico.com/browse/{{.myitsm.key}}",
    "title": "{{.myitsm.key}}",
    "icon": {},
    "status": {
      "icon": {}
    }
  }
}
```

Values

```yml
myitsm:
  key: ENCHG-36960
  id: 196838
```

Send to the JIRA issue as POST request

```
https://jira.techno.ingenico.com/rest/api/2/issue/LBO-72594/remotelink
```


---

MYITSM -> JIRA

```json
{
  "globalId": "appId=527ae2b4-d715-3525-b23d-78aa641cf419&issueId=696071",
  "application": {
    "type": "com.atlassian.jira",
    "name": "Ingenico Jira Techno"
  },
  "relationship": "relates to",
  "object": {
    "url": "https://jira.techno.ingenico.com/browse/LBO-72594",
    "title": "LBO-72594",
    "icon": {},
    "status": {
      "icon": {}
    }
  }
}
```

Template

```json
{
  "globalId": "appId=527ae2b4-d715-3525-b23d-78aa641cf419&issueId={{.jira.id}}",
  "application": {
    "type": "com.atlassian.jira",
    "name": "Ingenico Jira Techno"
  },
  "relationship": "relates to",
  "object": {
    "url": "https://jira.techno.ingenico.com/browse/{{.jira.key}}",
    "title": "{{.jira.key}}",
    "icon": {},
    "status": {
      "icon": {}
    }
  }
}
```

Values

```yml
jira:
  key: LBO-72594
  id: 696071
```

Send to the MYITSM issue as POST request

```
https://myitsm.services.ingenico.com/rest/api/2/issue/ENCHG-36960/remotelink
```
