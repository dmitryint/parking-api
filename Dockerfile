FROM golang:1.16 as build

WORKDIR /go/src/app
COPY . .
ENV CGO_ENABLED "0"
RUN go build -v -o parking_api  app/*.go


FROM scratch
ENV GIN_MODE "release"
COPY --from=build /go/src/app/parking_api /parking_api
EXPOSE 5000
CMD ["/parking_api"]
