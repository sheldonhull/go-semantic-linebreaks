VERSION 0.6
FROM golang:1.17.6-alpine3.15
WORKDIR /proj

build:
    COPY main.go .
    RUN go build -o ./artifacts/gslb /proj/cmd/go-semantic-linebreaks/go-semantic-linebreaks.go
#    SAVE ARTIFACT build/proj /proj AS LOCAL build/proj

# docker:
#     COPY +build/proj .
#     ENTRYPOINT ["/proj/proj"]
#     SAVE IMAGE proj:latest
