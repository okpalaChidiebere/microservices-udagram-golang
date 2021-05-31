# We are Refactoring Monolith to Microservices

1. # Refactor the API
The monilith API code currently contains logic for both /users and /feed endpoints. I decompose the API code so that we can have two separate projects that can be run independent of one another. We ended up having a different prpkect for the /user called  udagram-api-user and /feed called udagram-api-feed. You may find yourself copying a lot of duplicate code into the separate projects -- this is expected! For now, focus on breaking apart the monolith and we can focus on cleaning up the application code afterwards.

2. # Implement Reverse Proxy to Direct Backend Requests
A reversproxy streamlines the experience so that a consumer only cares about communication with the reverseproxy and not the services that it points to. Our front-end communicates with the reverseproxy for our backend APIs. It makes it easy for our front-end to integrate. If we create more microservices at a later time, our front-end can have less overhead to handle new request. 
The front-end treats every microservice as a single api under the /api endpoint. The front end dont care what is behine the reverseproxy. The reverseProxy handles routing and mapping to the appropriate api behind the scenes. Nginx is a web server that can be used as a reverse proxy. Our nginx webserver listens for http request comming in at port 8080 on our local machine then forwards requests on behalf of the client to the appropriate microservice based on the pendpoint path and appears to the client as the origin of the responses

3.  # Containerize the Code
Start with creating Dockerfiles for the frontend and backend applications. Each project should have its own Dockerfile. Then after we used docker-compose to test the ochestration of the reverseproxy for thesame way it will perform at kubernetes. We used docker-compose to start all the containers for all for images at thesame time. Everything was working fine, then we can push all the images to dockerhub with docker-compose




# Docker-Compose tips
- When you push docker images to dockerhub after you had logged in in your console using `docker login`, if the images is not already in your docker hub account, one will be created else it will update an existing one. This is good because you dont have to explicity create an empty image in docker hub just to finally an image to it later which you can do. But nice to know
- You can search for images in dockerhum right from the terminal as long as the docker image is public. Just run `docker search [search term]`. It lists all the docker imagesnames in dockerhub with descriptions. This is why it is good to add descriptions to your docker images
- The images you run in containers locally dont necessary have to be present in your local machine. When you run `docker run [image_name]`, the imagename could be username/whatever. Docker first checks if an image with the imagename exits locally to run that. If does not exist locally, then docker will go to dockerhub, pull down the image then run the image in a container. This good to know.
- We feed environmental variables to the docker container at run time and NOT build time. Look at 
(https://docs.docker.com/compose/environment-variables/)[https://docs.docker.com/compose/environment-variables/]
(https://docs.docker.com/compose/compose-file/compose-file-v3/#env_file)[https://docs.docker.com/compose/compose-file/compose-file-v3/#env_file]
The environment file called .env.prod or .env.dev looks like this
```env
POSTGRES_USERNAME=somevalue
POSTGRES_PASSWORD=fancypassword
```
Using the .env file allows you to take a write-once-use-often approach to configuring your containers. Although you might not use the exact same variables for various containers, it allows you to create a single .env file and then easily edit the values, so it can be repurposed for other containers. This also makes for easier writing of docker-compose.yml files, as you're not having to hard-code all environment variables. 
- We use `docker-compose build --parallel` to build all our images in parallel
- We use `docker-compose --env-file ./config/.env.dev up` to start containers for all images at once
Do not commit your env file to github
- All other docker-compose commands (here)[https://docs.docker.com/compose/reference/]
- All compose config options for verison 3+ (here)[https://docs.docker.com/compose/compose-file/compose-file-v3/]