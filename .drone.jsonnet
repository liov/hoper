// local mode(mode="app") = if mode == "app" then "app" else "node";
local tpldir = "./build/k8s/app/";

local Pipeline(group, name, mode, workdir, sourceFile, opts) = {
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
      "refs/tags/"+name+"-*"
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
         path: "/var/run/"
       }
     },
     {
       name: "kube",
       host: {
           path: "/root/.kube/"
       }
     },
      {
        name: "minikube",
        host: {
            path: "/root/.minikube/"
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
        "git clone ${DRONE_GIT_HTTP_URL} .",
        "git checkout $DRONE_COMMIT_REF",
        "sed -i 's/$${app}/"+name+"/g' "+tpldir+mode+"/Dockerfile",
        "sed -i 's/$${opts}/"+std.join(" ,",["\""+opt+"\"" for opt in opts])+"/g' "+tpldir+mode+"/Dockerfile",
        "cat "+tpldir+mode+"/Dockerfile",
        "echo",
        "sed -i 's/$${app}/"+name+"/g' "+tpldir+mode+"/deployment.yaml",
        "sed -i 's/$${group}/"+group+"/g' "+tpldir+mode+"/deployment.yaml",
        "sed -i 's#$${image}#jyblsq/"+name+":${DRONE_TAG##"+name+"-}#g' "+tpldir+mode+"/deployment.yaml"
      ]
    },
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
         GOPROXY: "https://goproxy.io,https://goproxy.cn,direct"
      },
      commands: [
        "cd " + workdir,
        "go mod download",
        "go mod tidy",
        "go build -o /drone/src/"+name+" "+sourceFile
      ]
    },
    {
      name: "docker build",
      image: "plugins/docker",
      volumes: [
        {
            name: "dockersock",
            path: "/var/run/"
        }
      ],
      settings: {
        username: {
          from_secret: "docker_username"
        },
        password: {
          from_secret: "docker_password"
        },
       repo: "jyblsq/"+name,
       tags: "${DRONE_TAG##"+name+"-}",
       dockerfile: tpldir+mode+"/Dockerfile",
       force_tag: true,
       auto_tag: false,
       daemon_off: true,
       purge: true
      }
    },
    {
      name: "deploy",
      image: "bitnami/kubectl",
      user: 0, //文档说是string类型，结果"root"不行 k8s runAsUser: 0
      volumes: [
        {
          name: "kube",
          path: "/root/.kube/"
        },
        {
          name: "minikube",
          path: "/root/.minikube/"
        }
      ],
      commands: [
        "kubectl --kubeconfig=/root/.kube/config apply -f "+tpldir+mode+"/deployment.yaml"
      ]
    },
    {
       name: "dingtalk",
       image: "lddsb/drone-dingtalk-message",
       settings: {
          token: {
            from_secret: "token"
          },
          type: "markdown",
          message_color: true,
          message_pic: true,
          sha_link: true
       }
    }
  ]
};

[
  Pipeline("timepill","timepill","app","tools/server","./timepill/cmd/record.go",["-t"]),
  Pipeline("hoper","hoper","app","server/go/mod","",[])
]