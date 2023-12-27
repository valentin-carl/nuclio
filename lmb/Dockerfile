FROM ubuntu:22.04

WORKDIR /app

# Download all the files!
RUN apt-get update && \
    apt-get install -y wget

# Download and install Go
# TODO this might not be necessary anymore
RUN wget https://go.dev/dl/go1.21.5.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.21.5.linux-arm64.tar.gz && \
    rm go1.21.5.linux-arm64.tar.gz

# Add Go binaries to the system PATH
ENV PATH=$PATH:/usr/local/go/bin

# Verify Go installation
RUN go version

# Install RabbitMQ
RUN apt install -y rabbitmq-server systemctl

# Copy the startup script to set up the broker into the container
# ... and make it executable
COPY setup.sh /app
RUN chmod +x setup.sh

# Expose ports for RabbitMQ
# - 5672 for queue
# - 15672 for management dashboard
EXPOSE 5672
EXPOSE 15672

# Setup rabbitmq and leave the container running
CMD ./setup.sh