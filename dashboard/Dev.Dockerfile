FROM docker.wdf.sap.corp:51505/securityapprovedbaseimages/node:12.18.3 as build
WORKDIR /app

COPY package.json yarn.lock tsconfig.json tsconfig.prod.json tslint.json ./
RUN yarn install --frozen-lockfile

COPY public/ public/
COPY src/ src/
RUN yarn run build

FROM docker.wdf.sap.corp:51505/securityapprovedbaseimages/nginx:1.19.2-alpine

WORKDIR /app
# Serve the frontend
COPY --from=build /app/build ./

# # Install and run NGINX
# RUN apt-get -y update && apt-get -y install ca-certificates nginx && update-ca-certificates && \
#     rm -f /etc/nginx/conf.d/default.conf && \ 
#     mkdir /run/nginx && \
#     printf "\ndaemon off;" >> /etc/nginx/nginx.conf 

# FOR DEBUGGING
RUN cat /etc/passwd && \
    ls -la /run && \
    touch /test.txt && \
    chown -R nginx:nginx /test.txt

#END DEBUGGING

## add permissions for nginx user
RUN chown -R nginx:nginx /app &&  \
    chmod -R 755 /app && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /etc/nginx/conf.d && \
    touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/run/nginx.pid

# DEBUG
RUN cat /etc/passwd

USER nginx:nginx

# forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log \
    && ln -sf /dev/stderr /var/log/nginx/error.log

# server is executed as nginx user
# configured in /etc/nginx/nginx.config
# prints a warning as container runs with user nginx
#CMD ["nginx", "-g daemon off;"]
