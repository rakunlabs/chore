ARG GO_IMAGE
ARG BASE_IMAGE
ARG FRONTEND_IMAGE

######################### BUILDER FRONTEND
FROM $FRONTEND_IMAGE AS builder-frontend
WORKDIR /workspace

COPY _web _web
ARG NPM_PROXY

RUN cd _web && \
    pnpm run depend && \
    pnpm build

######################### BUILDER BACKEND
FROM $GO_IMAGE AS builder-backend

# DEV-BUILD-EDIT
WORKDIR /

ARG SSH_KEY_CICD_64
ARG GOPRIVATE=gitlab.test.igdcs.com
ARG GOPROXY=https://proxy.golang.org,direct
ARG SKIP_CERTS=N

RUN if [ "${SKIP_CERTS}" != "Y" ]; then \
    git config --global url."https://".insteadOf git:// && \
    git config --global http.sslVerify false && \
    mkdir -p -m 0600 ~/.ssh && \
    ssh-keyscan -H gitlab.test.igdcs.com >> ~/.ssh/known_hosts && \
    echo ${SSH_KEY_CICD_64} | base64 -d > ~/.ssh/id_rsa && \
    chmod 600 ~/.ssh/id_rsa && \
    git clone git\@gitlab.test.igdcs.com:finops/devops/infra-certificates.git \
    ; else mkdir -p /infra-certificates/certs; fi
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
COPY --from=builder-backend /infra-certificates/certs /etc/ssl/certs

# Run the binary
ENTRYPOINT ["/chore"]
