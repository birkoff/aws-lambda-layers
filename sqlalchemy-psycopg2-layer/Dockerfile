FROM --platform=linux/amd64 public.ecr.aws/lambda/python:3.11 AS builder

# Set the working directory
WORKDIR /var/task


ARG OUTPUT_FILE=sqlalchemy-psycopg2-layer.zip


RUN yum install -y \
    python3-devel \
    gcc \
    && yum clean all

# Install Python packages
RUN mkdir /opt/python/ && cd /opt/python/  && \
    pip install --upgrade pip && \
    pip install \
    sqlalchemy \
    psycopg2-binary \
    sqlalchemy-continuum \
    sqlalchemy-serializer \
    sqlalchemy-utils \
    -t .


FROM  --platform=linux/amd64 public.ecr.aws/lambda/python:3.11
ARG OUTPUT_FILE=sqlalchemy-psycopg2-layer.zip

RUN yum install -y openssl zip

COPY --from=builder /opt/python /opt/python

RUN find /opt/python -type d -name '__pycache__' -exec rm -rf {} +

RUN cd /opt && zip -r9 ./${OUTPUT_FILE} python

CMD ["/bin/bash", "-c", "cp /opt/${OUTPUT_FILE} /var/task/ && chown $(stat -c '%u:%g' /var/task) /var/task/${OUTPUT_FILE}"]
