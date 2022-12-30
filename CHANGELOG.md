# [2.0.0](https://github.com/shivanshkc/rosenbridge/compare/v1.1.0...v2.0.0) (2022-12-30)


### Bug Fixes

* **core:** add missing error logs and websocket origin check ([465ed0a](https://github.com/shivanshkc/rosenbridge/commit/465ed0af4460cfb16d5feb51e8b6ed6c3acbfa6f))
* **core:** bug fixes ([6dab53b](https://github.com/shivanshkc/rosenbridge/commit/6dab53bfac8d35d41c7b5115e08cbb70b83dc50b))


### Features

* **core:** accept client ID in route params for browser compatibility ([741db5a](https://github.com/shivanshkc/rosenbridge/commit/741db5adbd6eb99b5b96af8cb261f949b0fc71f9))
* **core:** client_id is now a query parameter ([11f0b0f](https://github.com/shivanshkc/rosenbridge/commit/11f0b0fb75e4700152b54d3a960f4086c70a084b))


### BREAKING CHANGES

* **core:** GetBridge API change
* **core:** connect api interface change

# [2.0.0](https://github.com/shivanshkc/rosenbridge/compare/v1.1.0...v2.0.0) (2022-12-29)


### Features

* **core:** accept client ID in route params for browser compatibility ([741db5a](https://github.com/shivanshkc/rosenbridge/commit/741db5adbd6eb99b5b96af8cb261f949b0fc71f9))


### BREAKING CHANGES

* **core:** connect api interface change

# [1.1.0](https://github.com/shivanshkc/rosenbridge/compare/v1.0.1...v1.1.0) (2022-12-22)


### Bug Fixes

* **ci:** cd-test does not need lint ([d66b3d3](https://github.com/shivanshkc/rosenbridge/commit/d66b3d32918e3480993f787618ee45b808945c9b))
* **ci:** config location fix ([cfd819b](https://github.com/shivanshkc/rosenbridge/commit/cfd819b0e60d57f5ffbcfb8b3153676da86912bf))
* **ci:** fix config location ([df99e6a](https://github.com/shivanshkc/rosenbridge/commit/df99e6a8beb8fcbecbe50a57edd456d6eff99bbe))
* **ci:** update ci action versions ([39117ce](https://github.com/shivanshkc/rosenbridge/commit/39117cefaf66e21fb276e20b16f54c0dd32e42bc))
* **core:** cleanup bridge manager impl ([1f5efba](https://github.com/shivanshkc/rosenbridge/commit/1f5efba67bf04bac023c60a543cf9e0590cecc31))
* **core:** correct nomenclature in comments ([cbf98fe](https://github.com/shivanshkc/rosenbridge/commit/cbf98fedf23b6af98cc4fcca0b9165e3142d3485))
* **core:** discovery addr resol fix ([354a76d](https://github.com/shivanshkc/rosenbridge/commit/354a76d31e062e51e284dbeac3aa95559b9cb49c))
* **core:** fix validation bug ([3c84a6f](https://github.com/shivanshkc/rosenbridge/commit/3c84a6f4873dae984fc289be067b50c56974de6a))
* **core:** import fix ([5e7dfd0](https://github.com/shivanshkc/rosenbridge/commit/5e7dfd01ebb2fae6dcec13260ee4d9e1362f8029))
* **core:** log only the file name, not path elements ([d99e9e9](https://github.com/shivanshkc/rosenbridge/commit/d99e9e9216ebd29494fafa36c520392363fb5d67))
* **core:** project id is fetched after server starts ([3d927c1](https://github.com/shivanshkc/rosenbridge/commit/3d927c19a3db391ed02776c578ab19496554bb16))
* **core:** remove discovery address job and extra ci file ([1a621a3](https://github.com/shivanshkc/rosenbridge/commit/1a621a31438c3188fe0dc563aac2a14ffd267527))
* **core:** response body read fix ([425763f](https://github.com/shivanshkc/rosenbridge/commit/425763f03add5d98d14b603de7ef961e7594e8c1))
* **core:** url formatting bug fix ([33c5bf1](https://github.com/shivanshkc/rosenbridge/commit/33c5bf164e0bce772a7d13c0c27437c8c5e2301e))


### Features

* **ci:** add cd test for cloud run testing ([d09a2bc](https://github.com/shivanshkc/rosenbridge/commit/d09a2bc472feb02a434b243e9e3f7cf181c0c189))
* **ci:** add github actions ([f30652d](https://github.com/shivanshkc/rosenbridge/commit/f30652da3945dffc793466eab2f5c2fa44d8c6bc))
* **ci:** cd-test only deploys ([f0da085](https://github.com/shivanshkc/rosenbridge/commit/f0da08543f4fa10c28684e77bd23bffe53444adb))
* **core:** add access layer for get bridge ([e190992](https://github.com/shivanshkc/rosenbridge/commit/e1909921425f57c10bbfcf004a143c2b3e96017b))
* **core:** add access layer for list bridges ([b54a617](https://github.com/shivanshkc/rosenbridge/commit/b54a6173bdcdab2146ff4ae181b09e24eb500b8c))
* **core:** add access layer for post message ([95f0a06](https://github.com/shivanshkc/rosenbridge/commit/95f0a064baa890297479cb2fc72a84e1d4b974a1))
* **core:** add access layer for post message internal ([1be3bd9](https://github.com/shivanshkc/rosenbridge/commit/1be3bd9a306f21fdfc66a5de127b5e77a719e0f4))
* **core:** add addr fetch logic in gcp resolver ([7ee2dff](https://github.com/shivanshkc/rosenbridge/commit/7ee2dff3de57448d85e2753eaad2f2b8ffdd907f))
* **core:** add all dep definitions ([12520e9](https://github.com/shivanshkc/rosenbridge/commit/12520e933df3e99e095c8997aa8ec762d7692cb3))
* **core:** add bridge and bridge-manager deps ([6aa1502](https://github.com/shivanshkc/rosenbridge/commit/6aa150213643af5ba6ebf987e6be515f78fc357c))
* **core:** add bridge-database dep ([e3b4f77](https://github.com/shivanshkc/rosenbridge/commit/e3b4f777c3be3fbc07a072b70b32cc92b6156c4f))
* **core:** add configs ([b5d2c5e](https://github.com/shivanshkc/rosenbridge/commit/b5d2c5eaf2c07c74f1892e454b170274a00a3260))
* **core:** add create bridge method ([b3e709e](https://github.com/shivanshkc/rosenbridge/commit/b3e709e61f6ed47ecb586282e69fb54b19827c60))
* **core:** add dep skeletons and database connectivity ([edb7f5e](https://github.com/shivanshkc/rosenbridge/commit/edb7f5ededf49fde993cb3cfcec1ae7a85907e75))
* **core:** add disovery address resolver local ([ed774a1](https://github.com/shivanshkc/rosenbridge/commit/ed774a113d8f8e0587a766872efb3f59532fdee8))
* **core:** add handlers and router ([714e70d](https://github.com/shivanshkc/rosenbridge/commit/714e70dbd6bbb62f896405765a787c9eb9e10d8f))
* **core:** add intercom dep ([312ff1e](https://github.com/shivanshkc/rosenbridge/commit/312ff1e9bc8e110dda85d5c4cfbd26efda9bf35e))
* **core:** add intro handler ([b717142](https://github.com/shivanshkc/rosenbridge/commit/b717142fb336e15ab084e95b00f5b01a9846c37f))
* **core:** add list bridges api def ([5cb3d6b](https://github.com/shivanshkc/rosenbridge/commit/5cb3d6bf2ef06f31fec67cb098116fe8012a4b37))
* **core:** add list bridges function ([6812294](https://github.com/shivanshkc/rosenbridge/commit/6812294375e42d5deb0d0700d9d93a6fe39f884f))
* **core:** add list bridges handler def and response headers in bridge create params ([db158fb](https://github.com/shivanshkc/rosenbridge/commit/db158fb8a36df766c4e10554661a49f4dc621728))
* **core:** add logger ([f0bea1b](https://github.com/shivanshkc/rosenbridge/commit/f0bea1b5ec1d5515711af84678930c0a7bc13929))
* **core:** add middleware ([464dba3](https://github.com/shivanshkc/rosenbridge/commit/464dba31dc80d0a98f8e5112a4c6bdef1552655f))
* **core:** add more gcp api ([02a4cdf](https://github.com/shivanshkc/rosenbridge/commit/02a4cdfc2a51a849d67281899d512846b238abc6))
* **core:** add post message func ([a8b824c](https://github.com/shivanshkc/rosenbridge/commit/a8b824c4af453eadd840e24413bc85c28e15461a))
* **core:** add post message internal func ([77cc2ba](https://github.com/shivanshkc/rosenbridge/commit/77cc2ba5863260872394828fd632acf43903e135))
* **core:** add solo mode config option ([ea71e81](https://github.com/shivanshkc/rosenbridge/commit/ea71e813a4a004dcf8946ca8c392552c8e12d8cd))
* **core:** add some validations ([136f8b2](https://github.com/shivanshkc/rosenbridge/commit/136f8b29993565717f948ef85a34e6c5ae5aead1))
* **core:** add struct fields in impls ([06daaa8](https://github.com/shivanshkc/rosenbridge/commit/06daaa86a0b052c695932f7438cc1172e46b3862))
* **core:** add validations ([0d9dd42](https://github.com/shivanshkc/rosenbridge/commit/0d9dd42e0cdaf0fb3ce30a885f814f9226fcbad5))
* **core:** added docker files ([5540092](https://github.com/shivanshkc/rosenbridge/commit/5540092e65be175027c8e23360b2ede563dfb23e))
* **core:** bridge create params do not accept bridge limits ([72f7af3](https://github.com/shivanshkc/rosenbridge/commit/72f7af3e4a0a24d60442285ec15f4f3ce0814b71))
* **core:** discovery address resolution ([1d9421f](https://github.com/shivanshkc/rosenbridge/commit/1d9421ff7ad42f0bd14146917f1aa6e44cd2eae0))
* **core:** include node addr in bridge status ([6a24caa](https://github.com/shivanshkc/rosenbridge/commit/6a24caa9d14034993b5787c3a550f22febbb8b9d))
* **core:** progress on gcp addr resolver ([d07e9f5](https://github.com/shivanshkc/rosenbridge/commit/d07e9f56dd4385eae3153980d21e1bf43fdb13af))
* **core:** progress with post message func ([2252328](https://github.com/shivanshkc/rosenbridge/commit/225232812077bcffe6dbdf0da8f80d3c90d9ed45))
* **core:** remove most boilerplate, progress on core ([275e471](https://github.com/shivanshkc/rosenbridge/commit/275e47171d56a205c8578d82e1010697c9461709))
* **core:** remove separate struct for create bridge params ([9f118e9](https://github.com/shivanshkc/rosenbridge/commit/9f118e97dad39121052c0baa48fdf864328a3d64))
* **core:** restructure core and api ([833e285](https://github.com/shivanshkc/rosenbridge/commit/833e285c1f88f31cbb7a39d1ce5db6841e8429a4))
* **core:** roll back post message internal and update bridge manager interface ([6cd6fea](https://github.com/shivanshkc/rosenbridge/commit/6cd6fea939ba1353a6cdd9ac20caa4c56e35c47e))
