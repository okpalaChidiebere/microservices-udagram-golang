#https://medium.com/@knoldus/deployment-with-docker-in-ionic-2797ebbe65ec

# Use NodeJS base image
FROM node:latest AS build

# Create app directory in Docker
WORKDIR /usr/src/app

# Install app dependencies by copying
# package.json and package-lock.json
COPY package*.json ./

# Install dependencies in Docker
RUN npm ci

# Install Ionic
RUN npm install -g ionic

# Copy app from local environment into the Docker image
COPY . .

# Build the Ionic application
RUN ionic build

# Set the application's port number. We expose the nginx for our ionic frontend will be listening for requests on port 80
EXPOSE 80

# Use Nginx base image
FROM nginx:alpine

#we are removing everything present in
RUN rm -rf /usr/share/nginx/html/*

# Copy build artifacts to Nginx html folder for static hosting
COPY --from=build /usr/src/app/www /usr/share/nginx/html

# Copy custom Nginx config, to fix an annoying bug :(
COPY nginx.conf /etc/nginx/nginx.conf