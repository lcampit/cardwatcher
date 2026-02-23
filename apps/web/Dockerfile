# Placeholder Dockerfile for web frontend
# This will be populated with the actual frontend build instructions
# when the web application is implemented

# Example for a React/Vue/Angular app (adjust as needed):
# FROM node:20-alpine AS build
# WORKDIR /app
# COPY package*.json ./
# RUN npm ci
# COPY . .
# RUN npm run build
#
# FROM nginx:alpine
# COPY --from=build /app/dist /usr/share/nginx/html
# EXPOSE 80
# CMD ["nginx", "-g", "daemon off;"]

# For now, use a minimal placeholder
FROM alpine:3.20

WORKDIR /app

# Install basic dependencies
RUN apk --no-cache add ca-certificates

# Placeholder for web files
COPY . .

EXPOSE 80

CMD ["sh", "-c", "echo 'Web frontend not yet implemented' && sleep 3600"]
