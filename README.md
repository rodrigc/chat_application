## chat_application

### Building

Need Go 1.13+ with module mode turned on by default.
I wrote this with Go 1.15.

To build the application type:

```
make clean
make
```

This should create a binary: *./chatserver*.


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

To access the REST endpoint, you can do an HTTP GET  on the following endpoints:


```
http://localhost:8081/room
http://localhost:8081/user
```
### TODO

I haven't finished these features yet.
* If a user tries to use the /nick command to choose an existing name, the operation should fail.
* Need to implement message logging
* Need to finish the REST API endpoints

### Further references

While implementing this program, I looked at a few web sites, including
some of the following:

* https://golang.org/pkg/net 
* https://yourbasic.org/golang/convert-string-to-byte-slice/
* https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html
* https://golang.org/pkg/time/
* https://gobyexample.com/time-formatting-parsing
