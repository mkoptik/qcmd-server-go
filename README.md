# qcmd - search for terminal command

qcmd project is a tool for people like me who don't remember terminal commands. It is a utility to full-text search
terminal command examples, so the search can be performed by non exact text query.

Project is divided into multiple repositories for server, client and commands list. 

## Server side

Server side part written in Golang is API provider using full text search to get a commands. As a fulltext
search engine is used bleve, offering no dependency to external service.

Index is stored in ~/.qcmd/... directory

### Compile sources and run

```
go build && ./qcmd-server-go
```

### Run in docker

Build docker image

```
docker build -t mkoptik/qcmd-server-go .
```

Start docker container

```
docker run -d -p 8888:8888 mkoptik/qcmd-server-go
```

## API endpoints

### Search commands
```
GET /command/search?search=<search-string>&tag=<tag-filter>
```

Parameters:

* search - mandatory
* tag - optional, multivalue (when multiple values specified, all must match)

### List all commands

```
GET /command/all?tag=<tag-filter>&tag=<tag-filter>
```

Parameters:

* tag - optional, multivalue (when multiple values specified, all must match)

### Search tags

```
GET /tag/search?search=<search-string>
```