FROM golang:1.21 as base

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

FROM blockstream/gdk-ubuntu-builder:c59e04d07ba0b61e70883ab7fe8bbeb4795ca48c as gdk

RUN git clone https://github.com/Blockstream/gdk --depth 1
RUN export PATH="/root/.cargo/bin:$PATH" && cd gdk && ./tools/build.sh --gcc --buildtype release --no-deps-rebuild --external-deps-dir /prebuild/gcc --parallel 1


#RUN gdk/docker/debian/install_deps.sh && gdk/docker/debian/install_rust_tools.sh
#RUN cd gdk && tools/build.sh --gcc

# Build
RUN CGO_ENABLED=1 go build main.go

# Run
CMD ["/example"]
