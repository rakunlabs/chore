Chore tool designed to help you manage all of your external requests.

Create a control flow with declaring templates and authentications and call it with REST-API.

![diagram](./assets/chore-diagram.svg)

## Control Flow

Flow always run with an __endpoint__ and stop until no run node left.
It is continue to run even return a respond.

### Endpoint

Endpoint is starting point of the control flow.  
This is required name for sending request to chore to start flow.

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
 └─────────────────────────┘
```

### Template

Go template with swag functionality. Usually this is a json with some values.

Template name needs to set before trigger to control flow.

For testing goto __[repeatit.io](https://repeatit.io)__

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

In URL, method and headers can be usuable with go template. Values should be send on V input as json/yaml.

POST method is default when not set any method.

#### INPUT

`V-` Values as yaml/json bytes form for fill URL, method and headers' template values.  
`_-` Input is binary bytes of payload, usually values of request to chore.
#### OUTPUT

`F-` Returned body as bytes but when status code not between [100-399].  
`_-` Returned body as bytes.  

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

#### OUTPUT

Not exists.

```
 ┌┬┬─────────────────────────┐
 │ │Respond                  │
 └┴┴─────────────────────────┘
```

### Log

Print message in chore server. It helps to record something and debugging.

Log node don't modify input just send to output.

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

_to_, _cc_, _bcc_ values should be comma seperated like (user1@example.com, user2@example.com)

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

If case want a statement and input value defined as `data` object.

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
