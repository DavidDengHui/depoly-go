# Some tips about this project

Host a simple Go web application in Netlify and Vercel

**Live Preview**

| **NUM** | **Host** | **Status** | **Preview** |
|:-------:|:--------:|:----------:|:----------:|
| **1**   | [Vercel](https://vercel.com/)   | [![Vercel Status](https://vercel.badge.hnest.eu.org/api/daviddenghui/depoly-go)](https://vercel.com/daviddenghui/depoly-go/deployments "Vercel Status")        | <https://api3.hnest.eu.org/> |
| **2**   | [Netlify](https://www.netlify.com/)  | [![Netlify Status](https://api.netlify.com/api/v1/badges/46442f41-b1d8-4ff6-b3a0-6ea3b301515b/deploy-status)](https://app.netlify.com/sites/depoly-go/deploys) | <https://api2.hnest.eu.org/> |


## Write a function

The [Vercel docs](https://vercel.com/docs/concepts/functions/serverless-functions/runtimes/go#) and [Netlify docs](https://docs.netlify.com/functions/create/?fn-language=go) for golang are a great starting point.

The `functions directory` deployed to `Vercel` is [***/api***](api).

The `functions directory` deployed to `Netlify` is [***/cmd***](cmd).

The shared `static pages directory` is [***/public***](public).

## Build the binary
Push code to GitHub and instruct the `Vercel` and `Netlify` CI on how to use the source code to build.

## Deploy to Vercel
- Configure the settings in a [***vercel.json***](vercel.json) file in the root of project. <br>For example:

```yaml
{
  "version": 2,
  "build": {
    "env": {
      "GO_BUILD_FLAGS": "-ldflags '-s -w'"
    }
  },
  "builds": [
    {
      "src": "public/**",
      "use": "@vercel/static"
    },
    {
      "src": "api/**/*.go",
      "use": "@vercel/go"
    }
  ],
  "routes": [
    {
      "src": "/(.*\\..*)",
      "dest": "/public/$1"
    },
    {
      "src": "/([^.]*)",
      "dest": "/api/$1.go"
    }
  ]
}
```

- Configure the settings in a [***.vercelignore***](.vercelignore) file in the root of project. <br>For example:

```shell
# Golang components used in Netlify build
build.sh
.go-version
go.mod
go.sum
```

- You can specify the [*Ignored Build Step*](https://vercel.com/docs/concepts/projects/overview#ignored-build-step) that need to be ignored in the project `Settings -> Git` of `Vercel`. <br>For [example](https://vercel.com/daviddenghui/depoly-go/settings/git), specifying only the submission of `/api` and `/public` folder changes for deployment.

```shell
git diff --quiet HEAD^ HEAD ./api ./public
```

## Deploy to Netlify
- Configure the settings in a [***netlify.toml***](netlify.toml) file in the root of project. <br>For example:

```toml
[build]
  command = "chmod a+x build.sh && ./build.sh"
  functions = "functions"
  publish = "public"
  ignore = "git diff --quiet $CACHED_COMMIT_REF $COMMIT_REF ./cmd ./public"

[build.environment]
  GO_IMPORT_PATH = "github.com/DavidDengHui/depoly-go"
  GO111MODULE = "on"

[[redirects]]
  from = "/api/*"
  to = "/.netlify/functions/gateway/:splat"
  status = 200
```

- You can customize the `Netlify` build step in the [***build.sh***](build.sh). <br>For example:

```shell
set -euxo pipefail

mkdir -p "$(pwd)/functions"
rm -rf "$(pwd)/api"
GOBIN=$(pwd)/functions go install ./...
chmod +x "$(pwd)"/functions/*
go env
```