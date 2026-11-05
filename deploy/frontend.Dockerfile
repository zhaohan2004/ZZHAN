FROM node:24-alpine AS user-build
WORKDIR /src
COPY ZZHAN_web/user/package*.json ./
RUN npm ci
COPY ZZHAN_web/user/ ./
ENV VITE_API_BASE=/api/v1
RUN npm run build

FROM node:24-alpine AS admin-build
WORKDIR /src
COPY ZZHAN_web/admin/package*.json ./
RUN npm ci
COPY ZZHAN_web/admin/ ./
ENV VITE_API_BASE=/api/v1
RUN npm run build

FROM nginx:stable-alpine
COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=user-build /src/dist /usr/share/nginx/html/user
COPY --from=admin-build /src/dist /usr/share/nginx/html/admin
EXPOSE 3000 4000
HEALTHCHECK --interval=15s --timeout=5s --retries=4 \
  CMD wget -q -O /dev/null http://127.0.0.1:3000/ && wget -q -O /dev/null http://127.0.0.1:4000/ || exit 1
