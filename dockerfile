#BUILDING BLOCK!

#Choose the image in which the app will be builded
FROM golang:1.16-alpine as build

#Apk is the package manager from alpine, so this line just installs git
RUN apk add --no-cache git

#Create a working directory, where all subsequent commands will be run
WORKDIR /src 

RUN go mod init github.com/RubensMetteJr/rabbit-link

#Here we copy the app written in go to the /src directory inside the container
COPY main.go ./

#Now we build the previus copied go app
RUN go build main.go

#RUNTIME BLOCK!

#Choose the image in which our container will run
FROM alpine as runtime

#Copy the binary code generated in the build phase (publisher)
COPY --from=build /src/main /app/main

#Run on cmd the publisher binary as soon as the container is started
CMD [ "/app/main" ]