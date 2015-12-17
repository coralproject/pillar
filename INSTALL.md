# Installation

### Clone Repository

```
git clone git@github.com:CoralProject/pillar.git
```

### Configure

```
export MONGODB_URL=username:password@host/database
```

### Compile

```
go build
```
### Run the server

```
./server
```

### Endpoints

* http://127.0.0.1:8080/api/import/asset
* http://127.0.0.1:8080/api/import/user
* http://127.0.0.1:8080/api/import/comment

### Run the Client

To test that your local service is working correctly you can run the client.

```
cd client
go build
./client
```
