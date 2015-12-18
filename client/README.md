# Client
A command-line client module to test the services provided by the ```Pillar Server```.

Internally this client module makes REST calls using http (GET or POST).

At this time, the module can read data from an external source such as a Mongo Database.

**Important**
Before you run the client, you must provide an external data source. Use the following format:

~~~
ENV MONGODB_URL mongodb://<user>:<password>@<host>:<port>/echo
~~~


### Running Client
Use the following command to build and run the Client.

~~~
> cd $GOPATH/src
> go install github.com/coralproject/pillar/client
~~~

