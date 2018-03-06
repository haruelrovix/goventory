# goventory
Inventory REST API built in Go, SQLite, Docker and Kubernetes (hopefully)

### Prerequisite
1. [Go](https://golang.org/)
```sh
$ go version
go version go1.10 darwin/amd64
```
2. [SQLite](https://www.sqlite.org/index.html)
```sh
$ sqlite3 version
SQLite version 3.16.0 2016-11-04 19:09:39
```

### How to Run
1. Clone this repository
```sh
$ git clone https://github.com/haruelrovix/goventory.git && cd goventory
```

2. Execute `goventory.sh`
```sh
$ ./goventory.sh
```

3. If it asks to accept incoming network connections, allow it.
<img src="https://i.imgur.com/FqfijBf.png" alt="Accept incoming network connections" width="30%" />

4. `goventory` listening on port 3000
```sh
Starting goventory on Port :3000. Press Ctrl-C to quit.
```

### Test the API
1. Catatan Nilai Barang: `http://0.0.0.0:3000/api/items`

<img src="https://user-images.githubusercontent.com/17120764/37026626-8a536fd2-2161-11e8-92fd-7c3cd14e025b.png" title="Catatan Nilai Barang" width=500 />

2. Catatan Barang Masuk: `http://0.0.0.0:3000/api/barangmasuk`

<img src="https://user-images.githubusercontent.com/17120764/37026731-dda0d1f2-2161-11e8-90ed-f544d1eeacd3.png" title="Catatan Barang Masuk" width=500 />

3. Catatan Barang Keluar: `http://0.0.0.0:3000/api/barangkeluar`

<img src="https://user-images.githubusercontent.com/17120764/37026797-0f9e5828-2162-11e8-903e-0a7c5b884dd3.png" title="Catatan Barang Keluar" width=500 />

4. Laporan Nilai Barang: `http://0.0.0.0:3000/api/nilaibarang`

<img src="https://user-images.githubusercontent.com/17120764/37026887-552ea532-2162-11e8-8212-2de0d4ff8527.png" title="Laporan Nilai Barang" width=500 />

5. Laporan Penjualan: `http://0.0.0.0:3000/api/penjualan?startdate=2017-12-01&enddate=2017-12-31`

Both `startdate` and `enddate` requires date in `YYYY-MM-DD` format.

<img src="https://user-images.githubusercontent.com/17120764/37026971-97acd50a-2162-11e8-915f-67a1b57c59e3.png" title="Laporan Penjualan" width=500 />

Either way, it throws `400 Bad Request`.

<img src="https://user-images.githubusercontent.com/17120764/37027110-07611672-2163-11e8-9117-4df9334da90f.png" title="Bad Request" width=500 />

### Debugging
VS Code and [Delve](https://github.com/derekparker/delve), a debugger for the Go programming language.

<img src="https://user-images.githubusercontent.com/17120764/37027460-e8303b74-2163-11e8-8701-79d543afaefa.png" alt="Debugging" width=500 />
