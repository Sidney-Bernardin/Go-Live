FROM nginx:latest AS builder

# Download and install basic dependensies.
RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  wget \
  git \
  build-essential \
  libpcre3-dev \
  zlib1g-dev \
  libssl-dev \
  libzstd-dev

# Download Nginx.
RUN cd /usr/src && \
    wget https://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz && \
    tar zxf nginx-${NGINX_VERSION}.tar.gz && \
    rm nginx-${NGINX_VERSION}.tar.gz

# Download Nginx modules.
RUN cd /usr/src && \
    git clone --recursive https://github.com/arut/nginx-rtmp-module.git rtmp

# Install Nginx and the modules.
RUN cd /usr/src/nginx-${NGINX_VERSION} && \
    ./configure --with-compat \
    --add-dynamic-module=../rtmp && \
    make && make install

# ============================================================================

FROM nginx:latest

# Use the builder's Nginx modules.
COPY --from=builder /usr/local/nginx/modules/* /etc/nginx/modules/

COPY ./nginx.conf /etc/nginx/nginx.conf
RUN mkdir -p /etc/nginx/logs
