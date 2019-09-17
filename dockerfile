FROM golang:latest

# 安装librdkafka
RUN cd /opt \
    && git clone https://github.com/edenhill/librdkafka.git \
    && cd librdkafka \
    && ./configure --prefix /usr \
    && make \
    && make install