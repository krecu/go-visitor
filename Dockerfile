FROM alpine:3.6

EXPOSE 8091:8090
EXPOSE 8092:8091

ADD ./databases /databases
ADD ./build /app

CMD ["/app/visitor","-config=/app/"]