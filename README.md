# Acache
## Is a simple service for working with in-memory storage

### Simple test

### **TERM1**
```bash
$ docker run -p 127.0.0.1:11211:11211/tcp -d --name memcached1 memcached:latest

$ go run cmd/acache/main.go 
```

### **TERM2**

```bash
$ make client_run
```
