FROM node:latest AS base

WORKDIR /app

COPY package*.json .
RUN npm i

COPY . .

# ============================================================================

FROM node:latest AS build

ARG rtmp_url
ARG hls_url
ARG users_url
ARG rooms_url

ENV RTMP_URL=$rtmp_url
ENV HLS_URL=$hls_url
ENV USERS_URL=$users_url
ENV ROOMS_URL=$rooms_url

WORKDIR /app
COPY --from=base /app .
RUN npm run build

# ============================================================================

FROM node:latest AS dev

ARG rtmp_url
ARG hls_url
ARG users_url
ARG rooms_url

ENV RTMP_URL=$rtmp_url
ENV HLS_URL=$hls_url
ENV USERS_URL=$users_url
ENV ROOMS_URL=$rooms_url

WORKDIR /app
COPY --from=base /app .
CMD ["npm", "run", "dev"]

# ============================================================================

FROM nginx:1.23.1-alpine AS prod

COPY ./nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
