A **go** server for temporary **l**ocal **f**ile sharing.

```
go get -u github.com/tapirliu/golf
```

Run it:
```
lx@localhost:~/project/go-wasm$  golf
Serving folder:
   /home/user001/project/go-wasm
Running at:
   http://localhost:9999
   http://127.0.0.1:9999
   http://192.168.1.111:9999
```

Two options are supported
* `-port=9999`: specify a listening port
* `-b`: open in a browser automatically
