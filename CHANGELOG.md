# [0.13.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.12.0...v0.13.0) (2023-07-31)


### Bug Fixes

* **blueprint:** fix content and module guard, allow variables in interpolation ([21ceb5d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/21ceb5df37e9b6f8d54253aec26dcf6b60970673))
* **ep-3291:** fixing interpolation vars in generic variables ([ae9d325](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/ae9d325ffbdb6e4134cf611749ba13ae0727939c))
* **generic_variable:** type in generic/global variable should always have a type ([1445dd8](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1445dd804840598f5e9f7df1c2a7ad7ef918dddd))
* **iac:** `display_order` generates error if values is empty ([831a573](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/831a573dab50b25719932b29743487cba01c6d34))
* **iac:** change error message when an invalid "config" value is passed in ([0ffe28d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/0ffe28d864ac7087580853428cbe7f0c823f4a39))
* **iac:** conditionals are not working as expected ([681ab31](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/681ab31b4b8221c6832b566d5bcb8c15f433cac5))


### Features

* **autocloud_blueprint_config:** support terraform variable default values ([53d5154](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/53d515472142cd4f96e338a2624802594445537f))
* **blueprint_config:** create generic variable for list types ([652e48c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/652e48cc9257aab258e396acca3d9e2762c8afce))
* **generic:** add generic map with required values and fix generic shortText with value ([5ba4387](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/5ba43874ed977c4d74b787e1cb9d435050bc1112))
* **iac:** add minimal validations ([1af2e6c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1af2e6c2fc3d8c102aaadb78c3de9e1fcea0b54b))
* **iac:** support variable interpol in conditionals ([c6d679d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/c6d679db86f86ed9d78a304430356e0470fce6aa))
* **variable_interpolation:** detect when an unused variable in a template for interpolation ([98224e0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/98224e02cdb37027e29db94903963acb13c1411e))

# [0.13.0-pre.9](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.13.0-pre.8...v0.13.0-pre.9) (2023-07-26)


### Bug Fixes

* **ep-3291:** fixing interpolation vars in generic variables ([ae9d325](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/ae9d325ffbdb6e4134cf611749ba13ae0727939c))

# [0.13.0-pre.8](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.13.0-pre.7...v0.13.0-pre.8) (2023-07-21)


### Bug Fixes

* **iac:** `display_order` generates error if values is empty ([831a573](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/831a573dab50b25719932b29743487cba01c6d34))

# [0.13.0-pre.7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.13.0-pre.6...v0.13.0-pre.7) (2023-07-19)


### Bug Fixes

* **iac:** conditionals are not working as expected ([681ab31](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/681ab31b4b8221c6832b566d5bcb8c15f433cac5))

# [0.13.0-pre.6](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.13.0-pre.5...v0.13.0-pre.6) (2023-07-11)


### Bug Fixes

* **blueprint:** fix content and module guard, allow variables in interpolation ([21ceb5d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/21ceb5df37e9b6f8d54253aec26dcf6b60970673))

# [0.13.0-pre.5](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.13.0-pre.4...v0.13.0-pre.5) (2023-07-08)


### Features

* **generic:** add generic map with required values and fix generic shortText with value ([5ba4387](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/5ba43874ed977c4d74b787e1cb9d435050bc1112))

# [0.13.0-pre.4](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.13.0-pre.3...v0.13.0-pre.4) (2023-07-07)


### Features

* **variable_interpolation:** detect when an unused variable in a template for interpolation ([98224e0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/98224e02cdb37027e29db94903963acb13c1411e))

# [0.13.0-pre.3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.13.0-pre.2...v0.13.0-pre.3) (2023-06-30)


### Features

* **blueprint_config:** create generic variable for list types ([652e48c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/652e48cc9257aab258e396acca3d9e2762c8afce))

# [0.13.0-pre.2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.13.0-pre.1...v0.13.0-pre.2) (2023-06-29)


### Bug Fixes

* **iac:** change error message when an invalid "config" value is passed in ([0ffe28d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/0ffe28d864ac7087580853428cbe7f0c823f4a39))


### Features

* **autocloud_blueprint_config:** support terraform variable default values ([53d5154](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/53d515472142cd4f96e338a2624802594445537f))

# [0.13.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.12.0...v0.13.0-pre.1) (2023-06-23)


### Bug Fixes

* **generic_variable:** type in generic/global variable should always have a type ([1445dd8](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1445dd804840598f5e9f7df1c2a7ad7ef918dddd))


### Features

* **iac:** add minimal validations ([1af2e6c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1af2e6c2fc3d8c102aaadb78c3de9e1fcea0b54b))
* **iac:** support variable interpol in conditionals ([c6d679d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/c6d679db86f86ed9d78a304430356e0470fce6aa))

# [0.12.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.11.0...v0.12.0) (2023-06-21)


### Bug Fixes

* **tree:** use aliases hash to map parent and child blueprintconfigs ([627835e](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/627835e2c2a9b466e2e0730e40fef3981ca9b2b1))


### Features

* **config:** set prod url as default endpoint ([a77da18](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/a77da18d6fa651e7c4cd75dc2b1f6688987809f1))
* **iac:** support for generatedName and external variables from blueprint_config ([266e688](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/266e68822cb9ad53fe990d0321ab71207cf33b8d))
* **validations:** add scope to regex ([9850dd1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9850dd1c116f2ecb2f8137be6cbad8d2b228e207))

# [0.12.0-pre.4](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.12.0-pre.3...v0.12.0-pre.4) (2023-06-21)


### Features

* **validations:** add scope to regex ([9850dd1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9850dd1c116f2ecb2f8137be6cbad8d2b228e207))

# [0.12.0-pre.3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.12.0-pre.2...v0.12.0-pre.3) (2023-06-20)


### Features

* **iac:** support for generatedName and external variables from blueprint_config ([266e688](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/266e68822cb9ad53fe990d0321ab71207cf33b8d))

# [0.12.0-pre.2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.12.0-pre.1...v0.12.0-pre.2) (2023-06-14)


### Bug Fixes

* **tree:** use aliases hash to map parent and child blueprintconfigs ([627835e](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/627835e2c2a9b466e2e0730e40fef3981ca9b2b1))

# [0.12.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.11.0...v0.12.0-pre.1) (2023-06-12)


### Features

* **config:** set prod url as default endpoint ([a77da18](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/a77da18d6fa651e7c4cd75dc2b1f6688987809f1))

# [0.11.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.10.0...v0.11.0) (2023-06-09)


### Bug Fixes

* **iac:** apply generic/global variables in interpolation too ([63eec39](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/63eec39790de1c9f506452aba29eaff1b40e34b7))
* **iac:** remove formShapeMap legacy feature ([4d396b1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/4d396b18828b53d9ae30f34bf2351f712b8d008a))


### Features

* **iac:** support global variable values as reference to get other vars value ([9f7121b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9f7121b59f912a2d86adbf5335889b888b6ec625))

# [0.11.0-pre.3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.11.0-pre.2...v0.11.0-pre.3) (2023-06-09)


### Bug Fixes

* **iac:** apply generic/global variables in interpolation too ([63eec39](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/63eec39790de1c9f506452aba29eaff1b40e34b7))

# [0.11.0-pre.2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.11.0-pre.1...v0.11.0-pre.2) (2023-06-09)


### Bug Fixes

* **iac:** remove formShapeMap legacy feature ([4d396b1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/4d396b18828b53d9ae30f34bf2351f712b8d008a))

# [0.11.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.10.0...v0.11.0-pre.1) (2023-06-08)


### Features

* **iac:** support global variable values as reference to get other vars value ([9f7121b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9f7121b59f912a2d86adbf5335889b888b6ec625))

# [0.10.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.9.2...v0.10.0) (2023-06-07)


### Features

* **blueprint:** add variable interpolation ([67eaa2c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/67eaa2c39aaf5cf8ba794973b5ec400d47e0afab))

# [0.10.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.9.2...v0.10.0-pre.1) (2023-06-07)


### Features

* **blueprint:** add variable interpolation ([67eaa2c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/67eaa2c39aaf5cf8ba794973b5ec400d47e0afab))

## [0.9.2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.9.1...v0.9.2) (2023-05-31)


### Bug Fixes

* Set variables interpolation block as optional ([124f013](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/124f013c0c5c588b0497601f42540307b7489e05))

## [0.9.2-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.9.1...v0.9.2-pre.1) (2023-05-31)


### Bug Fixes

* Set variables interpolation block as optional ([124f013](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/124f013c0c5c588b0497601f42540307b7489e05))

## [0.9.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.9.0...v0.9.1) (2023-05-16)


### Bug Fixes

* Fixed display_order with references ([8198ff7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/8198ff71bc5b75f5e8e2712d5b677a5c77139747))

## [0.9.1-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.9.0...v0.9.1-pre.1) (2023-05-16)


### Bug Fixes

* Fixed display_order with references ([8198ff7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/8198ff71bc5b75f5e8e2712d5b677a5c77139747))

# [0.9.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.8.0...v0.9.0) (2023-05-16)


### Features

* Enabled editor type for overrides ([0525940](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/052594099e0a4c8fce7091ac7877d73913fa5e9d))
* Overrided data type as editor ([dd1337e](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/dd1337eaa3fb83e04e1ef8b0bea1b761c02d2885))

# [0.9.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.8.0...v0.9.0-pre.1) (2023-05-16)


### Features

* Enabled editor type for overrides ([0525940](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/052594099e0a4c8fce7091ac7877d73913fa5e9d))
* Overrided data type as editor ([dd1337e](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/dd1337eaa3fb83e04e1ef8b0bea1b761c02d2885))

# [0.8.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.7.0...v0.8.0) (2023-05-05)


### Bug Fixes

* Fixed error on update file definitions ([712b384](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/712b3844f0c9c12a86a87939f3c6a775c4ea19ed))
* Fixed error on update file definitions ([be31dc0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/be31dc010bb625bdf2f8030913ba5767d030ed91))
* Returned errors when loading existing state ([ebf6a86](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/ebf6a8698ac10689a156882b8cb5b3abf4d97ed1))
* Warning on not compatibile schema ([67d9081](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/67d90812be87f906e9fc0b735bba7134040658ea))


### Features

* Retrieved state from TFC ([16f0313](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/16f0313c824c616c2be6a884b3209ea2343dc55c))

# [0.8.0-pre.2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.8.0-pre.1...v0.8.0-pre.2) (2023-05-04)


### Bug Fixes

* Warning on not compatibile schema ([67d9081](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/67d90812be87f906e9fc0b735bba7134040658ea))

# [0.8.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.7.1-pre.1...v0.8.0-pre.1) (2023-05-04)


### Bug Fixes

* Fixed error on update file definitions ([712b384](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/712b3844f0c9c12a86a87939f3c6a775c4ea19ed))
* Returned errors when loading existing state ([ebf6a86](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/ebf6a8698ac10689a156882b8cb5b3abf4d97ed1))


### Features

* Retrieved state from TFC ([16f0313](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/16f0313c824c616c2be6a884b3209ea2343dc55c))

## [0.7.1-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.7.0...v0.7.1-pre.1) (2023-05-03)


### Bug Fixes

* Fixed error on update file definitions ([be31dc0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/be31dc010bb625bdf2f8030913ba5767d030ed91))

# [0.7.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.6.0...v0.7.0) (2023-05-03)


### Features

* **iac:** add ability to add raw header and footer iac module ([51d9ffc](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/51d9ffca6696fcd234b9f6ec8c9705d1adde981d))

# [0.7.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.6.0...v0.7.0-pre.1) (2023-05-03)


### Features

* **iac:** add ability to add raw header and footer iac module ([51d9ffc](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/51d9ffca6696fcd234b9f6ec8c9705d1adde981d))

# [0.6.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0...v0.6.0) (2023-04-19)


### Bug Fixes

* Displayed an error when autocloud token is not set ([131dfe0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/131dfe0ab749d531ff98b279a54c07cd13acba4e))
* Removed autogenerated item when there are not options for list ([30ca4ff](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/30ca4ffe69ec7bdc2c1edb8145b943d72b734c5b))
* Updated lists example ([9789dc1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9789dc1555d5a62f17c704371bd6ef1e668a23d7))


### Features

* **displayOrder:** warning on duplicate priority ([6cf140f](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/6cf140fca44820f017a7fc613cad5c14f1b4d4b8))

# [0.6.0-pre.2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.6.0-pre.1...v0.6.0-pre.2) (2023-04-18)


### Bug Fixes

* Displayed an error when autocloud token is not set ([131dfe0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/131dfe0ab749d531ff98b279a54c07cd13acba4e))

# [0.6.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0...v0.6.0-pre.1) (2023-04-18)


### Bug Fixes

* Removed autogenerated item when there are not options for list ([30ca4ff](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/30ca4ffe69ec7bdc2c1edb8145b943d72b734c5b))
* Updated lists example ([9789dc1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9789dc1555d5a62f17c704371bd6ef1e668a23d7))


### Features

* **displayOrder:** warning on duplicate priority ([6cf140f](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/6cf140fca44820f017a7fc613cad5c14f1b4d4b8))

# [0.5.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.4.0...v0.5.0) (2023-04-12)


### Bug Fixes

* Extended file_content example using omit dot notation ([8a4bdca](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/8a4bdca9d5e637f1351213efb85cb33dd3d3b6ab))
* Fixed omit specific variable with dot notation ([b7e5d58](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/b7e5d5869c089cbabc57d350e9aaa0db4ea9faf6))
* Left map convertion logic to API ([570ff2e](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/570ff2e021fa81873945122c2e2c5ffe76572066))
* use correct ARCH for darwin_arm64 ([53b8013](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/53b80132c3a84dd26ed654ac6e3a680c57766e3a))


### Features

* add display order ([f376a10](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f376a10fc01f8b4d9869bad4cfa1c4b9ea76209d))
* Added example with header and footer ([32e29fc](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/32e29fc7d43c8e2cc18352387cbc8b646c4f1dab))
* Added file block missing attributes validation ([d349dd3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/d349dd32835dcbd56f6aa6dc5d50b65d13e32032))
* Added file_content example ([49c9206](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/49c9206bf930b230774551dce5b0d44b59df4598))
* Added raw type ([7362e66](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/7362e66e2b6ade2015506b4a50c5014d5ab37104))
* Added unit tests for omitted by reference variable ([5016fff](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/5016fffe90d6a4a832bb6bdbcb8c13700d825b51))
* Added validation for non empty modules using footer/header attributes ([152f525](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/152f525488887800bcd0f59dec772867b0fc7f5b))
* Added warning when using content property on fileblock ([afd038f](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/afd038fd559945e9744023217847d659c3390f1a))
* Extended File block with header and footer attributes ([1ca9009](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1ca900977225741f97b09265248472188c8d6f39))
* Extended filecontent example ([0eb3090](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/0eb30900a67c7235a305c1a403fc2028cb9a28c9))
* Extended fileDefinition schema ([13f7da9](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/13f7da9135966f48187141e6d9de339997588227))
* Extended list-basic example with raw datatype ([59a78af](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/59a78af2b8d89c468fc56284ecefc8e94b4c1aab))
* Included new validation rules for maxLength and minLength ([0dcc3a1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/0dcc3a1a486c2480f9a4604e205111894c4ef923))
* integrate new git repository service ([15a45b4](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/15a45b47867eb8ab48d5dbad909f98339d3f443a))
* **module:** do not force tags name on modules ([fc47175](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/fc4717570411b4cbc39a53f9ea3ddaef48dd0252))
* Omitted just reference variable ([9da1aea](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9da1aea4206df0dd71ec997898efc6dd728d0ca3))
* **sdk:** update tag ([5a4aa68](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/5a4aa68fc4ac4d734f5d04e46c6c1452b573377d))
* Set raw type as hcl-expression ([a2526b1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/a2526b134bb7814b9f76290dbf44a0e7ea39a689))
* Set raw variableType ([2d2ddf7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/2d2ddf77e3865ad963344c44a54ae34f10049666))
* Updated list-basic example ([6e12daf](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/6e12dafd8577ac9163e4c22813c8a1717fe6af43))

# [0.5.0-pre.10](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0-pre.9...v0.5.0-pre.10) (2023-04-12)


### Bug Fixes

* Extended file_content example using omit dot notation ([8a4bdca](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/8a4bdca9d5e637f1351213efb85cb33dd3d3b6ab))
* Fixed omit specific variable with dot notation ([b7e5d58](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/b7e5d5869c089cbabc57d350e9aaa0db4ea9faf6))


### Features

* add display order ([f376a10](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f376a10fc01f8b4d9869bad4cfa1c4b9ea76209d))

# [0.5.0-pre.9](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0-pre.8...v0.5.0-pre.9) (2023-04-12)


### Features

* Added warning when using content property on fileblock ([afd038f](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/afd038fd559945e9744023217847d659c3390f1a))
* Extended filecontent example ([0eb3090](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/0eb30900a67c7235a305c1a403fc2028cb9a28c9))

# [0.5.0-pre.8](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0-pre.7...v0.5.0-pre.8) (2023-04-10)


### Features

* Added example with header and footer ([32e29fc](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/32e29fc7d43c8e2cc18352387cbc8b646c4f1dab))
* Added validation for non empty modules using footer/header attributes ([152f525](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/152f525488887800bcd0f59dec772867b0fc7f5b))
* Extended File block with header and footer attributes ([1ca9009](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1ca900977225741f97b09265248472188c8d6f39))

# [0.5.0-pre.7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0-pre.6...v0.5.0-pre.7) (2023-04-06)


### Features

* Added file block missing attributes validation ([d349dd3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/d349dd32835dcbd56f6aa6dc5d50b65d13e32032))
* Added file_content example ([49c9206](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/49c9206bf930b230774551dce5b0d44b59df4598))
* Extended fileDefinition schema ([13f7da9](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/13f7da9135966f48187141e6d9de339997588227))

# [0.5.0-pre.6](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0-pre.5...v0.5.0-pre.6) (2023-04-04)


### Features

* Added raw type ([7362e66](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/7362e66e2b6ade2015506b4a50c5014d5ab37104))
* Extended list-basic example with raw datatype ([59a78af](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/59a78af2b8d89c468fc56284ecefc8e94b4c1aab))
* Set raw type as hcl-expression ([a2526b1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/a2526b134bb7814b9f76290dbf44a0e7ea39a689))
* Set raw variableType ([2d2ddf7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/2d2ddf77e3865ad963344c44a54ae34f10049666))

# [0.5.0-pre.5](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0-pre.4...v0.5.0-pre.5) (2023-04-03)


### Features

* **module:** do not force tags name on modules ([fc47175](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/fc4717570411b4cbc39a53f9ea3ddaef48dd0252))
* **sdk:** update tag ([5a4aa68](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/5a4aa68fc4ac4d734f5d04e46c6c1452b573377d))

# [0.5.0-pre.4](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0-pre.3...v0.5.0-pre.4) (2023-03-31)


### Features

* Added unit tests for omitted by reference variable ([5016fff](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/5016fffe90d6a4a832bb6bdbcb8c13700d825b51))
* Included new validation rules for maxLength and minLength ([0dcc3a1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/0dcc3a1a486c2480f9a4604e205111894c4ef923))
* Omitted just reference variable ([9da1aea](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9da1aea4206df0dd71ec997898efc6dd728d0ca3))
* Updated list-basic example ([6e12daf](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/6e12dafd8577ac9163e4c22813c8a1717fe6af43))

# [0.5.0-pre.3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0-pre.2...v0.5.0-pre.3) (2023-03-27)


### Bug Fixes

* use correct ARCH for darwin_arm64 ([53b8013](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/53b80132c3a84dd26ed654ac6e3a680c57766e3a))

# [0.5.0-pre.2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.5.0-pre.1...v0.5.0-pre.2) (2023-03-22)


### Bug Fixes

* Left map convertion logic to API ([570ff2e](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/570ff2e021fa81873945122c2e2c5ffe76572066))

# [0.5.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.4.0...v0.5.0-pre.1) (2023-03-21)


### Features

* integrate new git repository service ([15a45b4](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/15a45b47867eb8ab48d5dbad909f98339d3f443a))

# [0.4.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.3.0...v0.4.0) (2023-03-21)


### Bug Fixes

* Exposed config for autocloud_module ([4ec8c5d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/4ec8c5d72fc4eddb215d331f1726e1761f5ea9ae))
* **iac:** adapt repository pattern to the module resource ([bf31b68](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/bf31b688ee6b708ffe704eca12772dc9796e464e))


### Features

* **iac:** add ability to fetch local modules ([d8682c7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/d8682c762870b1db4f3f78311039f2ca4e612c52))

# [0.4.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.3.1-pre.2...v0.4.0-pre.1) (2023-03-20)


### Features

* **iac:** add ability to fetch local modules ([d8682c7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/d8682c762870b1db4f3f78311039f2ca4e612c52))

## [0.3.1-pre.2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.3.1-pre.1...v0.3.1-pre.2) (2023-03-17)


### Bug Fixes

* **iac:** adapt repository pattern to the module resource ([bf31b68](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/bf31b688ee6b708ffe704eca12772dc9796e464e))

## [0.3.1-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.3.0...v0.3.1-pre.1) (2023-03-17)


### Bug Fixes

* Exposed config for autocloud_module ([4ec8c5d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/4ec8c5d72fc4eddb215d331f1726e1761f5ea9ae))

# [0.3.0](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/v0.2.0...v0.3.0) (2023-03-09)


### Bug Fixes

* Added better response for http error fetching repositories ([6609a73](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/6609a73eb85bcb86407fd6c7e19f04f31780b144))
* Added default value for tags field ([ab3c5a7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/ab3c5a7389a9bd6d89e69522349f13f421a7309b))
* adding AllowConsumerToEdit to true when overriding a variable ([487845f](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/487845fc20c5e75e48918ebe9e3d56f441dffddd))
* **conditionals:** use the same content shape as variable for conditionals ([d3072c3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/d3072c3a1e1650236450d3fb10b9d2ee20da6d56))
* Default options from module for override variables ([257124a](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/257124a3895134c62d09968a65aef8afe13904ae))
* Does not rewrite fielddatatype ([b5df9c3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/b5df9c3bf17cc543112dfb2c65e84ad12409d308))
* **ep-2499:** fixing tests ([501d21a](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/501d21aebf8178c2c588592aa7a90bd18fb2c840))
* **ep-2768:** iac file block vars ([73e63e5](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/73e63e545d021da68fa052221d891c90e2db0e35))
* **ep-2773:** iac shared overridden variables ([10ab06a](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/10ab06aa1a7bcfaaac52aaa3dd05d69c5f04646a))
* **examples:** add stub basic example ([08fd505](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/08fd505ac489bbcf3c6971c516dc010aa09d2dbb))
* Fixed issues with field types ([aec64c4](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/aec64c486071c343f32ff516611cc3ff19da8127))
* Fixed linter issues ([ba0772b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/ba0772b49635c74a324d664833e760b037b1cb78))
* Fixed missing types for overrides variables ([50bdd4d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/50bdd4d44fe0dbd9dc2fa7d69dfc69089ec96095))
* Fixed source_blueprint example ([cf22f35](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/cf22f35030bfe942946500c1177fde142583ee1f))
* **hcl:** when a user overrides a variable, should be always used in hcl ([dae126c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/dae126c22fb6432864e605e1aa1945119cb902f1))
* linting issues ([e879caf](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e879cafb0eff3cbd98e498de10efc728263df866))
* **makefile:** comment version ([e87f68c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e87f68cb2002d0e050bbadf664601bd4f2657139))
* **omits:** double omits are not impacting usedHCL if is overriden ([9b4a8ef](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9b4a8efa95f3debf01b1867a9a31de8dd1e217f7))
* Renabled omit and override variables for new blueprint_config ([55c1664](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/55c16644cc31ab0f400f29caac5fac82ff60bfae))
* Renamed terraform_processor to blueprint_config ([8905154](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/8905154165f4110f3019092c4001f2345ea89269))
* **sdk:** get develop version of sdk ([55244ee](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/55244ee57e76a97b5fc34c5c8e14065665ffe781))
* Turned off linter for gocyclo because of GetBlueprintConfigFromSchema function complexity ([e124be7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e124be787a516740a18f7e98c9befd951d0069ee))
* Updated override_var example ([22cac2c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/22cac2cdb558e0af5be363d480cca3d954db70ac))


### Features

* Added display_order field to module schema ([e1cf88c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e1cf88c6eab6ebc25f3ee16e824fa45606307ef1))
* Added examples with simple objects fields ([f733591](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f733591eccb6e261247d9a8526758035b503bc69))
* Added list-objects example ([3417e00](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/3417e00e8ea7496b11c2062cb0ec876123e3b097))
* Added map type to provider ([4b58396](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/4b5839634b1b2cafa59796a7583d38e8caf90f1a))
* Added maps example ([2b169e3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/2b169e3db8d8c8adc890c86accadb0d192826d9f))
* Added reference from blueprint variables example ([fe03d65](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/fe03d65d3f856c1b230206773c3a54fc7f2fc2d0))
* Added tags_variables field to module ([f2c6d6b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f2c6d6b1a7a36da2ea1e3c905c3d9a4d28c85500))
* **blueprint:** add context variable to blueprint resource ([1897d52](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1897d52efe8ed78afbb11e1b93a95bab3e98dc53))
* **blueprint:** add tf example ([7be7c32](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/7be7c3202e2c268dc51dbe1ed64e6e6977a1ab4a))
* **blueprint:** add tf source example ([0809365](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/0809365dd68511a3f1235014a162d238a0c07266))
* **blueprint:** get form shape from tree ([3a2136b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/3a2136b854e5e049b717de4839666032f7942f57))
* **blueprint:** use "source" keyword ([34051f5](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/34051f5f96571b88833f0be93b8f86e63929a542))
* Built generic form using custom variables ([9491b9b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9491b9b23b248dbbc004a82e375b87e6d3e5deb2))
* **conditionals:** integrate conditionals with tree ([1dc8f08](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1dc8f08d8e5877882b393353b7f7d2fa7b1ef209))
* Defined list type for default cases ([9c87d25](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9c87d250a1dde6ec887b015aa311683a067265e6))
* Defined TODO for refactors ([654942e](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/654942ea541e2e534d824c3d2c42e3a1a3137eb3))
* **eks_demo:** provisional eks demo ([6c185a9](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/6c185a9aad91c71c676d241669c20992e614bd02))
* **EP-2466:** modules variables id output ([82f023b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/82f023b77e772535ed1f00863381a2e8232fefb5))
* **EP-2499:** iac blueprint module override ([2c998a3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/2c998a342145efb9ac564206413576d531167de5))
* **ep-2502:** form builder ([d0c7100](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/d0c71005bba0a007155ab45dbbeb218e300d3b71))
* **EP-2503:** IAC token auth ([7058241](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/7058241d6758c66c88fb24f481a07d12995b4654))
* **ep-2559:** blueprint: limit git_config to 1 block ([73a282d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/73a282deaa5da328d6c0d6c7287caa8ff24bd425))
* **ep-2562:** override variable: display name and helper text ([3943b28](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/3943b28b6a4e98acd43b1b78ee0ddba2e5a4772c))
* **ep-2575:** IAC Module composability - updating tests ([64ee6a8](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/64ee6a84a8d330d39ef1206887ab45d933dac004))
* **ep-2575:** IAC Module composability through the outputs of another module ([acb071d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/acb071d43c5b4d91392e1ec6389dddadf456e8f7))
* **ep-2577:** autocloud provider endpoint ([84ccd11](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/84ccd111f8cc0a306c2ef92068ef7b80e653a425))
* **ep-2578:** blueprint file block update ([9e20850](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9e20850031b477c570f79e7b1413981396f4c004))
* **ep-2588:** renaming autocloud_module form_config into blueprint_config ([83bb2fe](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/83bb2fe22de37909061ba035cacdef4f8f2ac7b6))
* **ep-2617:** iac conditionals ([af2d03d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/af2d03dbe822e1b667be1c8b374583c727e48afe))
* **ep-2617:** iac conditionals ([7ce70e3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/7ce70e36b5a0208b1ae29efcd72b93ffb7531c1c))
* **ep-2617:** iac conditionals (fmt to log) ([03fc304](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/03fc30421f2187130ec83ed62eb60db98a6514fe))
* **ep-2628:** iac lists ([979bfbb](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/979bfbb94d6a5bd5fdba25187b3aa5852c7eb54d))
* **ep-2685:** iac conditionals lists and required list values ([2a9cd10](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/2a9cd1011f642edba0b2d6956adedd9fb7e8581e))
* **ep-2769:** iac variable.value ([b512139](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/b512139fd61226a36b1bf6d0d9cae50ef6dae73d))
* **ep-2770:** iac conditionals content value ([fdaef35](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/fdaef3564f000b9fb3d4360628c8c8ed6df32d0f))
* **ep-2838:** iac override value ([c5d1b90](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/c5d1b900be4dce375b61c1455bf836b61445d917))
* **ep-2845:** iac blueprint module ids ([58bcf01](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/58bcf01bfa0ad4d024fc0a61a6f43db365cfeff9))
* **ep-2845:** iac blueprint module ids ([1ce8f34](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1ce8f341effc437bf01b3ed5478d54db06f1965e))
* **ep-2914:** iac list of objects ([52eaf11](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/52eaf116183ed534b317f247eaf54db334edbb47))
* Exposed variables from blueprint_config ([94bb92c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/94bb92c1d935ef703687ecff226f3bdc3c41e643))
* Extended list-conditionals example ([9416c59](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9416c59450d24752f554c3490a9b3e051a2312bb))
* **gitconfig:** avoid configuration errors with git default repo ([b5c3c9f](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/b5c3c9fc335ac9b5b4395b83b14b26645a0ec5f3))
* **gitconfig:** do not force github configuration ([b8542bb](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/b8542bbea114e71ef704a1203762644d501ea69a))
* **iac:** introduce lowlevel blueprintconfig ([3782d56](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/3782d56d0136a96cd93d6c6357b55e247e335e64))
* **omits:** process omitted variables ([f4c3f23](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f4c3f23e7e4d429fa30fb5db534771f156008226))
* **omits:** process omitted variables ([ba181fb](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/ba181fbe0f7490945cd0cc82d8bb436ad810af76))
* **plugin:** Add flag to gate the new blueprintconfig ([34960e3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/34960e3fff554baf4821db1df0097f24b1568690))
* **plugin:** EXPERIMENTAL first approach to support low level lib ([15513bc](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/15513bc789e3f74af7ea8da0d321aa31ebd969be))
* Renamed fields_options to options ([63aa877](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/63aa877bcfee416441aef5c6ed854af98bd259a2))
* Replaced override_variable block with variable ([f6789a6](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f6789a64a6283cbd7850b925af03041b55a0f88e))
* **schema:** remove form options ([e50f418](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e50f418e69bb76ec1109f3430058045992cc44ed))
* **schema:** remove form options ([bb04707](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/bb047072b7e4937c2359247ef36a7d94fe61aeb2))
* set variables map ([c3d2a65](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/c3d2a651efdd2cd7318bd5e2f5f0aff330bc1545))
* Supported generic and module blueprint_config ([8e69c5c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/8e69c5c1051d720b4ad6311a7e8ba6ff9470d6e9))
* **tree:** parse json to include overridevariables ([57b6332](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/57b6332ae5a8b356959a80c57ee04ce654a16766))
* **tree:** process variables by tree level ([28fd1b2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/28fd1b260bacc85eb10bd3f49ac7cc1d3e7cf1d8))
* Updated basic list example ([f775914](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f775914be41bc2a07428c915782d3aad7b64b84f))
* Updated reference example ([3132e94](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/3132e940e64d90e0d6bea60d8951460c33aa1904))
* **validation:** add validation to autocloud_module name ([4f3a18a](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/4f3a18addf03f6499060ffea044dafcaf4ed341d))
* **variable_override:** modify children's variables by reference ([a1a40aa](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/a1a40aaf61e8522befc78bcecbe7491859c9c23f))
* **variable_override:** set default value and add regex to detect ouputs ([e7f9253](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e7f925305c519e9495c3e357a717aab11946c892))

# [0.3.0-pre.54](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.53...0.3.0-pre.54) (2023-03-08)


### Features

* **gitconfig:** do not force github configuration ([b8542bb](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/b8542bbea114e71ef704a1203762644d501ea69a))

# [0.3.0-pre.51](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.50...0.3.0-pre.51) (2023-02-17)


### Features

* **iac:** introduce lowlevel blueprintconfig ([3782d56](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/3782d56d0136a96cd93d6c6357b55e247e335e64))

# [0.3.0-pre.50](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.49...0.3.0-pre.50) (2023-02-17)


### Features

* Added reference from blueprint variables example ([fe03d65](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/fe03d65d3f856c1b230206773c3a54fc7f2fc2d0))
* Exposed variables from blueprint_config ([94bb92c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/94bb92c1d935ef703687ecff226f3bdc3c41e643))
* set variables map ([c3d2a65](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/c3d2a651efdd2cd7318bd5e2f5f0aff330bc1545))
* Updated reference example ([3132e94](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/3132e940e64d90e0d6bea60d8951460c33aa1904))

# [0.3.0-pre.49](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.48...0.3.0-pre.49) (2023-02-08)


### Bug Fixes

* **hcl:** when a user overrides a variable, should be always used in hcl ([dae126c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/dae126c22fb6432864e605e1aa1945119cb902f1))
* **omits:** double omits are not impacting usedHCL if is overriden ([9b4a8ef](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9b4a8efa95f3debf01b1867a9a31de8dd1e217f7))

# [0.3.0-pre.48](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.47...0.3.0-pre.48) (2023-02-08)


### Features

* **ep-2914:** iac list of objects ([52eaf11](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/52eaf116183ed534b317f247eaf54db334edbb47))

# [0.3.0-pre.47](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.46...0.3.0-pre.47) (2023-02-06)


### Features

* **plugin:** Add flag to gate the new blueprintconfig ([34960e3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/34960e3fff554baf4821db1df0097f24b1568690))
* **plugin:** EXPERIMENTAL first approach to support low level lib ([15513bc](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/15513bc789e3f74af7ea8da0d321aa31ebd969be))

# [0.3.0-pre.46](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.45...0.3.0-pre.46) (2023-02-06)


### Bug Fixes

* **sdk:** get develop version of sdk ([55244ee](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/55244ee57e76a97b5fc34c5c8e14065665ffe781))


### Features

* **omits:** process omitted variables ([f4c3f23](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f4c3f23e7e4d429fa30fb5db534771f156008226))
* **omits:** process omitted variables ([ba181fb](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/ba181fbe0f7490945cd0cc82d8bb436ad810af76))

# [0.3.0-pre.45](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.44...0.3.0-pre.45) (2023-02-01)


### Features

* **ep-2845:** iac blueprint module ids ([58bcf01](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/58bcf01bfa0ad4d024fc0a61a6f43db365cfeff9))
* **ep-2845:** iac blueprint module ids ([1ce8f34](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1ce8f341effc437bf01b3ed5478d54db06f1965e))

# [0.3.0-pre.44](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.43...0.3.0-pre.44) (2023-01-26)


### Bug Fixes

* **conditionals:** use the same content shape as variable for conditionals ([d3072c3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/d3072c3a1e1650236450d3fb10b9d2ee20da6d56))

# [0.3.0-pre.43](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.42...0.3.0-pre.43) (2023-01-25)


### Features

* **ep-2838:** iac override value ([c5d1b90](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/c5d1b900be4dce375b61c1455bf836b61445d917))

# [0.3.0-pre.42](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.41...0.3.0-pre.42) (2023-01-20)


### Bug Fixes

* **ep-2773:** iac shared overridden variables ([10ab06a](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/10ab06aa1a7bcfaaac52aaa3dd05d69c5f04646a))

# [0.3.0-pre.41](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.40...0.3.0-pre.41) (2023-01-20)


### Features

* Added examples with simple objects fields ([f733591](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f733591eccb6e261247d9a8526758035b503bc69))

# [0.3.0-pre.40](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.39...0.3.0-pre.40) (2023-01-18)


### Bug Fixes

* **ep-2768:** iac file block vars ([73e63e5](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/73e63e545d021da68fa052221d891c90e2db0e35))

# [0.3.0-pre.39](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.38...0.3.0-pre.39) (2023-01-18)


### Bug Fixes

* **makefile:** comment version ([e87f68c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e87f68cb2002d0e050bbadf664601bd4f2657139))

# [0.3.0-pre.38](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.37...0.3.0-pre.38) (2023-01-17)


### Features

* **variable_override:** modify children's variables by reference ([a1a40aa](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/a1a40aaf61e8522befc78bcecbe7491859c9c23f))

# [0.3.0-pre.37](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.36...0.3.0-pre.37) (2023-01-17)


### Features

* **ep-2617:** iac conditionals ([af2d03d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/af2d03dbe822e1b667be1c8b374583c727e48afe))

# [0.3.0-pre.36](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.35...0.3.0-pre.36) (2023-01-17)


### Features

* **eks_demo:** provisional eks demo ([6c185a9](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/6c185a9aad91c71c676d241669c20992e614bd02))

# [0.3.0-pre.35](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.34...0.3.0-pre.35) (2023-01-16)


### Features

* **ep-2770:** iac conditionals content value ([fdaef35](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/fdaef3564f000b9fb3d4360628c8c8ed6df32d0f))

# [0.3.0-pre.34](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.33...0.3.0-pre.34) (2023-01-13)


### Bug Fixes

* Default options from module for override variables ([257124a](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/257124a3895134c62d09968a65aef8afe13904ae))
* Updated override_var example ([22cac2c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/22cac2cdb558e0af5be363d480cca3d954db70ac))

# [0.3.0-pre.33](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.32...0.3.0-pre.33) (2023-01-13)


### Features

* **ep-2769:** iac variable.value ([b512139](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/b512139fd61226a36b1bf6d0d9cae50ef6dae73d))

# [0.3.0-pre.32](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.31...0.3.0-pre.32) (2023-01-10)


### Features

* **ep-2685:** iac conditionals lists and required list values ([2a9cd10](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/2a9cd1011f642edba0b2d6956adedd9fb7e8581e))

# [0.3.0-pre.31](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.30...0.3.0-pre.31) (2023-01-10)


### Bug Fixes

* Fixed issues with field types ([aec64c4](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/aec64c486071c343f32ff516611cc3ff19da8127))
* Turned off linter for gocyclo because of GetBlueprintConfigFromSchema function complexity ([e124be7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e124be787a516740a18f7e98c9befd951d0069ee))


### Features

* Added map type to provider ([4b58396](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/4b5839634b1b2cafa59796a7583d38e8caf90f1a))
* Added maps example ([2b169e3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/2b169e3db8d8c8adc890c86accadb0d192826d9f))

# [0.3.0-pre.30](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.29...0.3.0-pre.30) (2023-01-06)


### Features

* **ep-2628:** iac lists ([979bfbb](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/979bfbb94d6a5bd5fdba25187b3aa5852c7eb54d))

# [0.3.0-pre.29](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.28...0.3.0-pre.29) (2023-01-05)


### Features

* **tree:** parse json to include overridevariables ([57b6332](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/57b6332ae5a8b356959a80c57ee04ce654a16766))

# [0.3.0-pre.28](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.27...0.3.0-pre.28) (2023-01-05)


### Bug Fixes

* Fixed source_blueprint example ([cf22f35](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/cf22f35030bfe942946500c1177fde142583ee1f))

# [0.3.0-pre.27](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.26...0.3.0-pre.27) (2023-01-04)


### Bug Fixes

* linting issues ([e879caf](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e879cafb0eff3cbd98e498de10efc728263df866))


### Features

* **conditionals:** integrate conditionals with tree ([1dc8f08](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1dc8f08d8e5877882b393353b7f7d2fa7b1ef209))
* **schema:** remove form options ([e50f418](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e50f418e69bb76ec1109f3430058045992cc44ed))
* **schema:** remove form options ([bb04707](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/bb047072b7e4937c2359247ef36a7d94fe61aeb2))
* **tree:** process variables by tree level ([28fd1b2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/28fd1b260bacc85eb10bd3f49ac7cc1d3e7cf1d8))

# [0.3.0-pre.26](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.25...0.3.0-pre.26) (2023-01-03)


### Features

* **ep-2617:** iac conditionals ([7ce70e3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/7ce70e36b5a0208b1ae29efcd72b93ffb7531c1c))
* **ep-2617:** iac conditionals (fmt to log) ([03fc304](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/03fc30421f2187130ec83ed62eb60db98a6514fe))

# [0.3.0-pre.25](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.24...0.3.0-pre.25) (2022-12-28)


### Bug Fixes

* Fixed missing types for overrides variables ([50bdd4d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/50bdd4d44fe0dbd9dc2fa7d69dfc69089ec96095))

# [0.3.0-pre.24](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.23...0.3.0-pre.24) (2022-12-22)


### Features

* **blueprint:** add tf source example ([0809365](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/0809365dd68511a3f1235014a162d238a0c07266))

# [0.3.0-pre.23](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.22...0.3.0-pre.23) (2022-12-22)


### Bug Fixes

* Renabled omit and override variables for new blueprint_config ([55c1664](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/55c16644cc31ab0f400f29caac5fac82ff60bfae))


### Features

* **blueprint:** add tf example ([7be7c32](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/7be7c3202e2c268dc51dbe1ed64e6e6977a1ab4a))
* **blueprint:** get form shape from tree ([3a2136b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/3a2136b854e5e049b717de4839666032f7942f57))

# [0.3.0-pre.22](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.21...0.3.0-pre.22) (2022-12-22)


### Bug Fixes

* Fixed linter issues ([ba0772b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/ba0772b49635c74a324d664833e760b037b1cb78))


### Features

* **blueprint:** add context variable to blueprint resource ([1897d52](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/1897d52efe8ed78afbb11e1b93a95bab3e98dc53))
* **blueprint:** use "source" keyword ([34051f5](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/34051f5f96571b88833f0be93b8f86e63929a542))
* Built generic form using custom variables ([9491b9b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9491b9b23b248dbbc004a82e375b87e6d3e5deb2))
* Defined TODO for refactors ([654942e](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/654942ea541e2e534d824c3d2c42e3a1a3137eb3))
* Supported generic and module blueprint_config ([8e69c5c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/8e69c5c1051d720b4ad6311a7e8ba6ff9470d6e9))

# [0.3.0-pre.21](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.20...0.3.0-pre.21) (2022-12-14)


### Features

* **ep-2588:** renaming autocloud_module form_config into blueprint_config ([83bb2fe](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/83bb2fe22de37909061ba035cacdef4f8f2ac7b6))

# [0.3.0-pre.20](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.19...0.3.0-pre.20) (2022-12-14)


### Features

* **variable_override:** set default value and add regex to detect ouputs ([e7f9253](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e7f925305c519e9495c3e357a717aab11946c892))

# [0.3.0-pre.19](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.18...0.3.0-pre.19) (2022-12-14)


### Bug Fixes

* adding AllowConsumerToEdit to true when overriding a variable ([487845f](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/487845fc20c5e75e48918ebe9e3d56f441dffddd))

# [0.3.0-pre.18](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.17...0.3.0-pre.18) (2022-12-13)


### Features

* Renamed fields_options to options ([63aa877](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/63aa877bcfee416441aef5c6ed854af98bd259a2))

# [0.3.0-pre.17](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.16...0.3.0-pre.17) (2022-12-13)


### Features

* Replaced override_variable block with variable ([f6789a6](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f6789a64a6283cbd7850b925af03041b55a0f88e))

# [0.3.0-pre.16](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.15...0.3.0-pre.16) (2022-12-08)


### Bug Fixes

* Renamed terraform_processor to blueprint_config ([8905154](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/8905154165f4110f3019092c4001f2345ea89269))

# [0.3.0-pre.15](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.14...0.3.0-pre.15) (2022-12-07)


### Bug Fixes

* Added better response for http error fetching repositories ([6609a73](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/6609a73eb85bcb86407fd6c7e19f04f31780b144))

# [0.3.0-pre.14](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.13...0.3.0-pre.14) (2022-12-06)


### Features

* **ep-2575:** IAC Module composability - updating tests ([64ee6a8](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/64ee6a84a8d330d39ef1206887ab45d933dac004))
* **ep-2575:** IAC Module composability through the outputs of another module ([acb071d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/acb071d43c5b4d91392e1ec6389dddadf456e8f7))

# [0.3.0-pre.13](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.12...0.3.0-pre.13) (2022-12-02)


### Bug Fixes

* Added default value for tags field ([ab3c5a7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/ab3c5a7389a9bd6d89e69522349f13f421a7309b))

# [0.3.0-pre.12](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.11...0.3.0-pre.12) (2022-12-01)


### Features

* **ep-2578:** blueprint file block update ([9e20850](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/9e20850031b477c570f79e7b1413981396f4c004))

# [0.3.0-pre.11](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.10...0.3.0-pre.11) (2022-11-29)


### Features

* **ep-2577:** autocloud provider endpoint ([84ccd11](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/84ccd111f8cc0a306c2ef92068ef7b80e653a425))

# [0.3.0-pre.10](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.9...0.3.0-pre.10) (2022-11-25)


### Features

* **ep-2559:** blueprint: limit git_config to 1 block ([73a282d](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/73a282deaa5da328d6c0d6c7287caa8ff24bd425))

# [0.3.0-pre.9](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.8...0.3.0-pre.9) (2022-11-25)


### Features

* **ep-2562:** override variable: display name and helper text ([3943b28](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/3943b28b6a4e98acd43b1b78ee0ddba2e5a4772c))

# [0.3.0-pre.8](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.7...0.3.0-pre.8) (2022-11-24)


### Features

* Added display_order field to module schema ([e1cf88c](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/e1cf88c6eab6ebc25f3ee16e824fa45606307ef1))

# [0.3.0-pre.7](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.6...0.3.0-pre.7) (2022-11-22)


### Features

* **validation:** add validation to autocloud_module name ([4f3a18a](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/4f3a18addf03f6499060ffea044dafcaf4ed341d))

# [0.3.0-pre.6](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.5...0.3.0-pre.6) (2022-11-18)


### Features

* **gitconfig:** avoid configuration errors with git default repo ([b5c3c9f](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/b5c3c9fc335ac9b5b4395b83b14b26645a0ec5f3))

# [0.3.0-pre.5](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.4...0.3.0-pre.5) (2022-11-17)


### Features

* Added tags_variables field to module ([f2c6d6b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/f2c6d6b1a7a36da2ea1e3c905c3d9a4d28c85500))

# [0.3.0-pre.4](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.3...0.3.0-pre.4) (2022-11-15)


### Features

* **ep-2502:** form builder ([d0c7100](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/d0c71005bba0a007155ab45dbbeb218e300d3b71))

# [0.3.0-pre.3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.2...0.3.0-pre.3) (2022-11-04)


### Features

* **EP-2503:** IAC token auth ([7058241](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/7058241d6758c66c88fb24f481a07d12995b4654))

# [0.3.0-pre.2](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.3.0-pre.1...0.3.0-pre.2) (2022-11-02)


### Bug Fixes

* **ep-2499:** fixing tests ([501d21a](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/501d21aebf8178c2c588592aa7a90bd18fb2c840))


### Features

* **EP-2499:** iac blueprint module override ([2c998a3](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/2c998a342145efb9ac564206413576d531167de5))

# [0.3.0-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.2.1-pre.1...0.3.0-pre.1) (2022-10-31)


### Features

* **EP-2466:** modules variables id output ([82f023b](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/82f023b77e772535ed1f00863381a2e8232fefb5))

## [0.2.1-pre.1](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/compare/0.2.0...0.2.1-pre.1) (2022-10-31)


### Bug Fixes

* **examples:** add stub basic example ([08fd505](https://gitlab.com/auto-cloud/infrastructure/public/terraform-provider/commit/08fd505ac489bbcf3c6971c516dc010aa09d2dbb))
