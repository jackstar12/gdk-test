FROM golang:1.21 as builder
FROM blockstream/gdk-ubuntu-builder@sha256:f75606b4fd1c681ea6bba452dcf2a1e119b36b027f131b63c0fb83ceab45015f as gdk

RUN git clone https://github.com/Blockstream/gdk --depth 1
RUN export PATH="/root/.cargo/bin:$PATH" && cd gdk && ./tools/build.sh --gcc --buildtype release --no-deps-rebuild --external-deps-dir /prebuild/gcc --parallel 1

COPY --from=builder /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:$PATH"

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./
RUN cp /root/gdk/gdk/build-gcc/libgreenaddress_full.a /app/wallet/lib/libgreenaddress_full.a

# Build
RUN CC=gcc CGO_ENABLED=1 go build main.go

# Run
CMD ["/example"]
