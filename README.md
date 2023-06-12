# README
![](static/logo.jpeg)

Host a simple Go web application in Netlify

**Live Preview**

[![Netlify Status](https://api.netlify.com/api/v1/badges/46442f41-b1d8-4ff6-b3a0-6ea3b301515b/deploy-status)](https://app.netlify.com/sites/depoly-go/deploys)

[![Vercel Status](https://vercel.badge.hnest.eu.org/api/daviddenghui/depoly-go)](https://vercel.com/daviddenghui/depoly-go/deployments "Vercel Status")<br>

## Write a function

The [Netlify docs](https://docs.netlify.com/functions/build-with-go/) for golang are a great starting point.

## Build the binary
Push code to GitHub and instruct the Netlify CI on how to use the source code to build.

## Deploy to Netlify
Configure the settings in a `netlify.toml` file in the root of project.

```toml
[build]
  command = "${build command}"
  functions = "functions"
  publish = "${static html file}"
```