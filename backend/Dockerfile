# Node.js/Express backend Dockerfile
FROM node:18

# Install dumb-init for proper signal handling and curl for health checks
# Also install build tools for native dependencies (argon2, sqlite3)
RUN apt-get update && apt-get install -y dumb-init curl python3 make g++ && rm -rf /var/lib/apt/lists/*

# Create app directory and user
RUN groupadd -r nodejs && useradd -r -g nodejs nodejs
WORKDIR /app
RUN chown nodejs:nodejs /app

# Copy package files
COPY --chown=nodejs:nodejs package*.json ./

# Install dependencies
RUN npm ci --silent

# Copy application code
COPY --chown=nodejs:nodejs . .

# Switch to non-root user
USER nodejs

# Set NODE_ENV
ENV NODE_ENV=production

# Expose port
EXPOSE 5007

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:5007/health || exit 1


# Start the application
CMD ["node", "src/server.js"]