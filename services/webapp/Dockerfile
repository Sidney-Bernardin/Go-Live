FROM node:22.20 AS dev

WORKDIR /app

COPY package*.json .
RUN npm i
COPY . .

CMD ["npm", "run", "dev"]

# =====

FROM node:22.20 AS build

WORKDIR /app

COPY package*.json .
RUN npm i --omit=dev
COPY . .

RUN npm run build
