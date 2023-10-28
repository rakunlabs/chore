ARG ALPINE
# alpine:3.18.4

FROM $ALPINE

RUN apk --no-cache --no-progress add tzdata ca-certificates

ENTRYPOINT [ "/chore" ]
COPY chore /
