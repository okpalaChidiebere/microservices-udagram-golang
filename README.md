# We are Refactoring Monolith to Microservices

1. **Refactor the API:**
The monilith API code currently contains logic for both /users and /feed endpoints. I decompose the API code so that we can have two separate projects that can be run independent of one another. We ended up having a different prpkect for the /user called  udagram-api-user and /feed called udagram-api-feed. You may find yourself copying a lot of duplicate code into the separate projects -- this is expected! For now, focus on breaking apart the monolith and we can focus on cleaning up the application code afterwards.

2. **Implement Reverse Proxy to Direct Backend Requests:**
A reversproxy streamlines the experience so that a consumer only cares about communication with the reverseproxy and not the services that it points to. Our front-end communicates with the reverseproxy for our backend APIs. It makes it easy for our front-end to integrate. If we create more microservices at a later time, our front-end can have less overhead to handle new request. 
The front-end treats every microservice as a single api under the /api endpoint. The front end dont care what is behine the reverseproxy. The reverseProxy handles routing and mapping to the appropriate api behind the scenes. Nginx is a web server that can be used as a reverse proxy. Our nginx webserver listens for http request comming in at port 8080 on our local machine then forwards requests on behalf of the client to the appropriate microservice based on the pendpoint path and appears to the client as the origin of the responses

3.  **Containerize the Code:**
Start with creating Dockerfiles for the frontend and backend applications. Each project should have its own Dockerfile. Then after we used docker-compose to test the ochestration of the reverseproxy for thesame way it will perform at kubernetes. We used docker-compose to start all the containers for all for images at thesame time. Everything was working fine, then we can push all the images to dockerhub with docker-compose

4. **Build CICD Pipeline with [Travis CI](https://docs.travis-ci.com/user/for-beginners/):**
After setting up your GitHub account to integrate with Travis CI, set up a GitHub repository with a .travis.yml file for a build pipeline to be generated. Once you have the travis.yml in a repo and have your travis account integrated to your github account, travis will look through all your gitbub repositories and when it detects any repo with .travis.yml file, it will recognize that project as something it will set up in the trevis dashboard. Remember to set the environmental varibales for your travis build process of the github repo. Be careful not to echo you secret environmental varibales; look at best practises [here](https://docs.travis-ci.com/user/best-practices-security/)  Looking at the commands we are running in our travis, we used docker-compose. This means any change we made to one microservice in this project repo and push to commit and push to github, the whole services will be build again. This necessary may not be a bad thing knowing that docker do cache your changes for build process to have shorter and shorter time but its something to consider. 




# Docker-Compose tips
- When you push docker images to dockerhub after you had logged in in your console using `docker login`, if the images is not already in your docker hub account, one will be created else it will update an existing one. This is good because you dont have to explicity create an empty image in docker hub just to finally an image to it later which you can do. But nice to know
- You can search for images in dockerhum right from the terminal as long as the docker image is public. Just run `docker search [search term]`. It lists all the docker imagesnames in dockerhub with descriptions. This is why it is good to add descriptions to your docker images
- The images you run in containers locally dont necessary have to be present in your local machine. When you run `docker run [image_name]`, the imagename could be username/whatever. Docker first checks if an image with the imagename exits locally to run that. If does not exist locally, then docker will go to dockerhub, pull down the image then run the image in a container. This good to know.
- We feed environmental variables to the docker container at run time and NOT build time. Look at 
[https://docs.docker.com/compose/environment-variables/](https://docs.docker.com/compose/environment-variables/)
[https://docs.docker.com/compose/compose-file/compose-file-v3/#env_file](https://docs.docker.com/compose/compose-file/compose-file-v3/#env_file)
The environment file called .env.prod or .env.dev looks like this
```env
POSTGRES_USERNAME=somevalue
POSTGRES_PASSWORD=fancypassword
```
Using the .env file allows you to take a write-once-use-often approach to configuring your containers. Although you might not use the exact same variables for various containers, it allows you to create a single .env file and then easily edit the values, so it can be repurposed for other containers. This also makes for easier writing of docker-compose.yml files, as you're not having to hard-code all environment variables. 
- We use `docker-compose build --parallel` to build all our images in parallel
- We use `docker-compose --env-file ./config/.env.dev up` to start containers for all images at once
Do not commit your env file to github
- All other docker-compose commands [here](https://docs.docker.com/compose/reference/)
- All compose config options for verison 3+ [here](https://docs.docker.com/compose/compose-file/compose-file-v3/)


# Kubernetes on AWS
- **Create an EKS on AWS console:** You from the console or from aws CLI. Important thing to note, the IAM user that created the EKS cluster om AWS (from console or cli) will be thesame user that will create the nodegroup from the CLI as well. It cannot be the root user creating the EKS from console, then in CLI, you have another user creating the nodegroup. Run `aws sts get-caller-identity` to see who is logged in at the cli and you will want to use that user to create the cluster and node group. The user must have admin access in general or at the very list EKS admin access.
- **We create Node Group in AWS console:** You create the node group after you have created the cluster. You can have many nodegroup inside a cluster. Nodegroups are used to run the pods that your cluster will be serving. Without nodegroup, you will not be able to deploy your docker to run through pods.


**Creating an EKS Cluster on AWS console**
- Create cluster in EKS
- Create and specify role for Kubernetes cluster. If you dont have one, this [link](https://docs.aws.amazon.com/eks/latest/userguide/service_IAM_role.html#create-service-role) helps you crate one
- Enable public access. To have less complications accessing this cluster
- Leave the rest of the setting as default. More guide [here](https://docs.aws.amazon.com/eks/latest/userguide/create-cluster.html) It takes 5 to 15 mins to create your cluster


**Creating a Node Group on AWS console** 
- Add Node Group in the newly-created cluster
- Create and specify role for IAM role for [node group](https://docs.aws.amazon.com/eks/latest/userguide/create-node-role.html#create-worker-node-role)
- Create and specify SSH key for node group. If you dont have SSH, this [link](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html#having-ec2-create-your-key-pair) will help you create one
- Set instance type to t3.micro for cost-savings as we learn how to use Kubernetes. This affects your pods for compute power
- Specify desired number of nodes. His Affects how your pods scales horixontally


**Interacting With Your Cluster** 
- Make sure you have AWS CLI already installed
- Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/)
- Set up [aws-iam-authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html). To login a user run `aws configure` and enter user credentials. To confirm the user that you have locally login into your awscli run `aws sts get-caller-identity` in terminal
- Set up kubeconfig. Look at this [link](https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html) to what commands to run. At this point you should have connected cluster on aws to kubectl commands. kubectl is not created by aws this is why you have to connect kubectl and your cluster so that when you run kubectl commands, you are runningit against the cluster you created!
- You can now start to deploy your yaml files to have your pods available in your cluster. You can add files that will configure open or secret(usually encoded in base64 so that human can't read it at glance) environment variables for your pods,  just run `kubectl apply -f <fileName.yml>`


**Managing Secrets in k8s**
- You can store environmental vairbles as open using the k8s ConfigMap or as k8s Secret where you have to encode the value in base64 so that human can read the real value.
- If you want to convert a text to base 64 you can use this [link](https://base64.guru/converter/encode/file) for file or this [link](https://base64.guru/converter/encode/text) for text or do them from commandline. Search for this on stackoverflow :)
[https://kubernetes.io/docs/tasks/configmap-secret/](https://kubernetes.io/docs/tasks/configmap-secret/)





If you want to see how to create eks cluster and node group from AWS CLI, but you will need [eksctl cli](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html)
watch this [youtube video](https://www.youtube.com/watch?v=aGTOVaVXz7k&t=474s) and [this video](https://www.youtube.com/watch?v=p6xDCz00TxU&t=664s).


Setting up your system the first time for k8s on mac. Run these commands insequence
- `brew install kubectl` necessary to be able to interact with pods
- `brew tap weaveworks/tap`
- `brew install weaveworks/tap/eksctl` not ncessary but good if you want to create clusters from the console
- `brew install aws-iam-authenticator` to confirm you have aws autheniticator
run `aws sts get-caller-identity` to confirm a user is logged in to aws in console. if no aws user is logged in run `aws configure` and log in a user credentials
- `aws eks --region <region-code. update-kubeconfig --name <cluster_Name>`
