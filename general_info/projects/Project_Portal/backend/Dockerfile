FROM busybox:latest

WORKDIR /app

# Copy the locally built binary
COPY main .

# Copy HTML templates
COPY templates/ ./templates/

# Create uploads directory
RUN mkdir -p uploads

# Expose port for information purpose only
EXPOSE 8081

# Run the binary
CMD ["./main"]