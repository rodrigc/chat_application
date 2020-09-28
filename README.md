## chat_application

### Building

```
make clean
make
```

### Running

To run the server:

```
./chatserver
```

To access the server, you must use a telnet client to access the telnet port on the server.
For example, if you are running the server on port 8080, you would do the following to access the server:

```
telnet localhost 8080
```

Type the following to get a list of commands:

```
/help
```


### Further references

While implementing this program, I looked at a few web sites, including
some of the following:

* https://golang.org/pkg/net 
* https://yourbasic.org/golang/convert-string-to-byte-slice/
* https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html
* https://golang.org/pkg/time/
* https://gobyexample.com/time-formatting-parsing
