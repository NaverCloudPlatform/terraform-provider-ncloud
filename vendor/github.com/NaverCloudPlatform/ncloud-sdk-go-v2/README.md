# ncloud-sdk-go-v2

ncloud-sdk-go-v2 is the official Naver Cloud Platform SDK for the Go programming language.

### Installing

```
go get github.com/NaverCloudPlatform/ncloud-sdk-go-v2
```

### Example

```
package main

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"log"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
)

func main() {

	client := server.NewAPIClient(server.NewConfiguration(&server.APIKey{
		AccessKey: "accessKey",
		SecretKey: "secretKey",
	}))

	// Create server instance
	req := server.CreateServerInstancesRequest{
		ServerImageProductCode:     ncloud.String("SPSW0LINUX000031"),
		ServerProductCode:          ncloud.String("SPSVRSTAND000049"),
		UserData:                   ncloud.String("#!/bin/sh\nyum -y install httpd"),
		IsProtectServerTermination: ncloud.Bool(false),
		ServerCreateCount:          ncloud.Int32(1),
	}


	if r, err := client.V2Api.CreateServerInstances(&req); err != nil {
		log.Println(err)
	} else {
		sList := r.ServerInstanceList
		log.Println(ncloud.StringValue(r.RequestId))
		log.Println(ncloud.StringValue(sList[0].ServerInstanceNo))
		log.Println(ncloud.StringValue(sList[0].ServerName))
	}
}
```

## Documentation for individual modules

| Services       | Documentation                                                                                                           |
| -------------- | ------------------------------------------ |
| _Server_       | [**Server**](server/README.md)             |
| _Loadbalancer_ | [**Loadbalancer**](loadbalancer/README.md) |
| _Autoscaling_  | [**Autoscaling**](autoscaling/README.md)   |
| _Monitoring_   | [**Monitoring**](monitoring/README.md)     |
| _CDN_          | [**CDN**](cdn/README.md)                   |
| _CloudDB_      | [**CloudDB**](clouddb/README.md)           |


### License

```
Copyright (c) 2018 NAVER BUSINESS PLATFORM Corp.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
