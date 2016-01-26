# ecr-login

Login tool for AWS Container Registry.

This is a lightweight golang version of the AWS command-line utility
`aws ecr get-login`, designed to build into a small scratch docker
image.

Can also produce output in other formats using golang templates.

## Installation

See build or docker image below.

## Usage

Login to your AWS Container Registry:

```
$ eval $(./ecr-login)
WARNING: login credentials saved in /Users/ric/.docker/config.json
Login Succeeded
```

Alternatively, you can use the included templates to output docker
config format directly and redirect output to `~/.docker/config.json`
or `~/.dockercfg`:

```
$ TEMPLATE=templates/config.tmpl ./ecr-login
{
        "auths": {
                "https://1234567890.dkr.ecr.us-east-1.amazonaws.com": {
                        "auth": "...",
                        "email": "none"
                }
         }
}
```

## Systemd example

This is an example of how I use `ecr-login` with systemd units on
CoreOS:

```
[Unit]
Description=Example

[Service]
User=core
Environment=AWS_REGION=us-east-1
ExecStartPre=/bin/bash -c 'eval $(docker run -e AWS_REGION rlister/ecr-login)'
ExecStartPre=-/usr/bin/docker rm example
ExecStartPre=/usr/bin/docker pull 1234567890.dkr.ecr.us-east-1.amazonaws.com/example:latest
ExecStart=/usr/bin/docker run --name example 1234567890.dkr.ecr.us-east-1.amazonaws.com/example:latest
ExecStop=/usr/bin/docker stop example
```

## Build from source

```
go build ./ecr-login.go
```

## Docker image

```
version=0.0.1
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ecr-login.go
docker build -t rlister/ecr-login:${version} .
docker tag -f rlister/ecr-login:${version} rlister/ecr-login:latest
docker push rlister/ecr-login
```