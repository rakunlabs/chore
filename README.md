# Chore

Chore tool help to send request with templates.

_Project in the baby steps everything can changable_

## Usages

Connect to the chore UI with browser and add template, authentication and binding.

### Template

Template is a text file format. Go template and sprig functions inside of it.

For example using some functions and flow inside of template.

```
ID: {{uuidv4}}
Name: {{.name | b64enc}}
{{if eq .name "golang" }}
Link: DeepCore
{{end}}

{{- range .specs}}
{{.name}} {{repeat .point "⭐"}}
{{- end}}
```

In here `name` is a key of a map or struct and it print value.

For testing in a playground try [repeatit.io](repeatit.io).

### Auth

This give us information about server URL and REST API specifications of target API.

id, headers, URL and method keywords exists.

### Bind

Combining Auth and Template with this table.  
When request is getting by `/send` endpoint, server will check auth and template with this entry.

## Examples

<details><summary>Test Server</summary>

Open test server with `go run _example/testServer/main.go`.

Add an auth entry to show this server.

```json
{"id":"secret","headers":"{\"Authorization\": \"Bearer <token>\"}","URL":"http://localhost:9090","method":"POST"}
```

Add an template and bind it.

```
hello {{.name}}
```

```json
{"id":"sendhi","authentication":"secret","template":"test"}
```

Now send values with curl or in the swagger documentation.

```sh
curl -X 'POST' \
  'http://localhost:3000/api/v1/send?key=sendhi' \
  -H 'accept: application/json' \
  -H 'Content-Type: */*' \
  -d 'name: test'
```

</details>

<details><summary>Create a ticket in JIRA</summary>

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
curl -H "Authorization: Bearer MjQ5Nzc3NTg3MjM4OosJndoCMilW9HAnAl4T2CfMEnbG" http://localhost:8282/rest/api/2/issue/SCRM-10
```

Now add auth to chore with giving this header and `POST` method.

https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/

</details>

## Development

Required services before to run.

<details><summary>Consul Setup</summary>

```sh
docker run -it --rm --name=dev-consul --net=host consul:1.10.4
```

If `chore` runs in the container set also `CONSUL_HTTP_ADDR` env variable.

</details>

Backend

```sh
./build.sh --swag
./build.sh --run
```

Frontend

```sh
cd _web
pnpm run dev -- --host
```

Build project

```sh
./build.sh --build-all
```
