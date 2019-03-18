FROM alpine

RUN mkdir /app
RUN apk update \
	&& apk upgrade \
	&& apk add --no-cache \
	ca-certificates \
	&& update-ca-certificates 2>/dev/null || true

ENV SOURCES /app/

COPY freez-app-rest ${SOURCES}
COPY .env ${SOURCES}

WORKDIR ${SOURCES}

CMD ${SOURCES}freez-app-rest

EXPOSE 8080