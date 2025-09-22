FROM node:24 AS build

ARG rtmp_url
ARG hls_url
ARG users_url
ARG rooms_url

ENV RTMP_URL=$rtmp_url
ENV HLS_URL=$hls_url
ENV USERS_URL=$users_url
ENV ROOMS_URL=$rooms_url

WORKDIR /app

COPY package*.json .
RUN npm i

COPY . .
RUN npm run build

# ============================================================================

FROM nginx:1.23.1-alpine

COPY ./nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
