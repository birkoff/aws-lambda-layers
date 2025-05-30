# Use the Amazon Linux 2 base image which matches the Lambda execution environment
FROM amazon/aws-lambda-nodejs:18

# Set the working directory
WORKDIR /var/task
ARG OUTPUT_FILE=e2e-tests-layer.zip

RUN yum install -y \
    amazon-linux-extras \
    && amazon-linux-extras install epel -y \
    && yum install -y \
    gtk3 \
    alsa-lib \
    nss \
    at-spi2-atk \
    cups-libs \
    libdrm \
    libXcomposite \
    libXdamage \
    libXfixes \
    libXrandr \
    libgbm \
    pango \
    atk \
    at-spi2-core \
    gcc-c++ \
    make \
    unzip \
    && yum clean all
# Create the nodejs directory structure required for Lambda layers
RUN mkdir -p nodejs/node_modules

# Set the working directory to the node_modules folder
WORKDIR /var/task/nodejs/node_modules

# Install your dependencies here
# Replace with the packages you want in your layer
RUN npm install aws-sdk axios lodash playwright-core puppeteer-core express @sparticuz/chromium nodemon node-stream-zip form-data cors archiver @aws-sdk/client-eventbridge
RUN npm install -g npx
RUN npx playwright install --with-deps
# Set the working directory back to /layer
WORKDIR /var/task
RUN cp -Rv /var/task/nodejs/* /var/task/
# Zip the nodejs folder (this is the format Lambda expects for layers)
RUN yum install -y zip && \
    zip -r /opt/${OUTPUT_FILE} nodejs

CMD ["/bin/bash", "-c", "cp /opt/${OUTPUT_FILE} /var/task/ && chown $(stat -c '%u:%g' /var/task) /var/task/${OUTPUT_FILE}"]

# Copy function code
COPY lambda_function.js ./

CMD ["lambda_function.handler"]