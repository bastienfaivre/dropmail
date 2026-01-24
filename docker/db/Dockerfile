# Minimal container for SQLite data volume
# This container ensures the data directory exists with correct permissions
FROM busybox:stable

# Create data directory
RUN mkdir -p /data && chmod 777 /data

# Volume for SQLite database
VOLUME /data

# Keep container running (optional - can be removed if only used for volume)
CMD ["tail", "-f", "/dev/null"]
