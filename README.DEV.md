# 🛰 START in DEV MODE

> **NOTE**: This DApp has been tested on **Ubuntu 18.04** and on **Windows 10 with WSL** and Golang 1.18 was used.

## Table of Contents

- [Configuration file](#config_file)
- [Get Started](#get_started)
    * [Start DApp in dev mode](#dev_ways)
- [Tech and packages](#tech)
- [Architecture](#arch)
- [SWAGGER](#swagger)

## 🛠️️ Configuration file (conf.yaml) <a name="config_file"></a>

👉🏾 [The config file](/conf/conf.sample.unix.yaml)

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

### 🚀 Start DApp in dev mode <a name="dev_ways"></a>

Run:
```bash
go mod download
go mod vendor
```

If you make changes to the Endpoint you must generate Swagger API Spec:

![swagger doc](/docs/swagger.md)

Build:

```bash
go build -o dapp
```

#### 🌍 Environment variables

The environment variable is exported with the location of the server configuration file.

If you have 🐧Linux or 🍎Dash, run:

```bash
source ./setenv.sh
```

but if it's the windows cmd, run:

```bash
source ./setenv.cmd
```

#### Configure the dapp

Dapp config file:
```bash
cp conf/conf.sample.unix.yaml conf/conf.yaml
```

Network profile:
```bash
cp conf/cpp.sample.unix.yaml conf/cpp.yaml
```

Modify the crypto-config paths in __conf/cpp.yaml__ file, in our case it is `/home/user/fabric-folder/fabric-testnet-nano-without-syschannel` 

#### 🏃🏽‍♂️ Start the server

Before it is recommended that you read more about the server configuration file in the section 👉🏾 .

Run the server:

```bash
./dapp
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
