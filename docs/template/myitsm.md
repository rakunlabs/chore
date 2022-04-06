# Create Template for MYITSM

Change request creating usually 2 time, one for stag and one for prod.

## Get information with API

TOKEN is base64 encoded `username:password`, this token special you cannot success with your own.

```sh
curl -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" "https://myitsm.services.ingenico.com/rest/api/2/issue/ENCHG-36961" | jq .
```

Get types
```sh
curl -H "Authorization: Basic ${TOKEN}" -H "Content-Type: application/json" "https://myitsm.services.ingenico.com/rest/api/2/issue/createmeta" | jq .
```

Get issuetypes for ENCHG
```sh
curl -H "Authorization: Basic ${TOKEN}" -H "Content-Type: application/json" "https://myitsm.services.ingenico.com/rest/api/2/issue/createmeta/ENCHG/issuetypes" | jq .
```

Get all information of ENCHG (10002), if not exists in here don't send that value.
```sh
curl -H "Authorization: Basic ${TOKEN}" -H "Content-Type: application/json" "https://myitsm.services.ingenico.com/rest/api/2/issue/createmeta/ENCHG/issuetypes/10002" | jq .
```

## Template

Example for change type for Prod

```json
{
  "fields": {
    "project":
    {
       "key": "ENCHG"
    },
    "summary": "Release - PROD - V13 TRS - Change transaction date",
    "description": null,
    "issuetype": {
      "name": "Change"
    },
    "customfield_10106": {
        "value": "Release Change - BO Applications"
    },
    "customfield_11030": {
        "value": "Already performed"
    },
    "priority": {
        "name": "Normal"
    },
    "customfield_11029": {
        "value": "Financial Operations"
    },
    "customfield_11024": "DeepCore",
    "customfield_11027": "See Runbook",
    "customfield_11013": "Deepcore squad #Squad_FinOps_NL_DeepCore <Squad_FinOps_NL_DeepCore@epay.ingenico.com>",
    "customfield_11015": "N/A",
    "customfield_12500": {
      "value": "Global Online"
    },
    "customfield_11006": "See Runbook",
    "customfield_11008": "N/A",
    "customfield_11003": "N/A",
    "customfield_11106": {
      "value": "Release",
      "child": {
        "value": "PROD"
      }
    },
    "customfield_11101": "See Runbook",
    "customfield_11102": {
        "value": "Required & attached"
    },
    "customfield_11103": {
      "value": "Document required & attached (TER)"
    },
    "customfield_10004": {
      "value": "0 - Insignificant"
    }
  }
}
```

Template is

```json
{
  "fields": {
    "project":
    {
       "key": "ENCHG"
    },
    "summary": "Release - {{.release}} - {{.title}}",
    "description": {{or (.description | quote) "null"}},,
    "issuetype": {
      "name": "Change"
    },
    "customfield_10106": {
        "value": "Release Change - BO Applications"
    },
    "customfield_11030": {
        "value": "{{.uat}}"
    },
    "priority": {
        "name": "Normal"
    },
    "customfield_11029": {
        "value": "Financial Operations"
    },
    "customfield_11024": "{{.squad}}",
    "customfield_11027": "See Runbook",
    "customfield_11013": "{{.monitoring}}",
    "customfield_11015": "N/A",
    "customfield_12500": {
      "value": "Global Online"
    },
    "customfield_11006": "See Runbook",
    "customfield_11008": "N/A",
    "customfield_11003": "N/A",
    "customfield_11106": {
      "value": "Release",
      "child": {
        "value": "{{.release}}"
      }
    },
    "customfield_11101": "See Runbook",
    "customfield_11102": {
        "value": "Required & attached"
    },
    {{if .risk -}}
    "customfield_10004": {
        "value": "{{.risk}}"
    },
    {{- end}}
    "customfield_11103": {
      "value": "Document required & attached (TER)"
    }
  }
}
```

Values

```yml
title: V13 TRS - Change transaction date
squad: DeepCore
description: null
monitoring: "Deepcore squad #Squad_FinOps_NL_DeepCore <Squad_FinOps_NL_DeepCore@epay.ingenico.com>"
# or "Pre-Prod/STAG"
release: PROD
# or null
risk: "0 - Insignificant"
# or "Yes"
uat: "Already performed"
```

Send post request this link as POST

```
https://myitsm.services.ingenico.com/rest/api/2/issue
```
