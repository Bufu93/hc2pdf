FROM alpine:latest

RUN apk --no-cache add ca-certificates bash chromium
WORKDIR /root/

CMD ["./app"]

