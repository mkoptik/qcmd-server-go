# qcmd - terminal command full text search

qcmd project is a hammer for people like me who don't remember terminal commands :-). Project is divided into
multiple repositories for server, client and commands list. 

## Server side

Server side part written in Golang is REST API provider using full text search to get a commands. For the full text
search is used bleve package, offering no dependency to external service.

### Start in docker

All changes must be pushed to git repository in order to correctly run this package in docker

```
docker run golang go get -v github.com/mkoptik/qcmd-server-go
```

### REST API endpoints

TODO: Finish