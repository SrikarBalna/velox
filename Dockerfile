FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc g++ \
    python3 python3-pip \
    nodejs npm \
    openjdk-17-jdk-headless \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && npm install -g typescript

RUN useradd -m runner

WORKDIR /app

COPY run.sh /run.sh
COPY run_single.sh /run_single.sh

RUN chmod +x /run.sh /run_single.sh

USER runner

ENTRYPOINT ["/run.sh"]