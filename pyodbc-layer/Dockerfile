FROM --platform=linux/amd64 public.ecr.aws/lambda/python:3.11 AS builder

ENV ODBCINI=/opt/odbc.ini
ENV LD_LIBRARY_PATH=/opt/microsoft/msodbcsql18/lib64:/opt/lib
ENV ODBCSYSINI=/opt
ARG UNIXODBC_VERSION=2.3.12
ARG OUTPUT_FILE=pyodbc-layer.zip

RUN yum install -y gzip tar openssl-devel gcc gcc-c++ make automake kernel-devel

RUN curl ftp://ftp.unixodbc.org/pub/unixODBC/unixODBC-${UNIXODBC_VERSION}.tar.gz -O \
    && tar xzvf unixODBC-${UNIXODBC_VERSION}.tar.gz \
    && cd unixODBC-${UNIXODBC_VERSION} \
    && ./configure --sysconfdir=/opt --disable-gui --disable-drivers --enable-iconv --with-iconv-char-enc=UTF8 --with-iconv-ucode-enc=UTF16LE --prefix=/opt \
    && make \
    && make install

RUN curl https://packages.microsoft.com/config/rhel/9/prod.repo > /etc/yum.repos.d/mssql-release.repo
RUN yum install -y e2fsprogs fuse-libs libss
RUN ACCEPT_EULA=Y yum install -y msodbcsql18

ENV CFLAGS="-I/opt/include"
ENV LDFLAGS="-L/opt/lib"

RUN mkdir /opt/python/ && cd /opt/python/ && pip install pyodbc -t .

FROM  --platform=linux/amd64 public.ecr.aws/lambda/python:3.11

RUN yum install -y openssl zip

COPY --from=builder /opt/python /opt/python
COPY --from=builder /opt/microsoft /opt/microsoft
COPY --from=builder /opt/lib /opt/lib

RUN cd /opt && zip -r9 pyodbc-layer.zip .

CMD ["/bin/bash", "-c", "cp /opt/${OUTPUT_FILE} /var/task/ && chown $(stat -c '%u:%g' /var/task) /var/task/${OUTPUT_FILE}"]
