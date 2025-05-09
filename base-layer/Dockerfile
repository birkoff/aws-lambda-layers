FROM --platform=linux/amd64 public.ecr.aws/lambda/python:3.11 AS builder

# Set the working directory
WORKDIR /var/task

ARG OUTPUT_FILE=base-layer.zip

RUN yum install -y \
    python3-devel \
    gcc \
    && yum clean all

# Create layer directory and install packages
RUN mkdir /opt/python/ && cd /opt/python/  && \
    pip install --upgrade pip && \
    pip install \
    aws-lambda-powertools==2.14.1 \
    aws-xray-sdk==2.13.0 \
    requests \
    xsdata==24.11 \
    PyJWT==2.9.0 \
    cffi==1.17.1 \
    regex \
    certifi \
    charset-normalizer \
    wrapt \
    oauth2client \
    -t .

FROM  --platform=linux/amd64 public.ecr.aws/lambda/python:3.11
ARG OUTPUT_FILE=base-layer.zip

RUN yum install -y openssl zip


COPY --from=builder /opt/python /opt/python

RUN find /opt/python -type d -name '__pycache__' -exec rm -rf {} +

RUN cd /opt && zip -r9 ./${OUTPUT_FILE} python

CMD ["/bin/bash", "-c", "cp /opt/${OUTPUT_FILE} /var/task/ && chown $(stat -c '%u:%g' /var/task) /var/task/${OUTPUT_FILE}"]
