# Use Ubuntu as base image
FROM ubuntu:20.04

# Update package list and install required packages including perf tools
RUN apt-get update && apt-get install -y \
    linux-tools-common \
    linux-tools-generic \
    bash \
    && rm -rf /var/lib/apt/lists/*

# Default command to run when starting the container
CMD ["/bin/bash"]
