# We are Refactoring Monolith to Microservices

1. **Refactor the API:**
The monilith API code currently contains logic for both /users and /feed endpoints. I decompose the API code so that we can have two separate projects that can be run independent of one another. We ended up having a different prpkect for the /user called  udagram-api-user and /feed called udagram-api-feed. You may find yourself copying a lot of duplicate code into the separate projects -- this is expected! For now, focus on breaking apart the monolith and we can focus on cleaning up the application code afterwards.
There are different microservices architecture desgins you can explore on your own

2. **Logging**
Use logs to capture metrics. This can help us with debugging. Make sure you implement best practices for logging like assigning a process id to each request. You can pass this id to the request context. This way you can trace all the process that is related to a request through out you application. You can also use exports your logs to external services like datadog, where they expect your logs in a particular format so that they can paerse it for you.

3. **Implement Reverse Proxy to Direct Backend Requests:**
- A reversproxy streamlines the experience so that a consumer only cares about communication with the reverseproxy and not the services that it points to. Our front-end communicates with the reverseproxy for our backend APIs. It makes it easy for our front-end to integrate. If we create more microservices at a later time, our front-end can have less overhead to handle new request. 
- The front-end treats every microservice as a single api under the /api endpoint. The front end dont care what is behine the reverseproxy. The reverseProxy handles routing and mapping to the appropriate api behind the scenes. Nginx is a web server that can be used as a reverse proxy. Our nginx webserver listens for http request comming in at port 8080 on our local machine then forwards requests on behalf of the client to the appropriate microservice based on the pendpoint path and appears to the client as the origin of the responses
- ReverseProxy is one way of Securing your backened services

4.  **Containerize the Code:**
- Start with creating Dockerfiles for the frontend and backend applications. Each project should have its own Dockerfile. Then after we used docker-compose to test the ochestration of the reverseproxy for thesame way it will perform at kubernetes. We used docker-compose to start all the containers for all for images at thesame time. Everything was working fine, then we can push all the images to dockerhub with docker-compose
- We build built and ran our container images using docker
- We store our docker images with a container registry like DockerHub

5. **Build CI/CD Pipeline with [Travis CI](https://docs.travis-ci.com/user/for-beginners/):**
- After setting up your GitHub account to integrate with Travis CI, set up a GitHub repository with a .travis.yml file for a build pipeline to be generated. Once you have the travis.yml in a repo and have your travis account integrated to your github account, travis will look through all your gitbub repositories and when it detects any repo with .travis.yml file, it will recognize that project as something it will set up in the trevis dashboard. Remember to set the environmental varibales for your travis build process of the github repo. Be careful not to echo you secret environmental varibales; look at best practises [here](https://docs.travis-ci.com/user/best-practices-security/)  Looking at the commands we are running in our travis, we used docker-compose. This means any change we made to one microservice in this project repo and push to commit and push to github, the whole services will be build again. This necessary may not be a bad thing knowing that docker do cache your changes for build process to have shorter and shorter time but its something to consider.
- Travis helps us take care of the CI portion of the CICD pilpeline. We integrate github as part of CI/CD and we automate testing using CI
- Other Alternatives to CI are Jenkins

6. **Implement a Health Endpoint:**
This help us to determine when a pod is not healthy. When a port is not healthy, k8s will try to terminate it and generate another


7. **Deploying MicroServices with K8s:**
At this point, you should have a cluster created following these [steps](https://blog.juadel.com/2020/05/15/create-a-kubernetes-cluster-in-amazon-eks-using-a-reverse-proxy/) and the k8s ConfigMap and Secrets environmental varibales applied to your cluster from the kubectl CLI. This way when you deploy your pods, they can use the enviromental variable right away. Two things to consider, 
- A deployment yaml file: This specifies how we deploy an application. In practise, this is similar to how we use docker-compose file to set up docker containers that will run our pods. He explicitly specify in the yml file that this is a deployment file right at the second line. We define it in the "Kind" key of value "deployment"
- A service yml file: This specifies how we configure the service that exposes our application to our consumers. We explicitly define that this is a service file as well
- Look at this article(https://cloud.google.com/blog/products/containers-kubernetes/kubernetes-best-practices-setting-up-health-checks-with-readiness-and-liveness-probes) for best practices on health checks, readiness and liveness probs for your deployment
- Alternative Deployment strategies are [AWS ECS](https://aws.amazon.com/ecs/?whats-new-cards.sort-by=item.additionalFields.postDateTime&whats-new-cards.sort-order=desc&ecs-blogs.sort-by=item.additionalFields.createdDate&ecs-blogs.sort-order=desc), [AWS Fargate](https://aws.amazon.com/fargate/?whats-new-cards.sort-by=item.additionalFields.postDateTime&whats-new-cards.sort-order=desc&fargate-blogs.sort-by=item.additionalFields.createdDate&fargate-blogs.sort-order=desc), [AWS App runner](https://aws.amazon.com/apprunner/)

8. **Secure the API:**
- we want the APIs to be consumed by the frontend web application. So, we set up ingress rules so that only web requests from our web application can make successful API requests.
- Kubernetes Ingress(Inbound web traffic) and Egress(Outbound web traffic) helps us set our least priviledged access for accessing the Kubernetes services like pods thesame way AWS Security Groups helps us set our least priviledged access for AWs resources like S3 and RDS
- Kubernetes natively has the ability to configure Ingress and Egress traffic. This concept is similar to the same way you specify Inbound and Outbound traffic for a security group which you assign to an AWS resources. The way you define these traffics by their Ip address and port number, it is thesame way you define them in K8s. We specify thins in a yaml file called NetworkPolicy and apply it to our cluster
- [https://thenewstack.io/kubernetes-ingress-for-beginners/](https://thenewstack.io/kubernetes-ingress-for-beginners/)
- [https://kubernetes.io/docs/concepts/services-networking/ingress/](https://kubernetes.io/docs/concepts/services-networking/ingress/)

9. **Scaling:**
- Horizontal Scaling is quite straight forward to set up. Run `kubectl autoscale deployment <SERVICE_NAME> --cpu-percent=<CPU_PERCENTAGE_IN_INTEGER> --min <MIN_REPLICAS> --max=<MAX_REPLICAS>
To confirm that autoscale is configured for the desired service run `kubectl get hpa`
- By Setting Horizontal Scaling we are allowing K8s to detects when we are running low on cpu and we want more pods to address this
- Horizontal Scaling is more cost effective than vertical scaling and vertical scallling is not compactible with reverseproxies anyways


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


# [Kubernetes](https://kubernetes.io/docs/tasks/configure-pod-container/) on AWS
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
- See Kubectl Cheatsheet [here](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)


**Managing Secrets in k8s**
- You can store environmental vairbles as open using the k8s ConfigMap or as k8s Secret where you have to encode the value in base64 so that human can read the real value.
- If you want to convert a text to base 64 you can use this [link](https://base64.guru/converter/encode/file) for file or this [link](https://base64.guru/converter/encode/text) for text or do them from commandline. Search for this on stackoverflow :)
[https://kubernetes.io/docs/tasks/configmap-secret/](https://kubernetes.io/docs/tasks/configmap-secret/)

**Debugging a Pod**
- Common reasons your Pod can be failing is due to envireonmental variables or secrets not passed properly to the pod. 
- Make sure you apply the service and deployment file for the pod. If there is no service connected to the pod, the pod will not work. run kubctl get pods or kubectl get services to see all pods and services. 
- One thing you want to do is check the logs of the pod! You run `kubectl logs <unique_pod_name>` This way you can see the error
- Another thing you can check is the environmental variables that are baked or exposed into the pod. This way you can see which environmental variable is not initialised properly. See this [link](https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/) on how you do this
- **FYI:** Sometimes when you update an environmental variable value in the Configmap file , apply the file and the pods running can still have the old value(s). To get the pod have the new env value, you will have to simulate restarting the pod. One way you can do is manually deleting the pod by running `kubectl delete pod <pod_name>` and when you delete the pod, k8s will create a new one because by default, k8s wants to keep the number of replicas you specified in the deployment file for that pod intact. Another way you can simulate creaing a new pod is to edit the deployment file a bit. The slightest change to that will will make k8s delete the old version and create a new pod with the new deployment settings!


**Interacting with your Pods**
- To access the pods endpoints, you can connect into the pod and curl the endpoint (you ideally do this when you want to do some tests) or you have a reverseproxy set up to route requests to the k8s service which will then can call the pod necessary!
- When we connect to the reverseproxy service, it should direct our request to the appropriate pod that we want to access. 
- If you want to expose the reverseproxy to the consumers directly so that that can access your services apis through the reverseproxy, you will have to assign a loadBalancer to the reverseproxy. This standard way to of doing things in AWS Cloud Computing. The LoadBalancer will give you an externalIp link for the reverseproxy that consumers can invoke. The reversproxy interal IP address can only be accessed by pods inside the cluster!






If you want to see how to create eks cluster and node group from AWS CLI, but you will need [eksctl cli](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html) watch this [youtube video](https://www.youtube.com/watch?v=aGTOVaVXz7k&t=474s) and [this video](https://www.youtube.com/watch?v=p6xDCz00TxU&t=664s).


Setting up your system the first time for k8s on mac. Run these commands insequence
- `brew install kubectl` necessary to be able to interact with pods
- `brew tap weaveworks/tap`
- `brew install weaveworks/tap/eksctl` not ncessary but good if you want to create clusters from the console
- `brew install aws-iam-authenticator` to confirm you have aws autheniticator
run `aws sts get-caller-identity` to confirm a user is logged in to aws in console. if no aws user is logged in run `aws configure` and log in a user credentials
- `aws eks --region <region-code> update-kubeconfig --name <cluster_Name>` to connect kubectl against your cluster created on aws. Kubernetes is not created by aws so oyu have to make this connect. This way, any kubectl commands you run here will be applied to your cluster


On your won explore 
- how to deploy private images to k8s
- Approcaches to pass database connection to controllers [here](https://techinscribed.com/different-approaches-to-pass-database-connection-into-controllers-in-golang/)
- nginx configurations [here](https://docs.viblast.com/player/cors/cors-on-nginx)
- Courses on more golang operations for APIS like DB migration, DB transaction, etc [here](https://dev.to/techschoolguru/how-to-create-and-verify-jwt-paseto-token-in-golang-1l5j)