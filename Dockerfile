ARG GO_IMAGE
ARG BASE_IMAGE
ARG FRONTEND_IMAGE=rytsh/frontend-pnpm:v0.0.3

######################### BUILDER FRONTEND
FROM $FRONTEND_IMAGE AS builder-frontend
WORKDIR /workspace

COPY _web _web
COPY docs docs

ARG NPM_PROXY
ARG IMAGE_TAG

RUN cd _web && \
    pnpm build-front

######################### BUILDER BACKEND
FROM $GO_IMAGE AS builder-backend

# DEV-BUILD-EDIT
WORKDIR /

ARG GOPRIVATE=gitlab.test.igdcs.com
ARG GOPROXY=https://proxy.golang.org,direct

## Add ca-certificates
RUN apk add --no-cache \
    ca-certificates git bash

# git configurations
RUN git config --global url."https://".insteadOf git:// && \
    git config --global http.sslVerify false
#####

WORKDIR /workspace

## Cache Part
COPY go.* .
RUN go mod download
#####

COPY . .
# Copy output of the frontend
COPY --from=builder-frontend /workspace/_web/dist /workspace/_web/dist
ARG IMAGE_TAG
RUN ./build.sh --build

######################### IMAGE
FROM $BASE_IMAGE

COPY --from=builder-backend /workspace/_out/linux/chore /chore
COPY --from=builder-backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY infra-certificates/certs /etc/ssl/certs

# Run the binary
ENTRYPOINT ["/chore"]
