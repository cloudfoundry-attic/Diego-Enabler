Diego-Enabler: CLI Plugin
=====================
**This plugin is for [CLI
v6.13.0+](https://github.com/cloudfoundry/cli/releases).** For CLI v6.12.4 and older, use the
[Diego-Beta](https://github.com/cloudfoundry-incubator/diego-cli-plugin) plugin
instead.

This plugin enables Diego support for an app running on Cloud Foundry. For more
information about running apps on Diego, see the [Migrating to Diego](https://github.com/cloudfoundry-incubator/diego-design-notes/blob/master/migrating-to-diego.md)
guide.

##Full Command List

Command             |Usage                             |Description
---                 |---                               |---
`enable-diego`      | `cf enable-diego App_Name`       |enable Diego for an app
`disable-diego`     | `cf disable-diego App_Name`      |disable Diego for an app
`has-diego-enabled` | `cf has-diego-enabled App_Name`  |check if Diego is enabled for an app

##Installation

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

