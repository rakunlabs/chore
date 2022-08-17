<img src="docs/logo/chore.svg" height="120" />

Chore tool help to send request with templates and customizable flow diagram.

__-__ [info page of ui](docs/info/intro.md)

Template playground go to [repeatit.io](https://repeatit.io)

If you need any feature, find a bug or fixing something send pull request or open issue we will handle it.

## Fast Start and Use

```sh
docker run -d --name postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust postgres:14.5-alpine
docker run -it --rm --name chore -p 8080:8080 -e STORE_HOST=172.17.0.1 -e STORE_SCHEMA=public ghcr.io/worldline-go/chore:latest
```

Open browser and go to http://localhost:8080

Login with `admin:admin` and after that create your flow, template and run it.

![chore](./docs/info/chore.gif)

## Usages

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
  # password: test
  dbName: postgres
  timeZone: UTC
  # also you can set with DSN name, if DSN name exists other values not using
  # dbDataSource: "postgres://postgres@127.0.0.1:5432/postgres?application_name=testdb"

# migrate same as store and copy undefined part in store value
migrate:
  password: formigration
  user: migration

# BasePath just required for swagger ui, this tool use relative paths at all
# basePath: /chore/ # to set mywebsite.com/chore/
# host: 0.0.0.0 # default
# port: 8080 # default
# logLevel: info # default
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

__WARN__ when secret changed, all previous tokens not usable after that.

Set config file path to `CONFIG_FILE` environment variable.

Chore can get your configurations from vault, consul, file or with environments values.

To work with vault and consul set `PREFIX_VAULT` and `PREFIX_CONSUL` to show the path of the config file and `APP_NAME` default is __chore__. Details check [igconfig](https://github.com/worldline-go/igconfig) library to see how it works.

And run chore on container or binary.

Connect to the chore UI with browser and add template, authentication and design own control flow.

<details><summary>Template, Auth, Control information</summary>

### Template

Template is a text file format. `Go template` and `sprig` functions supported.

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

With basic-auth(username and password) use this header `Authorization: Basic <base64 username:password>` but in 2FA status this cannot work so use Bearer Token(personal access token PAT) most of cases or ask IT to get new user which can work with api.

With Personal access token, generate token in the profile page and use with `Authorization: Bearer <TOKEN>`.

### Control

Flow diagram to create your algorithm in UI.

To start flow send request `/send` endpoint as POST request.  
Server will check __endpoint__ and __control__ values with your __JSON/YAML__ payload.

Example: (generate token in token section of chore)

```sh
curl -X POST -H "Authorization: Bearer ${TOKEN}" -d 'name: deepcore' "http://localhost:8080/api/v1/send?control=try&endpoint=test"
```

Or you can send as json value `-d '{"name":"deepcore"}'`

Or send file directly, (when sending yaml format always use binary format due to yaml has new line and ascii format not hold that values)

```sh
curl -X POST -H "Authorization: Bearer ${TOKEN}" --data-binary @values.yml "http://localhost:8080/api/v1/send?control=try&endpoint=test"
```

</details>

## Development

<details><summary>Build and run</summary>

### Run

Required services (PostgreSQL) before to run.

```sh
cd _example/chore
docker-compose up
# for close run
# docker-compose down
```

Run command
```sh
# ./build.sh --run
# config file can be TOML, YAML or JSON
export CONFIG_FILE=_example/config/config.yml
go run cmd/chore/main.go
```

Frontend
```sh
cd _web
pnpm run dev -- --host
```

After this step just go to the `localhost:3000` address.

__NOTE__ frontend(`localhost:3000`) has proxy and `/api` path request goes to the server.

### Build

#### Build with goreleaser

```sh
goreleaser release --snapshot --rm-dist
```

### Build with script

Generate swagger (don't need if you didn't change related codes)
```sh
./build.sh --swag
```

Build project to generate binary
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

Use chore's record script to download/opload operation

Before to run script export __JWT_KEY__ variable with own chore token.

Change `-h` (help) parameter to any arguments of the shell script.

```sh
export JWT_KEY=""
curl -fksSL https://raw.githubusercontent.com/worldline-go/chore/main/data/record.sh | bash -s -- -h
```

Or first download it and after run.

```sh
curl -O -fksSL https://raw.githubusercontent.com/worldline-go/chore/main/data/record.sh && chmod +x record.sh
```

Example arguments
```sh
# download just one item
--url http://localhost:8080 --mode download --auth jira
# update all auths, controls and templates files
--url http://localhost:8080 --mode download --auths --controls --templates
# upload all auths folder
--url http://localhost:8080 --mode upload --auths
# upload just one item
--url http://localhost:8080 --mode upload --template confluence/ter
```

Get temporary JWT key with username and password

```sh
export JWT_KEY="$(curl -fksSL -u admin:admin http://localhost:8080/api/v1/login?raw=true)"
```

</details>

## Todo

- [ ] Activate group information
- [ ] Support custom method entries
