FROM ubuntu:latest

RUN apt-get update && apt-get install -y curl bash postgresql-client

# Install golang migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz \
    | tar -xvz && \
    mv migrate /usr/local/bin/migrate

COPY ./deployments/docker/migrations/migrations.sh /migrations.sh
COPY ./sql/migrations /migrations
RUN chmod +x /migrations.sh

ENTRYPOINT [ "/bin/bash", "/migrations.sh" ]
