# ##############################################################################
# Build the front-end
# ##############################################################################
FROM node:24-slim AS frontend

WORKDIR /app
COPY react/astria-web/app/package.json react/astria-web/app/package-lock.json ./
RUN npm install
COPY react/astria-web/app/ ./
RUN npm run build

# ##############################################################################
# Build the server
# ##############################################################################
FROM debian:trixie AS build

RUN apt-get -y update && \
    apt-get install -y --no-install-recommends ca-certificates curl g++ gcc git libc6-dev make unzip && \
    curl -sL https://golang.org/dl/go1.26.1.linux-amd64.tar.gz | tar -C /usr/local -xz

ENV PATH="/usr/local/go/bin:${PATH}"

COPY . /go/src/gitlab.com/olaris/olaris-server
WORKDIR /go/src/gitlab.com/olaris/olaris-server
COPY --from=frontend /app/dist/ ./react/build/

RUN make generate build-local

# ##############################################################################
# Build the release image
# ##############################################################################
FROM debian:trixie AS release

RUN apt-get -y update && \
    apt-get install -y --no-install-recommends ca-certificates ffmpeg sudo && \
    apt-get autoremove && apt-get clean

COPY --from=build /go/src/gitlab.com/olaris/olaris-server/build/olaris /opt/olaris/olaris
COPY ./docker/entrypoint.sh /

# Create a non-root user and group
RUN groupadd --gid 1000 astria \
    && useradd --uid 1000 --gid 1000 --create-home astria

RUN mkdir -p /home/astria/.config/astria && chown astria:astria /home/astria/.config/astria

WORKDIR /home/astria
VOLUME /home/astria/.config/astria
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh", "/opt/olaris/olaris"]
