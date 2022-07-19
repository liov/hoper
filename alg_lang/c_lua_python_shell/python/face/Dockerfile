# Builder Image
FROM python:3.8-slim-buster AS compile


RUN echo 'deb http://mirrors.cloud.tencent.com/debian/ buster main non-free contrib\n\
deb http://mirrors.cloud.tencent.com/debian-security buster/updates main\n\
deb http://mirrors.cloud.tencent.com/debian/ buster-updates main non-free contrib\n\
deb http://mirrors.cloud.tencent.com/debian/ buster-backports main non-free contrib'\
> /etc/apt/sources.list

# Install Dependencies
RUN apt-get -y update && apt-get install -y --fix-missing \
    build-essential \
    cmake \
    gfortran \
    git \
    wget \
    curl \
    graphicsmagick \
    libgraphicsmagick1-dev \
    libatlas-base-dev \
    libavcodec-dev \
    libavformat-dev \
    libgtk2.0-dev \
    libjpeg-dev \
    libwebp-dev \
    liblapack-dev \
    libswscale-dev \
    pkg-config \
    python3-dev \
    python3-numpy \
    software-properties-common \
    zip \
    && apt-get clean && rm -rf /tmp/* /var/tmp/*

# Virtual Environment
ENV VIRTUAL_ENV=/opt/venv
RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

# Install Dlib
ENV CFLAGS=-static

RUN pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple

RUN python3 -m pip install --upgrade pip
#    git clone -b 'v19.24' --single-branch https://github.com/davisking/dlib.git && \
#    cd dlib/ && \
#    python3 setup.py install --set BUILD_SHARED_LIBS=OFF
RUN pip3 install dlib
RUN pip3 install face_recognition


# Runtime Image
FROM python:3.8-slim-buster

COPY --from=compile /opt/venv /opt/venv
COPY --from=compile \
    # Sources
    /lib/x86_64-linux-gnu/libpthread.so.0 \
    /lib/x86_64-linux-gnu/libz.so.1 \
    /lib/x86_64-linux-gnu/libm.so.6 \
    /lib/x86_64-linux-gnu/libgcc_s.so.1 \
    /lib/x86_64-linux-gnu/libc.so.6 \
    /lib/x86_64-linux-gnu/libdl.so.2 \
    /lib/x86_64-linux-gnu/librt.so.1 \
    # Destination
    /lib/x86_64-linux-gnu/

COPY --from=compile \
    # Sources
    /usr/lib/x86_64-linux-gnu/libX11.so.6 \
    /usr/lib/x86_64-linux-gnu/libXext.so.6 \
    /usr/lib/x86_64-linux-gnu/libpng16.so.16 \
    /usr/lib/x86_64-linux-gnu/libjpeg.so.62 \
    /usr/lib/x86_64-linux-gnu/libwebp.so.6 \
    /usr/lib/x86_64-linux-gnu/libstdc++.so.6 \
    /usr/lib/x86_64-linux-gnu/libxcb.so.1 \
    /usr/lib/x86_64-linux-gnu/libXau.so.6 \
    /usr/lib/x86_64-linux-gnu/libXdmcp.so.6 \
    /usr/lib/x86_64-linux-gnu/libbsd.so.0 \
    # Destination
    /usr/lib/x86_64-linux-gnu/

# Add our packages
ENV PATH="/opt/venv/bin:$PATH"