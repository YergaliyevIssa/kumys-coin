# Start from a small base image
FROM alpine:3.14

# Install ca-certificates in case your app makes HTTPS calls
RUN apk add --no-cache ca-certificates

# Set the working directory in the container
WORKDIR /app

# Copy the pre-built binary from your host to the container
COPY bin/app .

# Copy the configuration file
COPY config/env.sh /app/config/

# Make sure the binary is executable
RUN chmod +x /app/app

# Source the environment file and run the binary
CMD ["/bin/sh", "-c", "source /app/config/env.sh && /app/app"]
