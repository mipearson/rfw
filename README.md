# Rotating File writer

[![wercker status](https://app.wercker.com/status/b834fcde90ed6fdf2fc5e0a3ca98d518/s/ "wercker status")](https://app.wercker.com/project/bykey/b834fcde90ed6fdf2fc5e0a3ca98d518) [![GoDoc](https://godoc.org/github.com/mipearson/rfw?status.png)](https://godoc.org/github.com/mipearson/rfw)

An `io.writer` & `io.Closer` compliant file writer that will always write to the path that you give it, even if somebody deletes/renames that path out from under you.

Created so that Go programs can be used with the standard Linux `logrotate`: writes following a removal / rename will occur in a newly created file rather than the previously opened filehandle.

As the current implementation relies on remembering and checking the current inode of the desired file, this code will not work or compile on Windows.

### Example

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

### TODO

 * Use `exp/fsnotify` so that we don't need to call `stat` on every write
 * Use `exp/winfsnotify` to allow for Windows support
 * Trap SIGHUP and close/reopen all files
 * Force a check every minute so that programs that rarely write can release filehandles

