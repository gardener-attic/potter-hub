FROM eu.gcr.io/sap-gcp-cp-k8s-stable-hub/3rd/node:12.22.12 as build
WORKDIR /app

COPY package.json yarn.lock tsconfig.json tsconfig.prod.json tslint.json ./
RUN yarn install --frozen-lockfile

COPY public/ public/
COPY src/ src/
RUN yarn run build

# check distroless image
# FROM cgr.dev/chainguard/nginx:1.23.2
# ensure that there the nginx user is used, too
FROM eu.gcr.io/sap-gcp-cp-k8s-stable-hub/3rd/nginx:1.23.2-alpine

WORKDIR /app
# Serve the frontend
COPY --from=build /app/build ./

## add permissions for nginx user
RUN chown -R nginx:nginx /app &&  \
    chmod -R 755 /app && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /etc/nginx/conf.d && \
    touch /var/log/nginx/nginx.pid && \
    chown -R nginx:nginx /var/log/nginx/nginx.pid

# Change the path to the nginx.pid since /var/log is not persistent during builds
RUN sed -i 's/\/var\/run\/nginx.pid/\/var\/log\/nginx\/nginx.pid/g' /etc/nginx/nginx.conf

USER nginx:nginx
# forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log \
    && ln -sf /dev/stderr /var/log/nginx/error.log

# server is executed as nginx user
# configured in /etc/nginx/nginx.config
# prints a warning as container runs with user nginx
#CMD ["nginx", "-g daemon off;"]
