FROM node:20-alpine AS build

WORKDIR /app

ENV NUXT_TELEMETRY_DISABLED=1
ARG NUXT_PUBLIC_API_BASE_URL
ENV NUXT_PUBLIC_API_BASE_URL=${NUXT_PUBLIC_API_BASE_URL:-http://localhost:8080}

COPY . .

RUN if [ ! -f package.json ]; then echo "package.json missing" && exit 1; fi && \
    if [ -f package-lock.json ]; then npm ci; \
    elif [ -f npm-shrinkwrap.json ]; then npm ci; \
    elif [ -f pnpm-lock.yaml ]; then corepack enable && pnpm install --frozen-lockfile; \
    elif [ -f yarn.lock ]; then corepack enable && yarn install --frozen-lockfile; \
    else npm install; \
    fi && \
    npm run build

FROM node:20-alpine AS runtime

WORKDIR /app

ENV NODE_ENV=production
ENV NUXT_TELEMETRY_DISABLED=1

RUN apk add --no-cache curl

COPY --from=build /app/.output ./.output
COPY --from=build /app/package.json ./package.json

EXPOSE 3000

HEALTHCHECK --interval=10s --timeout=5s --retries=3 CMD curl -f http://localhost:3000 || exit 1

CMD ["node", ".output/server/index.mjs"]
