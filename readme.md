# Go todo + htmx + tailwind + lit

This project is a simple todo app using go, htmx, and lit web components. It's the first time I've ever used any of these technologies
but I gotta say, it was a blast 🚀!

## getting started

To run the development environment, follow the steps below:

### prerequisites

Make sure you have the following software installed on your machine:

- docker
- node.js (I recommend volta.sh to use as your node version manager)
- go
- air (go live-reload tool)

### installation

1. install the node dependencies:
   ```shell
   npm install
   ```

### running the development environment

If you want to use google OAuth, please set it up on your end; I have not included
my credientals in my repo 😊

To run the development environment, execute the following commands in order:

1. start the docker containers:
   ```shell
   make docker:start
   ```
2. start the server using `Makefile`:
   ```shell
   make dev
   ```
3. view the app at `localhost:8080`

## explanation

This project is a simple todo app built using go, htmx, and lit web components. it allows users to create, update, and delete tasks in a user-friendly interface. the backend is written in go, providing restful apis for managing the tasks. the frontend utilizes htmx and lit web components to provide a dynamic and responsive user experience.
