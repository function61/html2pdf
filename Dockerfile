FROM joonas/alpine:f4fddc471ec2-nodejs

EXPOSE 80

RUN npm install -g express body-parser connect-multiparty prom-client

RUN apk add --update curl

RUN curl https://s3.amazonaws.com/infrastructure-cdn.xs.fi/packages/wkhtmltopdf/wkhtmltopdf-0.12.3-rootfs.tar.gz | tar -C / -zxf -

# stupid nodejs
ENV NODE_PATH=/usr/lib/node_modules

# memory optimization tips from https://blog.heroku.com/node-habits-2016
CMD node --optimize_for_size --gc_interval=1000 /app/src/index.js

COPY src/ /app/src
