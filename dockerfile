# BUILDING BLOCK!

# Choose the image in which the app will be built
FROM golang:1.16-alpine as build

# Set the working directory inside the container
WORKDIR /src

# Copy the go.mod file to the container
COPY go.mod ./
COPY go.sum ./

# Copy the rest of the project to the container
COPY . .

# Build the Go application
RUN go build main.go 

# RUNTIME BLOCK!

# Choose the image in which our container will run
FROM alpine as runtime

# Copy the binary code generated in the build phase (publisher)
COPY --from=build /src/main /app/main

# Run the publisher binary as soon as the container is started
CMD ["/app/main"]