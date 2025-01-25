# Build
FROM docker.io/golang:1.20-bullseye AS build

WORKDIR /project

COPY go.* ./

RUN go mod download

COPY ./ ./

RUN go build -o ./bin/app ./cmd/app

# Distroless
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /project/bin/app /project/bin/app

USER nonroot:nonroot

CMD [ "/project/bin/app" ]
