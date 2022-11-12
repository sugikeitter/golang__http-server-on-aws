FROM public.ecr.aws/bitnami/golang:latest AS build
#Get the hello world package from a GitHub repository
RUN go env -w GOPROXY=direct
# Clear GOPATH for go.mod
ENV GOPATH=
# cache dependencies
ADD go.mod go.sum ./
RUN go mod download
# Build the project and send the output to /bin/HelloWorld
ADD . .
RUN go build -o /bin/go-http

FROM public.ecr.aws/debian/debian:stable-slim
#Copy the build's output binary from the previous build container
COPY --from=build /bin/go-http /bin/go-http
# `docker run -p <host_port>:<container_port> IMAGE_NAME [<container_port>]`
CMD ["80"]
ENTRYPOINT ["/bin/go-http"]