# 📙 OpenAPI Specification

## To install Swagger api documentation running this command:

```shell
go install github.com/swaggo/swag/cmd/swag@v1.8.6
```

## To generate the OpenAPI specification run the fallowing command:

```shell
swag init --md docs/md_endpoints
```
[read more...](https://github.com/swaggo/swag/issues/817)

The current OpenAPI version used in this project is __2.2__. The integration package has not being ported to the 3 version yet.
For visiting the documentation open a browser in ``` http://127.0.0.1:8080/swagger/index.html ```.

___