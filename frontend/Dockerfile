# ---- Build stage ----
FROM node:20-alpine AS builder

WORKDIR /app

# Same caching trick as the backend: copy only the manifests first,
# so `npm install` (slow) is only re-run when dependencies actually change.
COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# The backend's URL is baked into the static JS bundle at build time
# (Vite replaces import.meta.env.VITE_API_URL with this literal value
# during `npm run build` — env vars can't be "read" at runtime by
# static files the way a server-side app would). We accept it as a
# build argument so docker-compose can control it without editing
# this file.
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

RUN npm run build

# ---- Final stage ----
# nginx:alpine is a minimal web server — just enough to serve static
# files, no Node.js or source code included.
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]