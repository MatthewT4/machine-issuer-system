FROM postgres:14.8-alpine3.17

COPY migrations/*.sql /docker-entrypoint-initdb.d/

CMD ["postgres" ]
