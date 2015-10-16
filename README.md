Diego-Enabler: CLI Plugin
=====================
**This plugin is for [CLI v6.13.0+](https://github.com/cloudfoundry/cli/releases)**

This plugin enable an app for Diego support, For more detail information of running apps on Diego, see [here](https://github.com/cloudfoundry-incubator/diego-design-notes/blob/master/migrating-to-diego.md)

##Installation

#####Install from Repo (v.6.10.0+)
  ```
  $ cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
  $ cf install-plugin Diego-Enabler -r CF-Community
  ```
  
#####Install from Url (v.6.8.0+)
OSX
  ```
  cf install-plugin https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/download/v1.0.0/diego-enabler_darwin_amd64
  ```

linux64:
  ```
  cf install-plugin https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/download/v1.0.0/diego-enabler_linux_amd64
  ```

linux32:
  ```
  cf install-plugin https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/download/v1.0.0/diego-enabler_linux_386
  ```

windows64:
  ```
  cf install-plugin https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/download/v1.0.0/diego-enabler_windows_amd64.exe
  ```
  
windows32:
  ```
  cf install-plugin https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/download/v1.0.0/diego-enabler_windows_386.exe
  ```


#####Install from Binary file (v.6.7.0)


- Download the binary [`win32`](https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/download/v1.0.0/diego-enabler_windows_386.exe) [`win64`](https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/download/v1.0.0/diego-enabler_windows_amd64.exe) [`osx`](https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/download/v1.0.0/diego-enabler_darwin_amd64) [`linux32`](https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/download/v1.0.0/diego-enabler_linux_386) [`linux64`](https://github.com/cloudfoundry-incubator/Diego-Enabler/releases/download/v1.0.0/diego-enabler_linux_amd64)
- Install plugin `$ cf install-plugin <binary_name>`



##Full Command List

| command | usage | description|
| :--------------- |:---------------| :------------|
|`enable-diego`| `cf enable-diego App_Name` |enable diego for an app|
|`disable-diego`| `cf disable-diego App_Name` |disable diego for an app|
|`has-diego-enabled`| `cf has-diego-enabled App_Name` |check if diego is enabled for an app|
