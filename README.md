# Rotating File writer

An `io.writer` & `io.Closer` compliant file writer that will always write to the path that you give it, even if somebody deletes/renames that path out from under you.

Created so that Go programs can be used with the standard Linux `logrotate`: writes following a removal / rename will occur in a newly created file rather than the previously opened filehandle.

## Example

``` go
package main

import (
  "github.com/mipearson/rfw"
  "log"
)

func main() {
  writer, err := rfw.Open("/var/log/myprogram", 0644)
  if err != nil {
    log.Fatalln("Could not open '/var/log/myprogram': ", err)
  }

  log := log.New(writer, "[myprogram] ", log.LstdFlags)
  log.Println("Logging as normal")
}
```
