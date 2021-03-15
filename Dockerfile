FROM golang:1 as build

COPY . /source

WORKDIR /source

ENV CGO_ENABLED 0
RUN make index-replicate

FROM scratch

COPY --from=build /source/index-replicate /

ENTRYPOINT ["/index-replicate"]
