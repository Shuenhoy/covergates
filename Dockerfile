FROM golang:alpine AS build
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.zju.edu.cn/g' /etc/apk/repositories
RUN apk --update add musl-dev
RUN apk --update add util-linux-dev
RUN apk --update add gcc g++
RUN apk --update add nodejs npm
WORKDIR /go/src/github.com/covergates/covergates
RUN npm config set registry https://registry.npm.taobao.org
COPY web/package.json ./web/package.json
COPY web/package-lock.json ./web/package-lock.json
RUN cd web && npm install 
RUN node --version && npm --version
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go env -w GOBIN=/bin
RUN go install github.com/bradrydzewski/togo@latest
COPY web ./web
RUN go generate ./web
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -v -o covergates ./cmd/server
FROM alpine
COPY --from=build /go/src/github.com/covergates/covergates/covergates /covergates
ENTRYPOINT [ "/covergates" ]
