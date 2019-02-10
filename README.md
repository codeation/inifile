# inifile

Package inifile implements parsing a simple INI file

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

# Environment variable

You can specify a full path to INI file via the environment variable. For example:

```
$ SAMPLE_INI=/etc/server_configuration.ini go run sample.go
```

The environment variable is checked only if
the parameter of inifile.Read does not contain the directory name.
The name of the environment variable is the file name in uppercase,
the dot is replaced by an underscore.

# Command output encapsulation

**Please note that executing external commands can lead to application vulnerabilities.
Command output encapsulation is disabled by default.**

You can use command output as INI file variable value. For example, as string of INI file:

```
one_time_password=$(pwgen -s 16)
```

To enable command output encapsulation, call inifile.Command in golang file:

```
    ini, _ := inifile.Read("sample.ini")
    ini.Command(true)
    fmt.Println("One time password is", ini.Get("", "one_time_password"))
```

See the [documentation](https://godoc.org/github.com/codeation/inifile) for details.
