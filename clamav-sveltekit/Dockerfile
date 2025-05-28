FROM node:18-alpine AS build

WORKDIR /app

# Install additional packages for better Node.js support
RUN apk add --no-cache libc6-compat openssl

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the app with experimental options
ENV NODE_OPTIONS="--experimental-vm-modules"
RUN npm run build

# Production stage
FROM node:18-alpine

WORKDIR /app

# Install additional packages for production
RUN apk add --no-cache libc6-compat openssl

# Copy built assets from the build stage
COPY --from=build /app/build ./build
COPY --from=build /app/package.json ./

# Install only production dependencies
RUN npm install --production

# Start the app
EXPOSE 3000
CMD ["node", "build"] 