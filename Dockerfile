FROM alpine

RUN mkdir /app

ENV SOURCES /app/

COPY freez-app-rest ${SOURCES}
COPY .env ${SOURCES}

WORKDIR ${SOURCES}

CMD ${SOURCES}freez-app-rest

EXPOSE 8080