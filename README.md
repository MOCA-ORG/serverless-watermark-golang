# Serverless Watermark Node.js
Watermarking images with Golang, lambda, api gateway, serverless framework.

## Dependence
### serverless
```shell
$ npm install -g serverless
```
### serverless(plugins)
```shell
$ sls plugin install -n serverless-offline
```

## Local test
```shell
$ sls offline start
```

## Deploy to AWS
```shell
# deploy to development environment.
$ make deploy

# deploy to production environment.
$ make deploy-prod
```
