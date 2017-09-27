FROM alpine:latest

RUN apk --no-cache add curl  

HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=1 CMD ["curl", "ya.ru"]