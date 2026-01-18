# Step 1: Modules caching
FROM golang:1.25-alpine3.21 AS modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN go mod download

# Step 2: Builder
FROM golang:1.25-alpine3.21 AS builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch

ENV TZ=Asia/Ho_Chi_Minh

COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app"]

# Example commands and in CI/CD
# docker build -t minhhoccode111/${{ secrets.APP_NAME }}:latest .
# docker push minhhoccode111/${{ secrets.APP_NAME }}:latest

# docker pull minhhoccode111/${{ secrets.APP_NAME }}:latest
# docker stop ${{ secrets.APP_NAME }} || true
# docker rm ${{ secrets.APP_NAME }} || true

# docker run -d \
#   --name ${{ secrets.APP_NAME }} \
#   --env-file <(cat <<EOF
# APP_NAME=${{ secrets.APP_NAME }}
# APP_VERSION=${{ secrets.APP_VERSION }}
# CORS_ALLOW_CREDENTIALS=true
# CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization
# CORS_ALLOW_METHODS=GET,POST,HEAD,PUT,DELETE,PATCH
# CORS_ALLOW_ORIGINS=${{ secrets.CORS_ALLOW_ORIGINS_DEV }}
# HTTP_PORT=${{ secrets.HTTP_PORT_DEV }}
# HTTP_USE_PREFORK_MODE=${{ secrets.HTTP_USE_PREFORK_MODE_DEV }}
# JWT_EXPIRATION=${{ secrets.JWT_EXPIRATION }}
# JWT_ISSUER=${{ secrets.JWT_ISSUER_DEV }}
# JWT_SECRET=${{ secrets.JWT_SECRET_DEV }}
# LOG_LEVEL=${{ secrets.LOG_LEVEL_DEV }}
# METRICS_ENABLED=${{ secrets.METRICS_ENABLED_DEV }}
# PG_POOL_MAX=${{ secrets.PG_POOL_MAX_DEV }}
# PG_URL=${{ secrets.PG_URL_DEV }}
# SWAGGER_ENABLED=${{ secrets.SWAGGER_ENABLED_DEV }}
# EOF
# ) \
#   -p ${{ secrets.HTTP_PORT_DEV }}:${{ secrets.HTTP_PORT_DEV }} \
#   minhhoccode111/realworld-fiber-clean:latest
