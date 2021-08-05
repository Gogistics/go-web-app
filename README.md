
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

## Envoy proxy (in progress)

```sh
# generate cert and key
$ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj '/CN=atai.com'

$ docker run -d \
      --name atai_envoy \
      -p 80:80 -p 443:443 -p 10000:10000 \
      --network atai_envoy \
      --ip "172.18.0.10" \
      --log-opt mode=non-blocking \
      --log-opt max-buffer-size=5m \
      --log-opt max-size=100m \
      --log-opt max-file=5 \
      alantai/envoy:v1

# update /etc/hosts
# e.g. add 0.0.0.0 dev.atai.com

```
Ref:
- https://hub.docker.com/r/envoyproxy/envoy-alpine-dev
- https://www.envoyproxy.io/docs/envoy/latest/start/docker

Issues:
- https://pjausovec.medium.com/the-v2-xds-major-version-is-deprecated-and-disabled-by-default-envoy-60672b1968cb
- https://stackoverflow.com/questions/63712716/envoy-proxy-v3-api-with-http-and-https-both


## Bazel 

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

# run app in container
$ docker run -d \
    -p 8443:443 \
    --name atai_api \
    --network atai_envoy \
    --ip "172.18.0.11" \
    --log-opt mode=non-blocking \
    --log-opt max-buffer-size=5m \
    --log-opt max-size=100m \
    --log-opt max-file=5 \
    alantai/api-app:atai-v0.0.0

# test
$ curl -k https://0.0.0.0:8443/api/v1/hello # {"Name":"Alan","Hobbies":["workout","programming","driving"]}

# after successfully push dokcer image to docker registry
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



## Remote caching (to be continued)

Ref:
- https://docs.bazel.build/versions/main/remote-caching.html



## React build by Bazel (to be continued)

Ref:
- https://github.com/salrashid123/go-grpc-bazel-docker
- https://github.com/thelgevold/react-bazel-example
- https://www.syntaxsuccess.com/viewarticle/large-react-production-bazel-build

