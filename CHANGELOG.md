# CØSMOS Faucet

## [2.0.0](https://github.com/okp4/cosmos-faucet/compare/v1.1.0...v2.0.0) (2023-06-13)


### ⚠ BREAKING CHANGES

* **cosmos:** bump cosmos sdk version to 0.46.13

### Features

* **cosmos:** bump cosmos sdk version to 0.46.13 ([911a649](https://github.com/okp4/cosmos-faucet/commit/911a64934641ee05e346efaf76afad8095c8a793))

# [1.1.0](https://github.com/okp4/cosmos-faucet/compare/v1.0.0...v1.1.0) (2022-08-30)


### Bug Fixes

* **grpc:** properly stop grpc conn when actor stops ([1e4ced0](https://github.com/okp4/cosmos-faucet/commit/1e4ced0bd505ad411804308f421c286e3229329b))
* make it build again ([dbaa637](https://github.com/okp4/cosmos-faucet/commit/dbaa637bf0d90e70f169a279d4e8a78e12df392c))
* make linters happy ([363dc8f](https://github.com/okp4/cosmos-faucet/commit/363dc8ff3ce98b253c3c5f72721325bc348053b4))
* minor imprecisions in tx handler ([835f283](https://github.com/okp4/cosmos-faucet/commit/835f28309d9cd57ca9c77216419be6f8178ae4b6))
* **pool:** avoid sending empty transactions ([c1c4697](https://github.com/okp4/cosmos-faucet/commit/c1c46974fa8b780280990375b2b6bebc882bcb22))
* **pool:** close channels in only one goroutine ([297b0f4](https://github.com/okp4/cosmos-faucet/commit/297b0f4786bd57aaac82fceb8193665d76d15021))


### Features

* add warn log on non 0 tx code ([b377e5f](https://github.com/okp4/cosmos-faucet/commit/b377e5fb5d7c31cf136f2158e77adfba3d7e6dec))
* allow all origins ([d9ef2e8](https://github.com/okp4/cosmos-faucet/commit/d9ef2e81ebc7d127d8aa49e98458ce64656e0f75))
* allow send & subscribe symulnateously ([8fece85](https://github.com/okp4/cosmos-faucet/commit/8fece8550f17464b3de18968ecedb3b216438284))
* allow to set deadline on tx triggers ([bbf6451](https://github.com/okp4/cosmos-faucet/commit/bbf64515f1d3bb9ca5deda393ca93028a9876dd8))
* **api:** implements send graphql subscription ([dad1aae](https://github.com/okp4/cosmos-faucet/commit/dad1aae2970d11acb885b261b344cd92e25823bb))
* **api:** introduce the void scalar ([230de82](https://github.com/okp4/cosmos-faucet/commit/230de820e835a0649db1871e97ad8fd3871821ca))
* **cli:** add temporal batch window flag ([566831c](https://github.com/okp4/cosmos-faucet/commit/566831c26ae728172212a6080c7e91e40f2bf43c))
* **cli:** implement transaction timeout config ([973570c](https://github.com/okp4/cosmos-faucet/commit/973570ca2a220367ce0ca3fb8e02dde84d4392ac))
* implement an atomic message pool ([99634d8](https://github.com/okp4/cosmos-faucet/commit/99634d888f61fd800044a3d769b46b5a8b3bca07))
* implement message pool tx submitter option ([c881f68](https://github.com/okp4/cosmos-faucet/commit/c881f68d71d83f9be1b7b1d10553d0b9f63624d7))
* implements msg batching ([77955fe](https://github.com/okp4/cosmos-faucet/commit/77955fea69098cce82b624ff6ae6b05a20afe6cb))
* **metrics:** disable prometheus middleware ([8504e60](https://github.com/okp4/cosmos-faucet/commit/8504e60e7e25d1719c2a96253490956feed5c4fb))
* setup basic faucet trigger channel ([6ef0d4f](https://github.com/okp4/cosmos-faucet/commit/6ef0d4fd8de1b91e9ce51f526ec0a28e43628453))
* wire the faucet with the message pool ([ea3da76](https://github.com/okp4/cosmos-faucet/commit/ea3da76cf185a0814e19b4599d6c59e1a901b899))

# 1.0.0 (2022-05-31)


### Bug Fixes

* add method on graphql route ([ea5674e](https://github.com/okp4/cosmos-faucet/commit/ea5674e40d853ab68861f240cacbf105f2d1c478))
* change log type of health handler ([c460a98](https://github.com/okp4/cosmos-faucet/commit/c460a98937dd92954c886ec19bcb671229ee3af9))
* effectively set grpc tls parameters ([390170b](https://github.com/okp4/cosmos-faucet/commit/390170b1085219bc1e33cb94990280176165aaac))
* error handling on config file ([594941f](https://github.com/okp4/cosmos-faucet/commit/594941fed7d788b08752aa115f0869c25f1d15e5))
* handle grpc error ([7f35eb5](https://github.com/okp4/cosmos-faucet/commit/7f35eb528feca26a8608504f6ace55cc9bd0fbec))
* linter on github action ([ba32216](https://github.com/okp4/cosmos-faucet/commit/ba322164aaebbbb66d88e95d7b9ca01e4bf9d0db))
* make graphql linter happy ([c951771](https://github.com/okp4/cosmos-faucet/commit/c951771c40fdf50cd2147fe7122972c47c534829))
* metrics now correctly instantiated ([80ae76d](https://github.com/okp4/cosmos-faucet/commit/80ae76da69511d3024b05c98e02b3ed20a4cda0e))
* move server related flags to cmd/start.go ([9757c30](https://github.com/okp4/cosmos-faucet/commit/9757c3087211d22cd0684da9e77725b7d4a11023))
* remove unused send http api handler ([ea3c8d5](https://github.com/okp4/cosmos-faucet/commit/ea3c8d50f4e88fae9da8334da618da2487ffb3b1))
* removed unused param in newSendRequestHandler ([ae4c014](https://github.com/okp4/cosmos-faucet/commit/ae4c014f597782fb31d0ee55108fbcf928afb67c))
* retrive account on each transaction ([b9fbe59](https://github.com/okp4/cosmos-faucet/commit/b9fbe59e43f4ed26364c4e39a62e84bf79f76b75))
* some linter error ([23397b6](https://github.com/okp4/cosmos-faucet/commit/23397b6d17b711ea7071a2a0847721bb6c78e053))
* use request context for make transacation call ([a16d4ec](https://github.com/okp4/cosmos-faucet/commit/a16d4ec6c1b54e12e0d60353b0b65352c661a5c2))


### Features

* add an interface to expose HttpServer ([3810d8e](https://github.com/okp4/cosmos-faucet/commit/3810d8e919e32d52830a858a9afb1c0a62f7fdcd))
* add captcha secret flag to start cmd ([f0aa7fe](https://github.com/okp4/cosmos-faucet/commit/f0aa7fe560553a79f43a5d92137d92f2fa38fbd2))
* add flags for metrics and health endpoint ([1c9688f](https://github.com/okp4/cosmos-faucet/commit/1c9688f1398fa0108ce41d3d1336e58071a18def))
* add health check endpoint and base for prometheus metrics ([bb6e73a](https://github.com/okp4/cosmos-faucet/commit/bb6e73a256c31075cc089e57182c952d3681858c))
* add logs and error handling ([ea9d97e](https://github.com/okp4/cosmos-faucet/commit/ea9d97eb1a9cc0c0c133ff46f67739b7f07cb559))
* add query to fetch configuration of faucet ([75e7145](https://github.com/okp4/cosmos-faucet/commit/75e71451cee994e9baf9bf1f18aece8d861b739a))
* add send command ([d8d1346](https://github.com/okp4/cosmos-faucet/commit/d8d134689abcfcddb9f36ef117020c44d869b29a))
* add some logs ([f5fe13f](https://github.com/okp4/cosmos-faucet/commit/f5fe13f69f19d1f174eb1864e7de146de9052662))
* add the start cmd ([827afce](https://github.com/okp4/cosmos-faucet/commit/827afce72ff8b876a991a86ed9cc8a226498f9c3))
* add UInt64 scalar ([e8d160c](https://github.com/okp4/cosmos-faucet/commit/e8d160c8d4da30e6c0d4b49c7edb3c97e0a98536))
* allow changing rest server port address ([eb72871](https://github.com/okp4/cosmos-faucet/commit/eb72871fea6cd26db3497798bd97ebf9677ffc4e))
* build the config file ([3a18495](https://github.com/okp4/cosmos-faucet/commit/3a184956067b59b07bf8181d1138a4dd3c1ddd5e))
* captcha check in graphql ([aacea5c](https://github.com/okp4/cosmos-faucet/commit/aacea5c2761c64c48e7705198b9d810b330af4da))
* captcha is now optional ([61228b1](https://github.com/okp4/cosmos-faucet/commit/61228b1217c2ac809053f8393fbd271311b86153))
* captcha middleware first impl ([369a45e](https://github.com/okp4/cosmos-faucet/commit/369a45e568fa9b14ccf74018b0576aa5f1d484ab))
* captcha middleware now use header field ([db8ea4f](https://github.com/okp4/cosmos-faucet/commit/db8ea4f203d8afb41aa9ec78d73ec2b55b65a382))
* changes on captcha behavior ([e66f92b](https://github.com/okp4/cosmos-faucet/commit/e66f92b5e3bc3dd488ed2c47b17a0c726aa98723))
* generate addres and keys from menmonic ([bd33a58](https://github.com/okp4/cosmos-faucet/commit/bd33a58f49e311a82378dc86b723b8f97ed87e71))
* handle health check response error if write fails ([6f92e73](https://github.com/okp4/cosmos-faucet/commit/6f92e73a3fade89955ed79dea9b9dab254836586))
* handle prometheus initialization errors ([8b43886](https://github.com/okp4/cosmos-faucet/commit/8b43886b48305364505b84d971e7311fb31026f0))
* handle rest error ([34afb2e](https://github.com/okp4/cosmos-faucet/commit/34afb2e2873a2663e3361ab7bef6b2b3d02887b0))
* handle send rest endpoint ([7ddd14c](https://github.com/okp4/cosmos-faucet/commit/7ddd14c4de34e30d879e43a43856d6f13a184556))
* import config from file or env ([beafb9e](https://github.com/okp4/cosmos-faucet/commit/beafb9e4d834e124d0a31f3e534a8e94b7c1db22))
* make address a graphql scalar ([a4c4d20](https://github.com/okp4/cosmos-faucet/commit/a4c4d20810aa8bc669376216370deddfc78cb019))
* move server related code to server package ([10d1f78](https://github.com/okp4/cosmos-faucet/commit/10d1f78531fa6b7d1bf12d7e9a27891c6c7cfb66))
* prometheus metrics ([c6e3f84](https://github.com/okp4/cosmos-faucet/commit/c6e3f84be62be920611acced54f3898be50cdabd))
* reimplemntation of metrics ([cb2ee0c](https://github.com/okp4/cosmos-faucet/commit/cb2ee0cd1a0b90648797fb1f706fb52e41d8bff3))
* send mutation return result of transaction ([050d32c](https://github.com/okp4/cosmos-faucet/commit/050d32c80c298f1ad3788c355ac4838e458eeb6c))
* send mutation with graphql ([5de61b6](https://github.com/okp4/cosmos-faucet/commit/5de61b60406810ae464a4bb183a9f3921d527b0d))
* sign and brodcast transaction ([3e590fa](https://github.com/okp4/cosmos-faucet/commit/3e590fa2cd595018e27b4d1d55cf3da0410e2b3b))
* tls by default and allow disable it ([043d450](https://github.com/okp4/cosmos-faucet/commit/043d4509deb7808b5febd3a68f8583467da60439))
* unexport server functions used only in server pkg ([f0c0a0e](https://github.com/okp4/cosmos-faucet/commit/f0c0a0e73c5cbce0caaaac25eda7a3229416f599))
