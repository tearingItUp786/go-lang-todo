# Stage 1: Build the front-end assets using Node.js
FROM node:latest as build-stage
# Set the working directory for the front-end build
WORKDIR /app
# Copy package.json and package-lock.json (if available)
COPY package*.json ./
# Install npm dependencies
RUN npm install
# Copy the front-end source code
COPY . .
# Build the front-end assets
RUN npm run build

# Stage 2: Build the Go application
FROM golang:latest
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the Go source code
COPY . .
# Copy the built front-end assets from the build-stage to the Go application directory
COPY --from=build-stage /app/dist /app/dist
# Build the Go app
RUN go build -o main .
# Expose the port the Go app runs on
EXPOSE 13255 
# Command to run the executable
CMD ["./main"]
