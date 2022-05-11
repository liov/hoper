local Pipeline(name, workdir, buildArg, dockerfile, deployment) = {
  kind: "pipeline",
  type: "kubernetes",
  name: name,
  metadata: {
    namespace: "default"
  },
  platform: {
    os: "linux",
    arch: "amd64"
  },
  trigger: {
    ref: [
      "refs/tags/timepill-*"
      ]
  },
  volumes: [
    {
      name: "gopath",
      host: {
         path: "/data/deps/gopath/"
      }
    },
     {
       name: "dockersock",
       host: {
         path: "/var/run/docker.sock"
       }
     },
     {
       name: "kube",
       host: {
           path: "/root/.kube/"
       }
     }
  ],
  clone: {
   disable: true
  },
  steps: [
    {
      name: "clone",
      image: "alpine/git",
      commands: [
        "git config --global https.proxy 'socks5://proxy.tools:1080'",
        "git clone https://github.com/octocat/hello-world.git .",
        "git checkout $DRONE_COMMIT_REF"
      ]
    }
    {
      name: "go build",
      image: "golang:1.18.1",
      volumes: [
        {
            name: "gopath",
            path: "/go/"
        }
      ],
      environment: {
         GOOS: "linux",
         GOARCH: "amd64",
         GOPATH: "/go",
         GOPROXY: "https://goproxy.io,https://goproxy.cn,direct"
      },
      commands: [
        "cd " + workdir,
        "go mod download",
        "go mod tidy",
        buildArg
      ]
    },
    {
      name: "docker build",
      image: "plugins/docker",
      volumes: [
        {
            name: "dockersock",
            path: "/var/run/docker.sock"
        }
      ],
      settings: {
        username: {
          from_secret: "docker_username"
        },
        password: {
          from_secret: "docker_password"
        },
       repo: "jyblsq/timepill",
       tags: ["latest","${DRONE_TAG##timepill-}"],
       dockerfile: dockerfile,
       force_tag: true
      }
    },
    {
      name: "deploy",
      image: "sinlead/drone-kubectl",
      volumes: [
        {
          name: "kube",
          path: "/root/.kube/"
        }
      ],
      commands: [
        deployment
      ]
    },
    {
       name: "notify",
       image: "plugins/slack",
       settings: {
          webhook: {
            from_secret: "wehook"
          }
       }
    }
  ]
};

[
  Pipeline("timepill","./tools/server", "go build -o ../../timepill ./timepill/cmd/record.go","./build/k8s/app/Dockerfile","kubectl apply -f ./build/k8s/app/timepill.yaml"),
]