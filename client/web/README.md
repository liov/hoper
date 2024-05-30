# lovestoryq

## Project setup
```
yarn install
```

### Compiles and hot-reloads for development
```
yarn serve
```

### Compiles and minifies for production
```
yarn build
```

### Run your unit tests
```
yarn test:unit
```

### Run your end-to-end tests
```
yarn test:e2e
```

### Lints and fixes files
```
yarn lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).

protoc -I=$DIR echo.proto \
--js_out=import_style=commonjs,binary:$OUT_DIR \
--grpc-web_out=import_style=typescript,mode=grpcweb:$OUT_DIR

- 暂时注释了wasm
- 找不到protoc-gen-js https://github.com/protocolbuffers/protobuf-javascript/releases