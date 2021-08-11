# Golang, Bazel, Envoy, etc.

## Golang app
1. create a container for running Goalng
2. go into the container and init a Golang module
```sh
# 2-1.
$ go mod init github.com/Gogistics/go-web-app

# 2-2. build a simple web app, main.go, under /api-app

# 2-3. make sure the app can run successfully
$ go run main.go

# 2-4. let's make it better with a test file, main_test.go
$ go test ./...
```

## Bazel 

### General setup and build
3. write WORKSPACE and its corresponding BUILD.bazel and run the following commands
```sh
# run the gazelle target specified in the BUILD file
$ bazel run //:gazelle

# update repos
$ bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies

```

4. build package and container
```sh
$ bazel build //api-app
# try to run the web app by bazel
$ bazel run //api-app

# build container
$ bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //api-app:atai-v0.0.0
$ bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //api-app:atai-v0.0.0

# after building the image, check if the image exists
$ docker images # in my case, the image repository is alantai/api-app and the tag is atai-v0.0.0

# run app in container; if you are going to test the app through Envoy, -p 8443:443 can be removed
$ docker run -d \
    -p 8443:443 \
    --name atai_api_1 \
    --network atai_envoy \
    --ip "172.18.0.11" \
    --log-opt mode=non-blocking \
    --log-opt max-buffer-size=5m \
    --log-opt max-size=100m \
    --log-opt max-file=5 \
    alantai/api-app:atai-v0.0.0

$ docker run -d \
    --name atai_api_2 \
    --network atai_envoy \
    --ip "172.18.0.12" \
    --log-opt mode=non-blocking \
    --log-opt max-buffer-size=5m \
    --log-opt max-size=100m \
    --log-opt max-file=5 \
    alantai/api-app:atai-v0.0.0

# test the golang app running in atai_api_1
$ curl -k https://0.0.0.0:8443/api/v1/hello # response: {"Name":"Alan","Hobbies":["workout","programming","driving"]}

# login to the registry and push the docker image to the container registry
$ bazel run //api-app:push
```

5. let bazel run tests
```sh
$ bazel test --sandbox_debug //api-app/...
```

6. .bazelrc, the Bazel configuration file
Bazel accepts many options. Some options are varied frequently (for example, --subcommands) while others stay the same across several builds (such as --package_path). To avoid specifying these unchanged options for every build (and other commands), you can specify options in a configuration file.


Ref:
- https://docs.bazel.build/versions/main/guide.html
- https://github.com/bazelbuild/bazel-gazelle
- https://www.even.com/blog/testing-fast-easy-bazel
- https://zhuanlan.zhihu.com/p/203325500
- https://gist.github.com/6174/9ff5063a43f0edd82c8186e417aae1dc

Issues:
- https://stackoverflow.com/questions/59019448/what-is-the-difference-between-importmap-and-importpath-in-bazel-build-f


### Remote caching (to be continued)

Ref:
- https://docs.bazel.build/versions/main/remote-caching.html
- https://github.com/buchgr/bazel-remote


### React build by Bazel (to be continued)


```sh
# build react app
$ bazel build //react-app

# unzip tar file to test
$ tar -xf react-all.tar.gz   

```
Ref:
- https://github.com/bazelbuild/rules_nodejs/tree/stable/examples/create-react-app
- https://bazelbuild.github.io/rules_nodejs/examples#react
- https://github.com/bazelbuild/rules_nodejs/blob/stable/examples/webapp/BUILD.bazel
- https://github.com/bazelbuild/bazel/blob/master/site/docs/skylark/faq.md
- https://github.com/thelgevold/react-bazel-example
- https://github.com/bazelbuild/rules_nodejs/tree/stable/examples/jest
- https://www.syntaxsuccess.com/viewarticle/large-react-production-bazel-build
- https://stackoverflow.com/questions/53734988/typescript-how-to-include-imported-images-in-the-output-directory
- https://duncanleung.com/typescript-module-declearation-svg-img-assets/


## Envoy proxy (in progress)

```sh
# generate cert and key
$ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj '/CN=atai.com'

# build image
$ docker build -t alantai/envoy:v1 -f ./infra/dockerfiles/Dockerfile.envoy .

# spin up container
$ docker run -d \
      --name atai_envoy \
      -p 80:80 -p 443:443 -p 10000:10000 -p 8001:8001 \
      --network atai_envoy \
      --ip "172.18.0.10" \
      --log-opt mode=non-blocking \
      --log-opt max-buffer-size=5m \
      --log-opt max-size=100m \
      --log-opt max-file=5 \
      alantai/envoy:v1

# update /etc/hosts of your local environment or the place you want to test the setup
# e.g. add 0.0.0.0 atai.com

# visit home page with a table of links to all available options => http://atai.com:8001/

# query to print a textual table of all available options; please refer to the official page for more information
$ curl http://atai.com:8001/help

# test the golang app through Envoy and the response looks as below
$ curl -k https://atai.com/api/v1/hello -vvv

# *   Trying 0.0.0.0...
# * TCP_NODELAY set
# * Connected to atai.com (127.0.0.1) port 443 (#0)
# * ALPN, offering h2
# * ALPN, offering http/1.1
# * successfully set certificate verify locations:
# *   CAfile: /etc/ssl/cert.pem
#   CApath: none
# * TLSv1.2 (OUT), TLS handshake, Client hello (1):
# * TLSv1.2 (IN), TLS handshake, Server hello (2):
# * TLSv1.2 (IN), TLS handshake, Certificate (11):
# * TLSv1.2 (IN), TLS handshake, Server key exchange (12):
# * TLSv1.2 (IN), TLS handshake, Request CERT (13):
# * TLSv1.2 (IN), TLS handshake, Server finished (14):
# * TLSv1.2 (OUT), TLS handshake, Certificate (11):
# * TLSv1.2 (OUT), TLS handshake, Client key exchange (16):
# * TLSv1.2 (OUT), TLS change cipher, Change cipher spec (1):
# * TLSv1.2 (OUT), TLS handshake, Finished (20):
# * TLSv1.2 (IN), TLS change cipher, Change cipher spec (1):
# * TLSv1.2 (IN), TLS handshake, Finished (20):
# * SSL connection using TLSv1.2 / ECDHE-RSA-CHACHA20-POLY1305
# * ALPN, server accepted to use h2
# * Server certificate:
# *  subject: CN=atai.com
# *  start date: Jun 17 20:00:04 2021 GMT
# *  expire date: Jun 17 20:00:04 2022 GMT
# *  issuer: CN=atai.com
# *  SSL certificate verify result: self signed certificate (18), continuing anyway.
# * Using HTTP2, server supports multi-use
# * Connection state changed (HTTP/2 confirmed)
# * Copying HTTP/2 data in stream buffer to connection buffer after upgrade: len=0
# * Using Stream ID: 1 (easy handle 0x7fcaef80d600)
# > GET /api/v1/hello HTTP/2
# > Host: atai.com
# > User-Agent: curl/7.64.1
# > Accept: */*
# > 
# * Connection state changed (MAX_CONCURRENT_STREAMS == 2147483647)!
# < HTTP/2 200 
# < content-type: applicaiton/json; charset=utf-8
# < content-length: 61
# < date: Fri, 06 Aug 2021 17:16:10 GMT
# < x-envoy-upstream-service-time: 1
# < server: envoy
# < 
# * Connection #0 to host atai.com left intact

```
Ref:
- https://www.envoyproxy.io/docs/envoy/latest/start/docker
- https://www.envoyproxy.io/docs/envoy/latest/operations/admin#
- https://www.envoyproxy.io/docs/envoy/latest/start/sandboxes/tls
- https://hub.docker.com/r/envoyproxy/envoy-alpine-dev
- https://pi3g.com/2019/01/17/envoy-as-http-2-front-proxy-enabling-http-2-for-envoy-aka-h2/
- https://myview.rahulnivi.net/api-gateway-envoy-docker/
- https://github.com/salrashid123/go-grpc-bazel-docker

Issues:
- https://github.com/gliderlabs/docker-alpine/issues/52
- https://pjausovec.medium.com/the-v2-xds-major-version-is-deprecated-and-disabled-by-default-envoy-60672b1968cb
- https://stackoverflow.com/questions/63712716/envoy-proxy-v3-api-with-http-and-https-both
- https://github.com/salrashid123/envoy_discovery/issues/3
