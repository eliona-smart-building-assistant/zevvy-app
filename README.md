# App Template

This template is a part of the Eliona App SDK. It can be used to create an app stub for an Eliona environment.

## Configuration

The app needs environment variables and database tables for configuration. To edit the database tables the app provides an own API access.


### Registration in Eliona ###

To start and initialize an app in an Eliona environment, the app has to be registered in Eliona. For this, entries in database tables `public.eliona_app` and `public.eliona_store` are necessary.

This initialization can be handled by the `reset.sql` script.


### Environment variables

<mark>Todo: Describe further environment variables tables the app needs for configuration</mark>

- `CONNECTION_STRING`: configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). Otherwise, the app can't be initialized and started (e.g. `postgres://user:pass@localhost:5432/iot`).

- `INIT_CONNECTION_STRING`: configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db) for app initialization like creating schema and tables (e.g. `postgres://user:pass@localhost:5432/iot`). Default is content of `CONNECTION_STRING`.

- `API_ENDPOINT`:  configures the endpoint to access the [Eliona API v2](https://github.com/eliona-smart-building-assistant/eliona-api). Otherwise, the app can't be initialized and started. (e.g. `http://api-v2:3000/v2`)

- `API_TOKEN`: defines the secret to authenticate the app and access the Eliona API.

- `API_SERVER_PORT`(optional): define the port the API server listens. The default value is Port `3000`. <mark>Todo: Decide if the app needs its own API. If so, an API server have to implemented and the port have to be configurable.</mark>

- `LOG_LEVEL`(optional): defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-utils/blob/main/log/README.md). The default level is `info`.

### Database tables ###

<mark>Todo: Describe other tables if the app needs them.</mark>

The app requires configuration data that remains in the database. To do this, the app creates its own database schema `template` during initialization. To modify and handle the configuration data the app provides an API access. Have a look at the [API specification](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/app-template/develop/openapi.yaml) how the configuration tables should be used.

- `template.configuration`: Contains configuration of the app. Editable through the API.

- `template.asset`: Provides asset mapping. Maps broker's asset IDs to Eliona asset IDs.

**Generation**: to generate access method to database see Generation section below.


## References

### App API ###

The app provides its own API to access configuration data and other functions. The full description of the API is defined in the `openapi.yaml` OpenAPI definition file.

- [API Reference](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/app-template/develop/openapi.yaml) shows details of the API

**Generation**: to generate api server stub see Generation section below.

### Configuring the app ###

To use the app it is necessary to create at least one configuration. A configuration points to one Zevvy Login.

A minimum configuration that can used by the app's API endpoint `POST /configs` is:

```json
{
  "authRootUrl": "https://iam.zevvy.org/realms/zevvy-prod",
  "apiRootUrl": "https://api.zevvy.org",
  "clientId": "client",
  "clientSecret": "secret",
  "enable": true
}
```

The configuration may include an optional `refreshToken` property provided by the user. In its absence, the application initiates the authorization sequence by generating a new login URL. This URL is communicated to the requester through a user notification and can also be accessed by making a `GET /configs` request.

Following user login and API access verification, the application finalizes the configuration in the background. Upon successful acquisition of both an access and a refresh token, the user receives a notification via the Eliona frontend.

It's important to note that both the `clientSecret` and `refreshToken` properties are write-only for enhanced security. Thus, they cannot be fully retrieved through GET requests.

### Define assets attributes ###

To ensure data is successfully reported to Zevvy, the necessary asset attributes must be correctly configured via `PUT /asset-attributes` request.

```json
{
  "configId": 1,
  "assetId": 4711,
  "subtype": "input",
  "attributeName": "power"
}
```

For each attribute all the data stored in Eliona will be sent to Zevvy with the corresponding timestamp.

The asset's GAI is used as device reference and the name off the attribute as register reference. The values can be overwritten by the optional `deviceReference` and `registerReference` properties.

## Tools

### Generate API server stub ###

For the API server the [OpenAPI Generator](https://openapi-generator.tech/docs/generators/openapi-yaml) for go-server is used to generate a server stub. The easiest way to generate the server files is to use one of the predefined generation script which use the OpenAPI Generator Docker image.

```
.\generate-api-server.cmd # Windows
./generate-api-server.sh # Linux
```

### Generate Database access ###

For the database access [SQLBoiler](https://github.com/volatiletech/sqlboiler) is used. The easiest way to generate the database files is to use one of the predefined generation script which use the SQLBoiler implementation.

```
.\generate-db.cmd # Windows
./generate-db.sh # Linux
```

### Generate asset type descriptions ###

For generating asset type descriptions from field-tag-annotated structs, [asse-from-struct tool](https://github.com/eliona-smart-building-assistant/dev-utilities) can be used.
