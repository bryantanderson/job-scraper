# Development build
FROM node:22-alpine as dev

WORKDIR /app

COPY package*.json .

RUN npm install

COPY . .

EXPOSE 8000

CMD ["npm", "run", "dev"]

# Build stage
FROM node:22-alpine AS build

WORKDIR /app

COPY package*.json .

RUN npm ci

COPY . .

RUN npm run build

# Production build
FROM node:22-alpine AS production

WORKDIR /app

COPY package*.json .

RUN npm ci --only=production

COPY --from=build /app/dist ./dist

CMD ["node", "dist/app.js"]

