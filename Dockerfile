FROM golang:1 as build

ADD . /source

WORKDIR /source

RUN make index-replicate

FROM centos

COPY --from=build /source/index-replicate /index-replicate

ENTRYPOINT "/index-replicate"
