# bodyless

Bodyless is cli (command line user interface for [BodylessCMS](https://github.com/chaitanya11/BodylessCMS/). we can create, build, and deploy BodylessCMS project with this cli. This is writen in go.


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
deploy
```


## Usage
### To create project.
```
bodyless create -w bodyless -P /tmp/ -N shadow -p default
```