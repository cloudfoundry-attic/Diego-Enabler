Diego-Enabler: CLI Plugin
=====================
**This plugin is for [CLI
v6.13.0+](https://github.com/cloudfoundry/cli/releases).** For CLI v6.12.4 and older, use the
[Diego-Beta](https://github.com/cloudfoundry-incubator/diego-cli-plugin) plugin
instead.

This plugin enables Diego support for an app running on Cloud Foundry. For more
information about running apps on Diego, see the [Migrating to Diego](https://github.com/cloudfoundry-incubator/diego-design-notes/blob/master/migrating-to-diego.md)
guide.

## Full Command List

Command             |Usage                                                                        |Description
---                 |---                                                                          |---
`enable-diego`      | `cf enable-diego App_Name`                                                  |Migrate app to the Diego runtime
`disable-diego`     | `cf disable-diego App_Name`                                                 |Migrate app to the DEA runtime
`has-diego-enabled` | `cf has-diego-enabled App_Name`                                             |Report whether an app is configured to run on the Diego runtime
`diego-apps`        | `cf diego-apps [-o ORG]`                                                    |Lists all apps running on the Diego runtime that are visible to the user
`dea-apps`          | `cf dea-apps [-o ORG]`                                                      |Lists all apps running on the DEA runtime that are visible to the user
`migrate-apps`      | <code>cf migrate-apps (diego &#124; dea) [-o ORG] [-p MAX_IN_FLIGHT]</code> |Migrate all apps to Diego/DEA

## Installation

To install the plugin from the CF Community repository:

```
$ cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
$ cf install-plugin Diego-Enabler -r CF-Community
```

To install the plugin from GitHub:

First, copy the URL or download the binary for your platform from the [latest release page](https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/latest).

Then call `cf install-plugin` with either the URL you copied or the binary you downloaded:

```
  cf install-plugin [URL|binary]
```

## Release

To create release.

1. [Gox](https://github.com/mitchellh/gox) must be installed
1. [Github release](https://github.com/aktau/github-release) must be installed
1. Run `bin/build`
1. Export Github personal access token to `$GITHUB_TOKEN`
1. Make sure the release version is updated in the file `VERSION`
1. Run `bin/release`
