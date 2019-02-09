# inifile

Package inifile implements parsing a simple ini-file

[![GoDoc](https://godoc.org/github.com/codeation/inifile?status.svg)](https://godoc.org/github.com/codeation/inifile)

# Installation

To install inifile package:

```
go get -u github.com/codeation/inifile
```

# Examples

sample.ini file:

```
port=8080
host=127.0.0.1

[node]
name=server.local
ip=10.0.0.11
```

sample.go file:

```
package main

import (
	"fmt"

	"github.com/codeation/inifile"
)

func main() {
	ini, err := inifile.Read("sample.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("port=", ini.Get("", "port"))
	fmt.Println("node.name=", ini.Get("node", "name"))
}
```

output:
```
$ go run sample.go 
port= 8080
node.name= server.local
```

# Installation

To install inifile package:

```
go get -u github.com/codeation/inifile
```

# Environment variable

You can specify a full path to ini-file through the environment variable. For example:

```
$ SAMPLE_INI=/etc/server_configuration.ini go run sample.go
```

The name of the environment variable must match the inifile.Read parameter.
The base file name is used, the directory name is ignored.
The letters of the file name are written in uppercase, the dot is replaced by an underscore.
