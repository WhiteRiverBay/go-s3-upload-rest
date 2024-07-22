FROM amazonlinux:2

# RUN amazon-linux-extras install epel -y
# RUN yum update -y
RUN yum install -y tar
RUN yum install -y gzip

# Install go1.22.4
RUN curl -O https://dl.google.com/go/go1.22.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz
RUN rm go1.22.2.linux-amd64.tar.gz

# Install git
RUN yum install -y git

# Set go path
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go


WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

# Copy the app source
COPY . .

# Compile the go app
RUN go build -o /go/bin/go-s3 ./cmd/main.go

WORKDIR /go/bin
# Remote the source
RUN rm -rf /go/src/app

# Run the app
ARG PORT=9090
ARG HOST="0.0.0.0"
CMD ["/go/bin/app", "-bind", "${HOST}:${PORT}" "-access", "${ACCESS_KEY}", "-secret", "${SECRET_KEY}", "-region", "${REGION}", "-bucket", "${BUCKET}", "-daily", "10", "-minute", "2"]
EXPOSE ${PORT}