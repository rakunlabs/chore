<img src="docs/logo/chore.svg" height="120" />

Chore tool help to send request with templates and customizable flow diagram.

[intro page of ui](docs/info/intro.md)

## Usages

Connect to the chore UI with browser and add template, authentication and design own control flow.

Chore uses PostgreSQL database.

_First initialization user and password is **admin:admin**, changable with configuration_

### Configuration

```yaml
secret: thisisfordevelopmenttestsecret
user:
  name: admin
  password: admin
store:
  type: postgres
  schema: chore
  host: "127.0.0.1"
  port: "5432"
  user: postgres
  dbName: postgres
```

Secret is important for tokens, to generate own token, use one of this commands:

With openssl
```sh
openssl rand -base64 32 | tr -- '+/' '-_'
```

With linux shell
```sh
dd if=/dev/urandom bs=32 count=1 2>/dev/null | base64 | tr -d -- '\n' | tr -- '+/' '-_'; echo
```

### Template

Template is a text file format. Go template and sprig functions supported.

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

For testing in a playground try [repeatit.io](https://repeatit.io), this webapp developed by us.

### Auth

This give us information about secret headers after that use with request flow node.

With basic-auth(username and password) use this header `Authorization: Basic <base64 username:password>` but in 2FA status this cannot work so use Bearer Token(personal access token PAT) most of cases.

With Personal access token, generate token in the profile page and use with `Authorization: Bearer <TOKEN>`.

### Control

Flow diagram to create your algorithm in UI.

To start flow send request `/send` endpoint.  
Server will check __endpoint__ and __control__ values with your JSON/YAML payload.

## Informations

[JIRA template](docs/template/jira.md)

## Development

<details><summary>Build and run</summary>

Required services (PostgreSQL) before to run.

```sh
cd _example/chore
docker-compose up
# for close run
# docker-compose down
```

Generate swagger (don't need if you didn't change related codes)
```sh
./build.sh --swag
```

Run command
```sh
# ./build.sh --run
export CONFIG_FILE=_example/config/config.yml
go run cmd/chore/main.go
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

Build docker

```sh
./build.sh --docker-build
```

Run image
```sh
# run postgres before to start
# to get latest build image name
IMAGE_NAME=$(./build.sh --docker-name)
docker run -it --rm --name="chore" -p 8080:8080 \
  --add-host=postgres:$(docker network inspect bridge | grep Gateway | tr -d '" ' | cut -d ":" -f2) \
  -v ${PWD}/_example/config/docker.yml:/etc/chore.yml \
  ${IMAGE_NAME}
```

</details>

<details><summary>Dummy-Whoami Server for Test</summary>

```sh
docker run --rm -it --name="whoami" -p 9090:80 traefik/whoami
```

</details>

<details><summary>Fill tables</summary>

Get a token and set to `JWT_KEY` value.

```sh
export JWT_KEY=""
./data/post-template.sh
```

</details>
