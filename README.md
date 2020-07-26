# bunnycdn-storage-sync

A simple command-line tool to synchronize local file system directories to [BunnyCDN Storages](https://bunnycdn.com/).

# Requirements

- Latest Go version
- Only tested in Unix (macOS)

# Installation

```bash
go get github.com/mtyurt/bunnycdn-storage-sync
```

# Usage
Assuming you have `$GOPATH/bin` in your path:
```bash
BCDN_APIKEY= <apikey> bunnycdn-storage-sync <local_dir> <zone_name>
```

Options:
```
-dry-run   Show the difference and exit
```


# TODO
* Parallel processing of filepath
* Leveled logger usage
* Dockerfile & docker usage
* Build sample github actions to run on push
* Prefix support in storage
