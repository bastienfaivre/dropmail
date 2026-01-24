# Nginx reverse proxy for routing and HTTPS termination
FROM nginx:alpine

# Remove default config
RUN rm /etc/nginx/conf.d/default.conf

# Create SSL directory for certificate mount
RUN mkdir -p /etc/nginx/ssl

# Copy custom config
COPY nginx.conf /etc/nginx/nginx.conf

# Expose HTTP and HTTPS ports
EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
