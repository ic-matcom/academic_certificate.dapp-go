# 🛰 GitHub Template DApp for Hyperledger Fabric

DApp to communicate with Hyperledger Fabric Blockchain Network

> **NOTE**: This DApp has been tested on **Ubuntu 18.04** and on **Windows 10 with WSL** and Golang 1.18 was used.

👉🏾 Read ![this doc](/README.DEV.md) to to start DApp in dev mode.

## Table of Contents

- [Configuration file](#config_file)
- [Get Started](#get_started)
    * [Deployment ways (2 ways)](#deploy_ways)
        - [Docker way](#docker_way)
        - [Manual way](#manual_way)
- [Tech and packages](#tech)
- [Architecture](#arch)
- [SWAGGER](#swagger)

## 🛠️️ Configuration file (conf.yaml) <a name="config_file"></a>

👉🏾 [The config file](/conf/conf.yaml)

| Param       | Description                                               | default value                 |
|-------------|-----------------------------------------------------------|-------------------------------|
| APIDocIP    | IP to expose the api (unused)                             | 127.0.0.1                     |
| DappPort    | app PORT                                                  | 7001                          |
| CronEnabled | active the cron job                                       | true                          |
| EveryTime   | time interval (in seconds) that the cron task is executed | 300 seconds (every 5 minutes) |

## ⚡ Get Started <a name="get_started"></a>

Download the github.template-fabric.dapp-go project and move to root of project:

```bash
git clone https://github.com/kmilodenisglez/github.template-fabric.dapp-go.git && cd github.template-fabric.dapp-go 
```

### 🚀 Deployment ways (2 ways)  <a name="deploy_ways"></a>

You can start the server in 2 ways, the first is using **docker** and **docker-compose** and the second is **manually**

#### 📦 Docker way <a name="docker_way"></a>

You will need docker and docker-compose in your system.

To builds Docker image from Dockerfile, run:

```bash
docker build --no-cache --force-rm --tag app_restapi .
```

Use docker-compose to start the container:

```bash
docker-compose up
```

#### 🔧 Manual way  <a name="manual_way"></a>

Run:

```bash
go mod download
go mod vendor
```

If you make changes to the Endpoint you must generate Swagger API Spec:

![swagger doc](/docs/swagger.md)

Build:

```bash
go build
```

#### 🌍 Environment variables

The environment variable is exported with the location of the server configuration file.

If you have 🐧Linux or 🍎Dash, run:

```bash
export SERVER_CONFIG=$PWD/conf/conf.yaml
```

but if it's the windows cmd, run:

```bash
set SERVER_CONFIG=%cd%/conf/conf.yaml
```

#### 🏃🏽‍♂️ Start the server

Before it is recommended that you read more about the server configuration file in the section 👉🏾 .

Run the server:

```bash
./github.template-fabric.dapp-go
```

and visit the swagger docs:

> http://localhost:7001/swagger/index.html

![swagger ui](/docs/images/swagger-ui.png)

You can then authenticate and test the remaining endpoints.

### 🧪 Unit or End-To-End Testing

Run:

```bash
go test -v
```

## 🔨 Tech and packages <a name="tech"></a>

* [Iris Web Framework](https://github.com/kataras/iris)
* [validator/v10](https://github.com/go-playground/validator)
* [govalidator](https://github.com/asaskevich/govalidator)
* [gocron](https://github.com/go-co-op/gocron)
* [swag](https://github.com/swaggo/swag)
* [Docker](https://docs.docker.com)
* [docker-compose](https://docs.docker.com/compose/)

## 📐 Architecture <a name="arch"></a>

This project has 3 layer :

- Controller Layer (Presentation)
- Service Layer (Business)
- Repository Layer (Persistence)

| Tag  | Path                                                         | Layer      |
|------|--------------------------------------------------------------|------------|
| Auth | [end_auth.go](/api/endpoints/end_auth.go)                    | Controller |
| Dapp | [end_dapp.go](/api/endpoints/end_dapp.go)                    | Controller |
|      |                                                              |            |
| Auth | [svc_authentication.go](/service/auth/svc_authentication.go) | Service    |
| Dapp | [svc_dapp.go](/service/svc_dapp.go)                          | Service    |
|      |                                                              |            |
| Dapp | [repo_dapp.go](/repo/repo_dapp.go)                           | Repository |

## 📐 Swagger <a name="swagger"></a>

Read ![swagger doc](/docs/swagger.md)