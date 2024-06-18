# Zevvy User Guide

### Introduction

> The Zevvy app provides integration and synchronization between Eliona and Zevvy services.

## Overview

This guide provides instructions on configuring, installing, and using the Zevvy app to manage resources and synchronize data between Eliona and Zevvy services.

## Installation

Install the Zevvy app via the Eliona App Store.

## Configuration

The Zevvy app requires configuration through Eliona’s settings interface. Below are the general steps and details needed to configure the app effectively.

### Registering the app in Zevvy Service

Create credentials in Zevvy Service to connect the Zevvy services from Eliona. All required credentials are listed below in the [configuration section](#configure-the-zevvy-app).  

Login into the Zevvy console and go to the settings page. Create a new client id and secret. 

### Configure the Zevvy app 

Configurations can be created in Eliona under `Apps > Zevvy > Settings` which opens the app's [Generic Frontend](https://doc.eliona.io/collection/v/eliona-english/manuals/settings/apps). Here you can use the appropriate endpoint with the POST method. Each configuration requires the following data:

| Attribute         | Description                                            |
|-------------------|--------------------------------------------------------|
| `authRootUrl`     | Root URL for the authentication process.               |
| `apiRootUrl`      | Root URL for the API access.                           |
| `clientId`        | Client ID for API access created in Zevvy console.     |
| `clientSecret`    | Client secret for API access created in Zevvy console. |
| `enable`          | Flag to enable or disable this configuration.          |
| `refreshInterval` | Interval in seconds for data synchronization.          |
| `requestTimeout`  | API query timeout in seconds.                          |

Example configuration JSON:

```json
{
  "authRootUrl": "https://iam.zevvy.org/realms/zevvy-prod",
  "apiRootUrl": "https://api.zevvy.org",
  "clientId": "123abc",
  "clientSecret": "s3cr3t",
  "enable": true,
  "refreshInterval": 60,
  "requestTimeout": 120
}
```

After the technical basics of the app have been configured, the app needs further information about which metrics should be reported to Zevvy.
To do this, it is necessary to configure the assets and the corresponding measurement attribute.

This can be done using the appropriate endpoint with the POST method. Each asset definition requires the following data:

| Attribute           | Description                                                                                                   |
|---------------------|---------------------------------------------------------------------------------------------------------------|
| `configId`          | The measurement is sent with this configuration..                                                             |
| `assetId`           | Measurement is taken from this asset.                                                                         |
| `subtype`           | Measurement data has this subtype.                                                                            |
| `attributeName`     | Name of the measurement attribute.                                                                            |
| `deviceReference`   | Name of the measurement's device reference in Zevvy. (Optionally, default device reference is asset's GAI)    |
| `registerReference` | Name of the measurement's register reference in Zevvy. (Optionally, register reference is the attribute name) |

Example JSON to configure a measurement data point for Zevvy

```json
{
  "configId": 1,
  "assetId": 4711,
  "subtype": "input",
  "attributeName": "power"
}
```

## Zevvy 

Once configured, the app starts sending periodically measurements taken from the configured assets and attributes to Zevvy.
