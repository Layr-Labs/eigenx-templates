FROM node:22-alpine

WORKDIR /app

COPY package.json ./
RUN npm install

COPY src/ ./src/
COPY tsconfig.json ./

RUN npm run build

EXPOSE 8080

CMD ["npm", "start"]
