# Bodyless-cli

Bodyless is cli (command line user interface for [BodylessCMS]("https://github.com/chaitanya11/BodylessCMS/")). we can create, build, and deploy BodylessCMS project with this cli. This is writen in go.


## How to Install
To install this package do the following.
```
$ git clone https://github.com/chaitanya11/bodyless-cli.git $GOPATH/src/bodyless-cli
$ cd $GOPATH/src/bodyless-cli/
$ go build -o bodyless
$ chmod +x bodyless
$ cp bodyless /usr/local/bin/
```
Now test installation as bellow.

```
$ bodyless
usage: bodyless <command> [<args>]
commands:
create
  -w, --CodeBucketName string
    	Name of the bucket where website code is deployed.
  -P, --Path string
    	Project Location. (default ".")
  -N, --ProjectName string
    	Name of the project.
  -p, --profile string
    	Name of the aws profile configured. (default "default")
  -r, --region string
    	Name of the aws region. (default "us-west-2")
build
  -p, --path string
    	Project Location. (default ".")
deploy
  -p, --path string
    	Project Location. (default ".")
remove
  -p, --path string
    	Project Location. (default ".")
```


## Usage
### To create project.
```
bodyless create -w bodyless -P /tmp/ -N shadow -p default
```

### To remove all created aws resources and project
```
bodyless remove -p /tmp/shadow/
```

If you want to do any customisation to [bodylesscms]("https://github.com/chaitanya11/BodylessCMS/"), go to the given path in create command and do necessary changes and to build or deploy follow these steps.


### To build [bodylesscms]("https://github.com/chaitanya11/BodylessCMS/") project (after any customisations)

```
bodyless build -p /tmp/shadow/
```


### To deploy [bodylesscms]("https://github.com/chaitanya11/BodylessCMS/") project (after any customisations)

```
bodyless deploy -p /tmp/shadow/
```