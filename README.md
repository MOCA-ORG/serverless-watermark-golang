# Serverless Watermark Node.js
Watermarking images with Golang, lambda, api gateway, serverless framework.

[Qiita](https://qiita.com/haxidoi/items/b278237d5aebfa889303)

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
