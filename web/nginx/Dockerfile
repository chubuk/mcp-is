FROM nginx:alpine

RUN apk --no-cache add openssl envsubst certbot

COPY nginx.conf /etc/nginx/nginx.conf
COPY snippets/security-headers.conf /etc/nginx/snippets/security-headers.conf
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 80 443

ENTRYPOINT ["/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
