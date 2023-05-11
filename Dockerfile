# Base container alpine in last stable and recommended version
# Alpine is a very light container with the minimum necessary for executing tests
# it can be built and executed very fast because of this
FROM alpine:latest

# Install dependencies with root user
USER root

# Install Go language, task runner and test framework
RUN apk update && apk upgrade \
    && adduser -D tropestogo \
    && apk add --no-cache 'go>1.20' curl cargo git bash \
    && curl https://sh.rustup.rs -sSf | bash -s -- -y \
    && source $HOME/.cargo/env \
    && source $HOME/.cargo/env \
    && rustup update \
    && mkdir -p /app/test \
    && chown tropestogo /app/test \
    # Cargo package manager is necessary for installing Mask task runner
    && cargo install mask --root /tropestogo/.cargo/ \
    # Delete unneeded packages
    && apk del curl perl wget git bash

RUN go version

USER tropestogo

WORKDIR /app/test

COPY . .

CMD ["/tropestogo/.cargo/bin/mask", "test"]