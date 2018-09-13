FROM node:9.6.1

RUN npm install -g serve
CMD serve -s build
EXPOSE 3000

COPY package.json package.json
RUN npm install

COPY . .

RUN npm run build --production
