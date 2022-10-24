Chore tool designed to help you manage all of your external requests.

Create a control flow with declaring templates and authentications and call it with REST-API.

![diagram](./assets/chore-diagram.svg)

## Control Flow

Flow always run with an __endpoint__ and stop until no run node left.
It is continue to run even return a respond.

If you use more than one respond node in flow, it will return in first message.

### Endpoint

Endpoint is starting point of the control flow.  
This is required name for sending request to chore to start flow.

If public is enabled, it can callable without authentication.

#### INPUT

Input is bytes of payload, usually values of request to chore.

#### OUTPUT

Directly send to bytes to other nodes.

```
 ┌─────────────────────────┐
 │ ENDPOINT                │
 ├─────────────────────────┤
 │ Enter endpoint name    ┌┼┐
 │ ┌────────────────────┐ └┼┘
 │ │                    │  │
 │ └────────────────────┘  │
 │ Methods                 │
 │ ┌────────────────────┐  │
 │ │                    │  │
 │ └────────────────────┘  │
 │ Public              │ │ │
 └─────────────────────────┘
```

### Template

Go template with sprig functionality and some extra functions.  
Usually this is a json with some values.

Template name needs to set before trigger to control flow. If not it just respond error.

For go template playground try this: __[repeatit.io](https://repeatit.io)__

#### INPUT

Input is bytes of previous node.

#### OUTPUT

Rendered bytes.

```
 ┌─────────────────────────┐
 │ TEMPLATE                │
 ├─────────────────────────┤
┌┼┐Enter template name    ┌┼┐
└┼┘┌────────────────────┐ └┼┘
 │ │                    │  │
 │ └────────────────────┘  │
 └─────────────────────────┘
```

### Request

Send http request. Set URL, method and headers with previously declared an auth header.

In URL, method and headers can be usuable with go template. We use `interface{}` in template, so it could be any type.

POST method is default when not set any method.

If `V` is connected and a request comes, it first wait `V` value to come to the request node.  
When V value is set, you can use that one for multiple requests.

Chore tries to set all nodes as pure and usable again.  
So if you want to pure function call just call `setValue` function in the main function it will use that value to render go templates.

#### INPUT

`V-` Values as yaml/json bytes form for fill URL, method and headers' template values.  
`_-` Input is binary bytes of payload, usually values of request to chore.
#### OUTPUT

`F-` Returned body as bytes but when status code not between [100-399].  
`_-` Returned body as bytes.  

Request can be used with respond node to return response directly of request's result.

```
 ┌───────────────────────────┐
 │ REQUEST                   │
 ├───────────────────────────┤
 │ Enter request url         │
 │ ┌────────────────────┐    │
 │ │                    │    │
 │ └────────────────────┘    │
 │ Enter method              │
┌┼┐┌────────────────────┐   ┌┼┐
│V││ POST               │   │F│
└┼┘└────────────────────┘   └┼┘
┌┼┐Enter additional headers ┌┼┐
└┼┘┌────────────────────┐   └┼┘
 │ │                    │    │
 │ └────────────────────┘    │
 │ Enter auth                │
 │ ┌────────────────────┐    │
 │ │                    │    │
 │ └────────────────────┘    │
 └───────────────────────────┘
```

### Script

Javascript code (ES5.1) for parsing, editing and managing control flow.

`Open Editor` button opens code editor window to write code better.

One `main` function must be set.

If function throw an error (`throw data`), script continue flow on false path.

<u>Predefined functions:</u>  
`toObject` convert byte to object  
`toString` convert byte to string  
`sleep` parameter such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
`setValue` set value for use in future go template.

A json/yaml entries automatically converting to the object/array not need to convert and not need to convert back to string.  
Functions just for corner cases not need to use.

Input entry could be more than one, that mean you can connect more than one node to script node and it run on last entry comes.

Use defination as `function main(input1, input2, input3)` to get more entry values.  
It will wait all entry values to come and last one will trigger the script.  
Usually use this when you want to combine request results.

Script-Editor also has playground to test your code.

#### INPUT

Bytes from previous nodes. Input count can be set when adding to flow.

First input is first argument of the code.

#### OUTPUT

Returned value.

__NOTE__ flow automatically convert JS object to byte when sending other nodes.

```
 ┌───────────────────────────┐
 │ Script                    │
 ├───────────────────────────┤
 │ |Open Editor            │┌┼┐
┌┼┐┌───────────────────────┐│F│
└┼┘│function main(data){   │└┼┘
 │ │  return data;         │┌┼┐
 │ │}                      │└┼┘
 │ └───────────────────────┘ │
 └───────────────────────────┘
```

### Respond

Return input to the control flow caller.

If at least one respond to the flow, caller will wait to return of respond.

If respond in flow not exists, `201` code return to caller.

#### INPUT

Bytes from previous nodes. Input count can be set when adding to flow.

First input is first argument of the code.

Respond include headers if you want to send as html just add header in there.

#### OUTPUT

Not exists.

```
 ┌┬┬─────────────────────────┐
 │ │Respond                  │
 └┴┴─────────────────────────┘

 + headers and get respond data
```

### Log

Print message in chore server. It helps to record something and debugging.

Log node don't modify input just send to output but type is changing so not usable with special nodes.

Select print data for printing input data in log message.

#### INPUT

Bytes from other nodes.

#### OUTPUT

Input value.

```
 ┌───────────────────────────┐
 │ Log                       │
 ├───────────────────────────┤
 │ Message                   │
 │ ┌────────────────────┐    │
 │ │                    │    │
┌┼┐└────────────────────┘   ┌┼┐
└┼┘ Log Level               └┼┘
 │ ┌────────────────────┐    │
 │ │ Debug           \/ │    │
 │ └────────────────────┘    │
 │───────────────────────────│
 │ Print data            │ │ │
 └───────────────────────────┘
```

### Email

Send email with a playload data.

_to_, _cc_, _bcc_ values should be comma or space seperated like (user1@example.com, user2@example.com)

Email server and authentication should be set with an admin account in chore.

#### INPUT

`V-` Values as yaml/json bytes form for fill all values.  
`_-` Input is binary bytes of payload, usually values of request to chore.

#### OUTPUT

Not exists.

```
 ┌───────────────────────────┐
 │ Email                     │
 ├───────────────────────────┤
 │ From                      │
 │ ┌───────────────────────┐ │
 │ │                       │ │
 │ └───────────────────────┘ │
 │ To                        │
 │ ┌───────────────────────┐ │
┌┼┐│                       │ │
│V│└───────────────────────┘ │
└┼┘CC                        │
┌┼┐┌───────────────────────┐ │
└┼┘│                       │ │
 │ └───────────────────────┘ │
 │ Bcc                       │
 │ ┌───────────────────────┐ │
 │ │                       │ │
 │ └───────────────────────┘ │
 │ Subject                   │
 │ ┌───────────────────────┐ │
 │ │                       │ │
 │ └───────────────────────┘ │
 └───────────────────────────┘
```

### IF

If case want a statement and input value defined as `data` value.

`data` is a special name to represent input value. It can be any type what you give.

#### INPUT

Bytes from previous nodes.

#### OUTPUT

Input value.

```
 ┌───────────────────────────┐
 │ IF                        │
 ├───────────────────────────┤
 │ Expression               ┌┼┐
┌┼┐┌───────────────────────┐│F│
└┼┘│                       │└┼┘
 │ │data > 0               │┌┼┐
 │ │                       ││T│
 │ └───────────────────────┘└┼┘
 └───────────────────────────┘
```

### For

For loop want a statement and should an array.  
Input value defined as `data` object.

For loop call output branch with iterating array.

#### INPUT

Bytes from previous nodes.

#### OUTPUT

For each of returned value.

```
 ┌───────────────────────────┐
 │ For Loop                  │
 ├───────────────────────────┤
 │ Return an array           │
┌┼┐┌───────────────────────┐┌┼┐
└┼┘│data                   │└┼┘
 │ └───────────────────────┘ │
 └───────────────────────────┘
```

### Note

Record some information to explain flow.

```
 ┌───────────────────────────┐
 │┌─────────────────────────┐│
 ││Important!               ││
 ││                         ││
 ││                         ││
 ││                         ││
 │└─────────────────────────┘│
 └───────────────────────────┘
```
