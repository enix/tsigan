FROM --platform=$BUILDPLATFORM golang:1.23.2 AS base
WORKDIR /opt/app

FROM --platform=$BUILDPLATFORM cosmtrek/air:v1.61.1 AS air

FROM base AS dev
ARG USER_ID=1000
ARG GROUP_ID=1000
RUN groupadd -g ${GROUP_ID} air \
    && useradd -l -u ${USER_ID} -g air air \
    && install -d -m 0700 -o air -g air /home/air
USER ${USER_ID}:${GROUP_ID}
COPY --from=air /go/bin/air /go/bin/air
CMD [ "/go/bin/air" ]

FROM base AS build
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download -x
COPY internal internal
COPY pkg pkg
COPY cmd cmd
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG GIT_COMMIT
ARG GIT_STATE
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
RUN go build -v -a -buildvcs=false \
    -o /usr/bin/tsigan \
    -tags osusergo,netgo \
    -ldflags "-X \"github.com/enix/tsigan/internal/product.version=${VERSION}\" \
    -X \"github.com/enix/tsigan/internal/product.gitCommit=${GIT_COMMIT}\" \
    -X \"github.com/enix/tsigan/internal/product.gitTreeState=${GIT_STATE}\" \
    -X \"github.com/enix/tsigan/internal/product.buildTime=$(date --iso-8601=seconds)\"" \
    ./cmd/

FROM busybox:1.37.0-glibc AS shell
COPY --from=build /usr/bin/tsigan /tsigan
EXPOSE 5353/udp
EXPOSE 5353/tcp
ENTRYPOINT [ "/tsigan" ]

FROM scratch
COPY --from=build /usr/bin/tsigan /tsigan
EXPOSE 5353/udp
EXPOSE 5353/tcp
ENTRYPOINT [ "/tsigan" ]
