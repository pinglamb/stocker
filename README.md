# Stocker

For local development, most of the time you need some sort of 3rd party services like Postgres, Redis, Elasticsearch, etc. Installing them to your local machine using Homebrew is okay at the beginning, but imagine how troublesome it is after running `brew upgrade`, or you have another project which depends on a version that differs from what you have installed. Running them in docker containers is a much better solution.

[Stocker](https://github.com/pinglamb/stocker) is a tool for helping you to bootstrap the `docker-compose.yml` set up, which you can use for booting up the services in containers.

## vs Vagrant

Vagrant is a great tool for setting up isolated development environment, but using VirtualBox is just too slow and resources intensive. Using docker for that should be more lightweight, especially with the release of [Native Docker on Mac and Windows](https://blog.docker.com/2016/03/docker-for-mac-windows-beta/) which uses the native virtualization of the OS.

## Install

```
curl -L https://github.com/pinglamb/stocker/releases/download/0.1/stocker-`uname -s`-`uname -m` > /usr/local/bin/stocker
```

P.S. It currently works on Mac only.

## How to use

Inside your project, you can run:

```
stocker add postgres
```

It will add `postgres` to your `docker-compose.yml`, it will also map the `ports` the service exposes to your host machine. You might need to update your configuration file for that.

---

```
stocker up
```

## Contribute

Feel Free.
