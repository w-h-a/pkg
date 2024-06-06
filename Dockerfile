FROM golang AS build
WORKDIR /go/src/runtime
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/runtime ./

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/runtime /bin/runtime
ENTRYPOINT [ "/bin/runtime" ]