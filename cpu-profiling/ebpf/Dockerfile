FROM ubuntu:bionic-20210325

ENV PATH /usr/share/bcc/tools:$PATH
RUN sed -i "s#deb http://deb.debian.org/debian buster main#deb http://deb.debian.org/debian buster main contrib non-free#g" /etc/apt/sources.list

ENV DEBIAN_FRONTEND=nonintercative

RUN apt-get update && apt-get install -y \
    ca-certificates \
    clang \
    curl \
    gcc \
    git \
    g++ \
    --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

# Install dependencies for libbcc
# FROM: https://github.com/iovisor/bcc/blob/master/INSTALL.md#install-build-dependencies
RUN apt-get update && apt-get install -y \
    debhelper \
    cmake \
    libllvm3.9 \
    llvm-dev \
    libclang-dev \
    libelf-dev \
    bison \
    flex \
    libedit-dev \
    clang-format \
    python \
    python3-pyroute2 \
    luajit \
    libluajit-5.1-dev \
    arping \
    iperf \
    ethtool \
    devscripts \
    zlib1g-dev \
    libfl-dev \
    containerd \
    --no-install-recommends \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common \
    && rm -rf /var/lib/apt/lists/*


ENV VERSION="v1.26.0"
RUN curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/$VERSION/crictl-${VERSION}-linux-amd64.tar.gz --output crictl-${VERSION}-linux-amd64.tar.gz && \
    tar zxvf crictl-$VERSION-linux-amd64.tar.gz -C /usr/local/bin && \
    rm -f crictl-$VERSION-linux-amd64.tar.gz

#linux-headers-azure for EBPF in the cloud :)
RUN curl -LO https://mirrors.edge.kernel.org/ubuntu/pool/main/l/linux-azure-5.4/linux-headers-5.4.0-1103-azure_5.4.0-1103.109~18.04.1_amd64.deb && \  
    curl -LO https://mirrors.edge.kernel.org/ubuntu/pool/main/l/linux-azure-5.4/linux-azure-5.4-headers-5.4.0-1103_5.4.0-1103.109~18.04.1_all.deb && \
    dpkg -i linux-azure-5.4-headers-5.4.0-1103_5.4.0-1103.109~18.04.1_all.deb && \
    dpkg -i linux-headers-5.4.0-1103-azure_5.4.0-1103.109~18.04.1_amd64.deb

RUN mkdir /work
WORKDIR /work

RUN git clone --depth 1 https://github.com/brendangregg/FlameGraph

ENV BCC_VERSION v0.20.0
RUN git clone --depth 1 --branch "$BCC_VERSION" https://github.com/iovisor/bcc.git /usr/src/bcc \
	&& ( \
		cd /usr/src/bcc \
		&& mkdir build \
		&& cd build \
		&& cmake .. -DCMAKE_INSTALL_PREFIX=/usr \
		&& make \
		&& make install \
	)
  
ENTRYPOINT ["/bin/bash"]
