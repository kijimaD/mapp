###########
# builder #
###########

FROM golang:1.20-buster AS builder
RUN apt update \
    && apt install -y --no-install-recommends \
    upx-ucl
RUN apt install -y \
    gcc \
    libc6-dev \
    libgl1-mesa-dev \
    libxcursor-dev \
    libxi-dev \
    libxinerama-dev \
    libxrandr-dev \
    libxxf86vm-dev \
    libasound2-dev \
    pkg-config \
    xorg-dev \
    libx11-dev \
    libopenal-dev

WORKDIR /build
COPY . .

RUN GO111MODULE=on go build -o ./bin/mapp . \
 && upx-ucl --best --ultra-brute ./bin/mapp

###########
# release #
###########

FROM gcr.io/distroless/static-debian11:latest AS release

COPY --from=builder /build/bin/mapp /bin/
WORKDIR /workdir
ENTRYPOINT ["/bin/mapp"]

########
# node #
########

FROM node:18 as releaser
RUN yarn install
