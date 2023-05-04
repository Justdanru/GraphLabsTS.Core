FROM golang:1.20-alpine AS build_stage
COPY . /go/src/glts/
WORKDIR /go/src/glts
RUN go install ./...

FROM alpine AS run_stage
WORKDIR /glts_binary
COPY --from=build_stage /go/bin/glts.core /glts_binary/
COPY ./templates/ /glts_binary/templates/
RUN chmod +x ./glts.core
EXPOSE 8081/tcp
ENTRYPOINT ./glts.core