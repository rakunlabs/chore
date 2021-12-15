# Chore

Chore tool help you send request with templates.

## Usages

### JIRA

For testing, you can use already installed jira or setup a jira from docker.

```sh
docker pull atlassian/jira-software:8.13.9
```

## Development

<details><summary>Consul Setup</summary>

```sh
docker run -it --rm --name=dev-consul --net=host consul:1.10.4
```

</details>