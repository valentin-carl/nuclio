FROM ubuntu:22.04

WORKDIR /app

# Download all the files!
RUN apt-get update && \
    apt-get install -y wget htop rabbitmq-server systemctl

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

# Now:
# docker build -t local-message-bus .
# docker run -d -it -p 5672:5672 -p 15672:15672 --name lmb local-message-bus
# docker exec -it lmb bash