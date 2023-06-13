# Changelog

## [0.3.8](https://github.com/shufo/backlog-cli/compare/v0.3.7...v0.3.8) (2023-06-13)


### Bug Fixes

* **deps:** update github.com/kenzo0107/backlog digest to 44e2d1d ([9ff0a20](https://github.com/shufo/backlog-cli/commit/9ff0a20616b6b555a0fb6728e2c314a23fac9359))
* **deps:** update module github.com/charmbracelet/bubbletea to v0.24.2 ([8efc39b](https://github.com/shufo/backlog-cli/commit/8efc39b3b2fc2c79584cbe8ea2445398d4ca656e))
* **deps:** update module github.com/urfave/cli/v3 to v3.0.0-alpha3 ([35123bc](https://github.com/shufo/backlog-cli/commit/35123bc0ae12e6897decfcb21632008122235dfc))


### Miscellaneous

* 🤖 change binary name from `bk` to `backlog` ([7d19e31](https://github.com/shufo/backlog-cli/commit/7d19e31305ec7b4c41081fa3359abae1d3fff989))

## [0.3.7](https://github.com/shufo/backlog-cli/compare/v0.3.6...v0.3.7) (2023-03-21)


### Bug Fixes

* 🐛 order issue comments to asc ([f1f0bff](https://github.com/shufo/backlog-cli/commit/f1f0bffae2cb454b705c3e8009fcea562a0d5d92))

## [0.3.6](https://github.com/shufo/backlog-cli/compare/v0.3.5...v0.3.6) (2023-03-21)


### Bug Fixes

* 🐛 invalid int flag ([e39fe9a](https://github.com/shufo/backlog-cli/commit/e39fe9a543dc938b098739d2e14e409eba787d9e))

## [0.3.5](https://github.com/shufo/backlog-cli/compare/v0.3.4...v0.3.5) (2023-03-21)


### Bug Fixes

* 🐛 create app config directory if not exists ([d257870](https://github.com/shufo/backlog-cli/commit/d2578701c83631155b863da91def614969a332de))

## [0.3.4](https://github.com/shufo/backlog-cli/compare/v0.3.3...v0.3.4) (2023-03-21)


### Bug Fixes

* 🐛 error when hosts config does not exists ([5b1dc63](https://github.com/shufo/backlog-cli/commit/5b1dc63935c069bf67980e5b7f5e3f7f6e282529))

## [0.3.3](https://github.com/shufo/backlog-cli/compare/v0.3.2...v0.3.3) (2023-03-21)


### Bug Fixes

* 🐛 exit status 1 on open url in windows ([688433d](https://github.com/shufo/backlog-cli/commit/688433d4675f0fd06e7b47522455a230026d3ee6))

## [0.3.2](https://github.com/shufo/backlog-cli/compare/v0.3.1...v0.3.2) (2023-03-21)


### Bug Fixes

* 🐛 auth login does not work on wsl on windows ([ab245cf](https://github.com/shufo/backlog-cli/commit/ab245cffc6fbdeb0a0bdb7cea55ce2ad0fadd9ab))

## [0.3.1](https://github.com/shufo/backlog-cli/compare/v0.3.0...v0.3.1) (2023-03-21)


### Bug Fixes

* 🐛 auth login does not work on windows ([09d2416](https://github.com/shufo/backlog-cli/commit/09d2416ba76f8f42fb7cafed30158280941882cb))

## [0.3.0](https://github.com/shufo/backlog-cli/compare/v0.2.1...v0.3.0) (2023-03-21)


### Features

* 🎸 `issue comment` command ([c1b855c](https://github.com/shufo/backlog-cli/commit/c1b855c5b88bb31a2ec5396eeb3adbed4296dbbb))


### Bug Fixes

* **deps:** update module github.com/shufo/find-config to v0.1.1 ([ef7d7fe](https://github.com/shufo/backlog-cli/commit/ef7d7feb33b3e6d24e866f80e74527e50e44a1fc))


### Miscellaneous

* 🤖 bump urfave/cli from v2 to v3 ([a008c38](https://github.com/shufo/backlog-cli/commit/a008c3838653d16bf5fca3bbb69f4cf4a756fa5b))
* 🤖 upgrade dependencies ([68e780f](https://github.com/shufo/backlog-cli/commit/68e780f27155ef822d55a21b5ab6f399cdaa225b))
* **deps:** bump golang.org/x/crypto ([7d43439](https://github.com/shufo/backlog-cli/commit/7d43439a49d21601c148eeb6e01b19228a5b9057))

## [0.2.1](https://github.com/shufo/backlog-cli/compare/v0.2.0...v0.2.1) (2023-03-20)


### Bug Fixes

* 🐛 change usage text from bl to bk ([cf4d4d8](https://github.com/shufo/backlog-cli/commit/cf4d4d808f2d76b0a66ab443dee5f1e8380be1a3))
* 🐛 spinner does not stopped on error ([7f23e35](https://github.com/shufo/backlog-cli/commit/7f23e35d401ef4323db289277bf80af223fb9299))
* 🐛 use uuid for auth code ([1eed17a](https://github.com/shufo/backlog-cli/commit/1eed17acd3acf492e6ca0a1d28bd6cd182c01d0f))


### Miscellaneous

* 🤖 add entrypoint for install command ([472753b](https://github.com/shufo/backlog-cli/commit/472753bd19574e64ecc15a9d6b4220b011306447))

## [0.2.0](https://github.com/shufo/backlog-cli/compare/v0.1.0...v0.2.0) (2023-03-19)


### Features

* 🎸 support --web flag for list command ([d104a0a](https://github.com/shufo/backlog-cli/commit/d104a0a32560f8b1a1191167b23df0496e7c1df4))


### Miscellaneous

* 🤖 change binary name from ba to bk ([e0a0322](https://github.com/shufo/backlog-cli/commit/e0a0322436d2adebb8d7c24657df15d550971894))
* 🤖 change binary name to bk ([76da1c0](https://github.com/shufo/backlog-cli/commit/76da1c0b9e8d20feabbdeb9811551b0a70513953))

## 0.1.0 (2023-03-18)


### Features

* 🎸 add installation script ([2a45c44](https://github.com/shufo/backlog-cli/commit/2a45c443c4d38fdc78a11bd9a7ffc2d8cbff20f1))


### Miscellaneous

* 🤖 bump go version to 1.18 ([462d179](https://github.com/shufo/backlog-cli/commit/462d17992743d9f78fd8a7ad0bd1a19452595113))
* Initialize ([4630906](https://github.com/shufo/backlog-cli/commit/4630906aa03a5068c6c6a8805ca6eabef1232eb6))
* release 0.1.0 ([cd9dc89](https://github.com/shufo/backlog-cli/commit/cd9dc893c2ad97eddae151710adba2ff85a89902))
