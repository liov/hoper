# Builder Image
FROM jybl/goprotoc AS source

FROM golang:latest AS upgrade

ADD --from=source /bin/protoc /bin/
ADD --from=source /bin/protoc-gen-enum /bin/
ADD --from=source /bin/protoc-gen-go /bin/
ADD --from=source /bin/protoc-gen-go-grpc /bin/
ADD --from=source /bin/protoc-gen-go-patch /bin/
ADD --from=source /bin/protoc-gen-govalidators /bin/
ADD --from=source /bin/protoc-gen-grpc-gateway /bin/
ADD --from=source /bin/protoc-gen-grpc-gin /bin/
ADD --from=source /bin/protoc-gen-openapiv2 /bin/
ADD --from=source /bin/generate /bin/

FROM golang:latest AS compile

RUN go build -trimpath -o  /work/${output} /work/${sourceFile}

ADD ./protoc /bin/
ADD ./protoc-gen-enum /bin/
ADD ./protoc-gen-go /bin/
ADD ./protoc-gen-go-grpc /bin/
ADD ./protoc-gen-go-patch /bin/
ADD ./protoc-gen-govalidators /bin/
ADD ./protoc-gen-grpc-gateway /bin/
ADD ./protoc-gen-grpc-gin /bin/
ADD ./protoc-gen-openapiv2 /bin/
