FROM debian:bookworm-slim

ARG TURSO_VERSION=latest
ARG TARGETARCH

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        jq \
        unzip \
        xz-utils \
    && rm -rf /var/lib/apt/lists/*

# Install flyctl
RUN curl -fsSL https://fly.io/install.sh | FLYCTL_INSTALL=/usr/local sh

# Install turso CLI
RUN curl -sSfL https://get.tur.so/install.sh | bash && \
    mv /root/.turso/turso /usr/local/bin/turso && \
    rm -rf /root/.turso

# Verify installations
RUN flyctl version && turso --version

ENTRYPOINT ["/bin/bash"]
