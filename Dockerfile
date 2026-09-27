# Подготовка.
FROM golang:alpine AS preper
ARG NAME
WORKDIR /usr/src/${NAME}
# Создание директорий.
RUN mkdir /etc/${NAME}
RUN mkdir /var/log/${NAME}
# Копирование файлов.
COPY . .
COPY ./config/prod.yaml /etc/${NAME}/prod.yaml
RUN go mod download

# Тестирование.
FROM preper AS testing
ARG DB_PASSWORD
RUN go test --tags=unit -v ./...
RUN go test --tags=integration -v ./...
RUN go test --tags=e2e -v ./...

# Сборка.
FROM preper AS building
RUN go build -v -o /usr/local/bin/${NAME} ./cmd/main.go

# Релиз.
FROM alpine AS release
EXPOSE 8000
ARG NAME
# Настройки.
RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
# Создание директорий.
RUN mkdir /var/log/${NAME}
RUN mkdir /etc/${NAME}
# Копирование конфигурации.
COPY ./config/prod.yaml /etc/${NAME}/prod.yaml
# Копирование исполняемого файла.
COPY --from=building /usr/local/bin/${NAME} /usr/local/bin/${NAME}
# Копирование точки входа.
COPY ./docker-entrypoint.sh /usr/local/bin
RUN echo "exec ${NAME}" >> /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT [ "docker-entrypoint.sh" ]
