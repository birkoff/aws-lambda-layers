FROM --platform=linux/amd64 public.ecr.aws/lambda/python:3.13 AS builder

ENV ODBCINI=/opt/odbc.ini
ENV LD_LIBRARY_PATH=/opt/microsoft/msodbcsql18/lib64:/opt/lib
ENV ODBCSYSINI=/opt
ARG UNIXODBC_VERSION=2.3.12

RUN dnf install -y gzip tar openssl-devel gcc gcc-c++ make automake kernel-devel

RUN curl ftp://ftp.unixodbc.org/pub/unixODBC/unixODBC-${UNIXODBC_VERSION}.tar.gz -O \
    && tar xzvf unixODBC-${UNIXODBC_VERSION}.tar.gz \
    && cd unixODBC-${UNIXODBC_VERSION} \
    && ./configure --sysconfdir=/opt --disable-gui --disable-drivers --enable-iconv --with-iconv-char-enc=UTF8 --with-iconv-ucode-enc=UTF16LE --prefix=/opt \
    && make \
    && make install

RUN curl https://packages.microsoft.com/config/rhel/9/prod.repo > /etc/yum.repos.d/mssql-release.repo
RUN dnf install -y e2fsprogs fuse-libs libss
RUN ACCEPT_EULA=Y dnf install -y msodbcsql18

ENV CFLAGS="-I/opt/include"
ENV LDFLAGS="-L/opt/lib"

RUN mkdir /opt/python/ && cd /opt/python/ && pip install pyodbc -t .

FROM  --platform=linux/amd64 public.ecr.aws/lambda/python:3.13

RUN dnf install -y openssl tree

COPY --from=builder /opt/python /opt/python
COPY --from=builder /opt/microsoft /opt/microsoft
COPY --from=builder /opt/lib /opt/lib

RUN ls -rla /opt/ && tree /opt/ && ls -rla /opt/python && tree /opt/python