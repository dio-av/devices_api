# syntax=docker/dockerfile:1
FROM golang:1.23 AS build
WORKDIR /cmd
COPY <<EOF /cmd/main.go

RUN go build -o /bin/devices_api ./main.go

FROM scratch
COPY --from=build /bin/devices_api /bin/devices_api
CMD ["/bin/heldevices_api"]