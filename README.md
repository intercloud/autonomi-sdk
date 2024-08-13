# Autonomi SDK

Autonomi allows to deploy automatically and asynchronously cloud resources.

## Getting started

### Prerequisites

Go 1.22 or higher

### Install Autonomi Go SDK

1. Use `go get` to install the latest version of the Milvus Go SDK and dependencies:

   ```shell
   go get -u https://github.com/intercloud/autonomi-sdk
   ```

2. Include the Autonomi Go SDK in your application:

```go
import autonomisdk "github.com/intercloud/autonomi-sdk"

//...other snippet ...
client, err := autonomisdk.NewClient(
    terms_and_conditions,
    autonomisdk.WithHTTPClient(&http.Client{}),
    autonomisdk.WithHostURL(hostURL),
    autonomisdk.WithPersonalAccessToken(personal_access_token),
)
if err != nil {
    // handle error
}
defer client.Close()

workspace, err := client.CreateWorkspace(ctx, payload)
if err != nil {
    // handle error
}
```

### Resources

Autonomi SDK allows to :

- Create, Read, Update and Delete a **Workspace**
- Create, Read, Update and Delete a **Node**
- Create, Read, Update and Delete a **Transport**

### To Do

- Create, Read, Update and Delete a **Attachment**
