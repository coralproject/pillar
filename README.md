# Pillar Server
Pillar is a REST based web-service module written in golang. It provides the following services:

* Imports external data into the coral data model
* Allows CRUD operation on coral data model
* Provides simple queries on coral data model

**Important**
You can configure the server using three environment variables. See ```config/dev.cfg``` file as an example. Make changes to the config file as needed and source it as follows

~~~
> source config/dev.cfg
~~~


## End-Points
Server provides the following end-points:

* /api/import/asset
* /api/import/user
* /api/import/comment

Here is a generic example how you might use these end-points:

~~~
> curl -i -H "Accept: application/json" -X POST -d '{"source.id": "original_id", "url": "url of the asset"}' http://localhost:8080/api/import/asset

HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 20 Oct 2015 15:25:12 GMT
Content-Length: 173

{"id":"a5efbb05-6ed7-455e-bc4c-37236614ac14","source.id": "original_id", "url": "url of the asset"}> 
~~~

### /api/import/asset
Imports an ```Asset``` from an external system and the caller must pass a json payload for an ```Asset``` in the following format:

~~~
{
  "source.id" : "42f215a2-066c-11e5-a428-c984eb077d4e",
  "url" : "http://washingtonpost.com/world/national-security/some-nsa-surveillance-powers-set-to-expire-sunday-unless-senate-acts/2015/05/31/42f215a2-066c-11e5-a428-c984eb077d4e_story.html"
}
~~~

### /api/import/user
Imports a ```User``` from an external system and the caller must pass a json payload for a ```User``` in the following format:

~~~
{
  "source.id" : "u6qTe%2BFQ%2BFli6rmbWJ6BEP%2BLRzrUEvutviR1VYa5PdNoGeVxxhJF5A%3D%3D",
  "user_name" : "sazcrin",
  "avatar" : "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
  "status" : "ModeratorApproved"
}
~~~

### /api/import/comment
Imports a ```Comment``` from an external system and the caller must pass a json payload for a ```Comment ``` in the following format:

~~~
{
  "body":"Drinking alcohol isn't an explicit constitutional right.  Better would be these prior restraints on writing an editorial, joining a congregation, or registering to vote.\n\nAll of which will come if these people have their way.",
  "status": "Untouched",
  "source": {
    "id":"f2582294-a4c1-461a-982f-9e63dffbae6a",
    "asset_id":"http://washingtonpost.com/world/national-security/some-nsa-surveillance-powers-set-to-expire-sunday-unless-senate-acts/2015/05/31/42f215a2-066c-11e5-a428-c984eb077d4e_story.html",
    "user_id":"u6qTe%2BFQ%2BFli6rmbWJ6BEP%2BLRzrUEvutviR1VYa5PdNoGeVxxhJF5A%3D%3D"
  },
  "date_created": "2015-11-10T00:00:02.626Z",
  "date_updated": "2015-11-10T00:00:02.626Z"
}
~~~


## Pillar Server as a Docker Container
Skip this section if you're not familiar or comfortable with Docker. This section helps you build and run a docker image of the Pillar Server.


### Create a Server Docker Image (Optional)

~~~
> cd $GOPATH/src/github.com/coralproject/pillar
> docker build -t pillar-server:1.0 .
~~~

### Run Server Container
Find the Docker Image with tag pillar-server:0.1 and run the IMAGE_ID.

You must pass the environment variables needed to run Pillar, using the env.list file. See ```config/dev.cfg``` file as an example. 

~~~
export PILLAR_ADDRESS=:8080
export PILLAR_HOME=/opt/pillar
export MONGODB_URL=mongodb://user:password@host:port/coral
~~~

Find the image id for ```pillar-server``` and run using the command below:

~~~
> docker images
REPOSITORY     TAG  IMAGE ID       CREATED         VIRTUAL SIZE
pillar-server  0.1  7b59c4c5efde   6 minutes ago   728.2 MB

> docker run --env-file ./env.list --publish 8080:8080 7b59c4c5efde
~~~
