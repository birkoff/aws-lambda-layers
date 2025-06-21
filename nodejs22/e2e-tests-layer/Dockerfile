# Use the Amazon Linux 2023 base image with Node.js 22
FROM amazon/aws-lambda-nodejs:22

# Set the working directory
WORKDIR /var/task
ARG OUTPUT_FILE=e2e-tests-layer.zip

# Install system dependencies (using dnf instead of yum)
RUN dnf install -y \
    amazon-linux-extras \
    && amazon-linux-extras install epel -y \
    && dnf install -y \
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
    && dnf clean all

# Create the nodejs directory structure required for Lambda layers
RUN mkdir -p nodejs/node_modules

# Set the working directory to the node_modules folder
WORKDIR /var/task/nodejs/node_modules

# Install your dependencies here
RUN npm install \
    aws-sdk \
    axios \
    lodash \
    playwright-core \
    puppeteer-core \
    express \
    @sparticuz/chromium \
    nodemon \
    node-stream-zip \
    form-data \
    cors \
    archiver \
    @aws-sdk/client-eventbridge

# Install global tools and Playwright
RUN npm install -g npx
RUN npx playwright install --with-deps

# Set the working directory back to /var/task
WORKDIR /var/task
RUN cp -Rv /var/task/nodejs/* /var/task/

# Install zip and create the layer package
RUN dnf install -y zip && \
    zip -r /opt/${OUTPUT_FILE} nodejs

CMD ["/bin/bash", "-c", "cp /opt/${OUTPUT_FILE} /var/task/ && chown $(stat -c '%u:%g' /var/task) /var/task/${OUTPUT_FILE}"]

# Copy function code
COPY lambda_function.js ./

CMD ["lambda_function.handler"]