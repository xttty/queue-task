FROM golang:1.14

# 安装librdkafka
RUN cd /opt \
    && git clone https://github.com/edenhill/librdkafka.git \
    && cd librdkafka \
    && ./configure --prefix /usr \
    && make \
    && make install \
    && rm -rf /opt/librdkafka