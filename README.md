# Pillar Server
Pillar is a REST based WebService written in `Golang`. It provides the following services:

* Imports external data into Coral data model
* Allows CRUD operation on Coral data model
* Provides simple queries on Coral data model


## Key Points

* Pillar APIs strongly adhere to [REST style](https://en.wikipedia.org/wiki/Representational_state_transfer).

* Pillar APIs only work with [JSON](http://www.json.org/) data.

* Regular `CRUD` API pattern is `/api/*`, where as import API pattern is `/api/import/*`.

* Import related APIs allow you to import data into Coral from an existing Source system. The key to a successful import and tracking lies in `ImportSource`. This structure keeps the original identifiers. Most top-level model e.g. `User` or `Comment` embeds this source data in a field named `Source`.

* We understand that an import process can be challenging, hence all import APIs `upsert` data. By doing so, each time you import it overwrites existing entries.


## Running Pillar
Pillar interacts with **Mongo** database as a data-store and  **RabbitMQ** for messaging. All this information can be passed to Pillar through `environment` variables. For convenience, we have an [example](https://github.com/coralproject/pillar/blob/master/config/dev.cfg.sample) file for you. Make a copy of this and `source` the file before running.

####Build Command
```
> cd $GOPATH/src
> go install github.com/coralproject/pillar/app/pillar/
```

####Run Command
```
> source <path>/myenv.cfg
> $GOPATH/bin/pillar
```

## Using Pillar End-Points

Here is a generic example of how you might use these end-points. See [model](https://github.com/coralproject/pillar/tree/master/pkg/model) for the structure of data to be passed for various APIs.

~~~
> curl -i -H "Accept: application/json" -XPOST -d '  {
    "name" : "IamSam",
    "avatar" : "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
    "status" : "New",
    "source" : {
      "id":"original-id-for-iam-sam"
    },
    "tags" : ["top_commentor", "powerball"]
  }
' http://localhost:8080/api/import/user
~~~

Here is a list of end-points

| Model         | Import                   | CRUD            |
|:------------- |:-------------------------|:----------------|
| User          |/api/import/user          |/api/user        |
| Asset         |/api/import/asset         |/api/asset       |
| Action        |/api/import/action        |/api/action      |
| Comment       |/api/import/comment       |/api/comment     |
| Tag           |None                      |/api/tag         |
| Search        |None                      |/api/search      |


## Install Pillar as a Docker Container
Skip this section if you're not familiar or comfortable with Docker. This section helps you build and run a docker image of the Pillar Server.

### Create a Server Docker Image (Optional)

~~~
> cd $GOPATH/src/github.com/coralproject/pillar
> docker build -t pillar-server:0.1 .
~~~

### Run Pillar as a Container
Find the Docker Image with tag pillar-server:0.1 and run the IMAGE_ID.

You must pass the `environment` variables needed to run Pillar, using the env.list file. See ```config/dev.cfg.sample``` file as an example. 

Now, find the image id for ```pillar-server``` and run using the command below:

~~~
> docker images
REPOSITORY     TAG  IMAGE ID       CREATED         VIRTUAL SIZE
pillar-server  0.1  7b59c4c5efde   6 minutes ago   728.2 MB

> docker run --env-file ./env.list --publish 8080:8080 7b59c4c5efde
~~~