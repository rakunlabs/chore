# Chore

Chore tool help you send request with templates.

_Project in the development step everything can changable_

## Usages

Connect to the chore UI and add template, authentication and binding.

### Template

Template is a text file format. Go template and sprig functions inside of it.

For example

```txt
Hello {{.name}}
```

In here `name` is a key of a map or struct and it print value.

For testing in a playground try [repeatit.io](repeatit.io).

### Auth

This give us information about server URL and REST API specifications of target API.

id, headers, URL and method keywords exists.

### Bind

Combining Auth and Template with this table.  
When request is getting by `/send` endpoint, server will check auth and template with this entry.

## Example usage



## Development

Required services before to run.

<details><summary>Consul Setup</summary>

```sh
docker run -it --rm --name=dev-consul --net=host consul:1.10.4
```

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
