ARG GO_IMAGE
ARG BASE_IMAGE

######################### BUILDER
FROM $GO_IMAGE AS builder

# DEV-BUILD-EDIT
WORKDIR /

ARG SSH_KEY_CICD_64
ARG GOPRIVATE=gitlab.test.igdcs.com
ARG GOPROXY=https://proxy.golang.org,direct

RUN git config --global url."https://".insteadOf git:// && \
    git config --global http.sslVerify false && \
    mkdir -p -m 0600 ~/.ssh && \
    ssh-keyscan -H gitlab.test.igdcs.com >> ~/.ssh/known_hosts && \
    echo ${SSH_KEY_CICD_64} | base64 -d > ~/.ssh/id_rsa && \
    chmod 600 ~/.ssh/id_rsa && \
    git clone git\@gitlab.test.igdcs.com:finops/devops/infra-certificates.git
#####

WORKDIR /workspace

## Cache Part
COPY go.* .
RUN go mod download
#####

COPY . .
ARG IMAGE_TAG
RUN ./build.sh --build

######################### IMAGE
FROM $BASE_IMAGE

COPY --from=builder /workspace/_out/linux/chore /chore
COPY --from=builder /infra-certificates/certs /etc/ssl/certs

# Add healthcheck
# HEALTHCHECK --interval=30s --start-period=5s --timeout=2s \
#     CMD /turna api --ping || exit 1

# Run the binary
ENTRYPOINT ["/chore"]
