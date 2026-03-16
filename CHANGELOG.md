# [3.0.0](https://github.com/shivanshkc/rosenbridge/compare/v2.0.0...v3.0.0) (2026-03-16)


### Bug Fixes

* **ci:** use latest golangci lint action ([d95bea0](https://github.com/shivanshkc/rosenbridge/commit/d95bea03beb49c5886295dd953031cb337604a79))
* **connect-api:** accept credentials as query params, skip CORS ([afcf190](https://github.com/shivanshkc/rosenbridge/commit/afcf190dc714d1ed15c1b5d8b953e439124e6aff))
* file database uses MarshalIndent ([34037e0](https://github.com/shivanshkc/rosenbridge/commit/34037e0e91cbbaa047c07ad52e61550a54a61101))
* **http:** reduce MaxHeaderBytes from 2 MB to 64 KB ([a73877a](https://github.com/shivanshkc/rosenbridge/commit/a73877a480cb1cae7b5faf045444e40f4a7b0746))
* **lint:** explicitly ignore connection close errors in test ([617140e](https://github.com/shivanshkc/rosenbridge/commit/617140e50b24a10d7df2d65b37281b03f98bbf71))
* make frontend config controllable from backend ([8ccacfc](https://github.com/shivanshkc/rosenbridge/commit/8ccacfc95c492ffa54a6c3abd304cd0dcb3ffa3a))


### Features

* add spa server ([0752b95](https://github.com/shivanshkc/rosenbridge/commit/0752b95903e92913b966b743812ceab80f85de4b))
* **ci:** add ci-cd pipeline ([6954c5c](https://github.com/shivanshkc/rosenbridge/commit/6954c5c3a1324d6669fa1862328b4d3ee85cbe7c))
* **ci:** containerization ([c691073](https://github.com/shivanshkc/rosenbridge/commit/c6910734651dc63f367a929ebbb36ad275d959e9))
* **connect-api:** add connect API ([c8603be](https://github.com/shivanshkc/rosenbridge/commit/c8603be0906771c49d5d622cca5aabf413c1fd3a))
* **create-user-api:** add body parsing and validation ([94f7ecd](https://github.com/shivanshkc/rosenbridge/commit/94f7ecdf2570b10c686e594560b712e82ec6a34b))
* **create-user-api:** add database layer ([6357ebf](https://github.com/shivanshkc/rosenbridge/commit/6357ebf14ccf149d85651ff1b286a4c5a9ae3c89))
* **send-message-api:** add send message api ([8392897](https://github.com/shivanshkc/rosenbridge/commit/83928976c35679d0d5a81b39f851e63f6c239b77))
* **web-client:** add browser client ([b621448](https://github.com/shivanshkc/rosenbridge/commit/b6214484c5abc6cdd82b619259f77d4641dac757))
* **web-client:** add new client ([d965bb4](https://github.com/shivanshkc/rosenbridge/commit/d965bb44de086091a620dd4d79ce68f68411c67c))
* **ws-manager:** outsource websocket logic to internal/ws package ([c86fd96](https://github.com/shivanshkc/rosenbridge/commit/c86fd963f2953fb958fb00fc5f75b779da54c11e))


### BREAKING CHANGES

* API contract modified
