[![Build Status](https://github.com/beyondstorage/go-service-ipfs/workflows/Unit%20Test/badge.svg?branch=master)](https://github.com/beyondstorage/go-service-ipfs/actions?query=workflow%3A%22Unit+Test%22)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/Xuanwo/storage/blob/master/LICENSE)
[![](https://img.shields.io/matrix/beyondstorage@go-storage:matrix.org.svg?logo=matrix)](https://matrix.to/#/#beyondstorage@go-storage:matrix.org)

# go-service-ipfs

[InterPlanetary File System(IPFS)](https://ipfs.io/) support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get github.com/beyondstorage/go-service-ipfs
```

## Usage

```go
import (
	"log"

	_ "github.com/beyondstorage/go-service-ipfs"
	"github.com/beyondstorage/go-storage/v4/services"
)

func main() {
	store, err := services.NewStoragerFromString("ipfs:///path/to/workdir?endpoint=<ipfs_http_api_endpoint>&gateway=<ipfs_http_gateway>")
	if err != nil {
		log.Fatal(err)
	}
	
	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/ipfs) about go-service-ipfs.
