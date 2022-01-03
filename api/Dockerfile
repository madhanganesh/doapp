FROM golang:1.17-buster AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY *.go ./
RUN go build -o /doapp

FROM gcr.io/distroless/base-debian10
COPY --from=build /doapp /doapp
ENTRYPOINT [ "/doapp" ]
