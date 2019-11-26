FROM postgres:12

USER postgres

RUN pg_createcluster 12 main

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker todoapp &&\
    /etc/init.d/postgresql stop

ADD init.sql /docker-entrypoint-initdb.d/

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/12/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/12/main/postgresql.conf
RUN echo "synchronous_commit = off" >> /etc/postgresql/12/main/postgresql.conf
RUN echo "shared_buffers = 256MB" >> /etc/postgresql/12/main/postgresql.conf
RUN echo "autovacuum = off" >> /etc/postgresql/12/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

EXPOSE 5432
