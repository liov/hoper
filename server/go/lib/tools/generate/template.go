package main

const ymlTpl = `schema:
  - ./*.graphqls

# Where should the generated server code go?
exec:
  filename: ../../{{.}}/generated.gql.go
  package: {{.}}

# Enable Apollo federation support
federation:
  filename: ../../{{.}}/federation.gql.go
  package: {{.}}

model:
  filename: ../../{{.}}/models.gql.go
  package: {{.}}

autobind:
  - "github.com/actliboy/hoper/server/go/protobuf/{{.}}"
  - "github.com/actliboy/hoper/server/go/lib/protobuf/response"
  - "github.com/actliboy/hoper/server/go/lib/protobuf/oauth"

models:
  ID:
    model:
      - github.com/actliboy/hoper/server/go/lib/utils/net/http/api/graphql.UInt64
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
  Int32:
    model:
      - github.com/99designs/gqlgen/graphql.Int32
  Int64:
    model:
      - github.com/99designs/gqlgen/graphql.Int64
  Uint8:
    model:
      - github.com/actliboy/hoper/server/go/lib/utils/net/http/api/graphql.Uint8
  Uint:
    model:
      - github.com/actliboy/hoper/server/go/lib/utils/net/http/api/graphql.Uint
  Uint32:
      model:
        - github.com/actliboy/hoper/server/go/lib/utils/net/http/api/graphql.Uint32
  Uint64:
      model:
        - github.com/actliboy/hoper/server/go/lib/utils/net/http/api/graphql.Uint64
  Float32:
    model:
      - github.com/actliboy/hoper/server/go/lib/utils/net/http/api/graphql.Float32
  Float64:
    model:
      - github.com/actliboy/hoper/server/go/lib/utils/net/http/api/graphql.Float64
  Float:
    model:
      - github.com/99designs/gqlgen/graphql.Float
  Bytes:
    model:
      - github.com/actliboy/hoper/server/go/lib/utils/net/http/api/graphql.Bytes
  HttpResponse_HeaderEntry:
    model:
      - github.com/actliboy/hoper/server/go/lib/utils/net/http/api/graphql.HttpResponse_HeaderEntry
`

//经过一番查找，发现yaml语法对格式是非常严格的，不可以有制表符！不可以有制表符！不可以有制表符！
//缩进也有要求
