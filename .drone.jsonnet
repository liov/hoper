// local mode(mode="app") = if mode == "app" then "app" else "node";
local tpldir = "./build/k8s/app/";
local codedir = "/root/code/app/hoper/";

local Pipeline(group, name, mode, protoc, workdir, sourceFile, opts) = {
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
        name: "codedir",
        host: {
           path: codedir
        }
    },
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
      volumes: [
          {
              name: "codedir",
              path: "/code/"
          }
      ],
      commands: [
        "git config --global http.proxy 'socks5://proxy.tools:1080'",
        "git config --global https.proxy 'socks5://proxy.tools:1080'",
        //"git clone ${DRONE_GIT_HTTP_URL} .",
        "cd /code",
        "git fetch --all && git reset --hard origin/master && git pull",
        "cd /drone/src/",
        "git clone /code .",
        "git checkout -b deploy $DRONE_COMMIT_REF",
        local buildfile = "/code/"+workdir+"/protobuf/build";
        if protoc then "if [ -f "+buildfile+" ]; then cp -r /code/"+workdir+"/protobuf  /drone/src/"+workdir+"; fi" else "echo",
        "sed -i 's/$${app}/"+name+"/g' "+tpldir+mode+"/Dockerfile",
        local cmd = ["./"+name]+opts;
        "sed -i 's#$${cmd}#"+std.join(" ,",["\""+opt+"\"" for opt in cmd])+"#g' "+tpldir+mode+"/Dockerfile",
        "sed -i 's/$${app}/"+name+"/g' "+tpldir+mode+"/deployment.yaml",
        "sed -i 's/$${group}/"+group+"/g' "+tpldir+mode+"/deployment.yaml",
        "sed -i 's#$${image}#jyblsq/"+name+":${DRONE_TAG##"+name+"-}#g' "+tpldir+mode+"/deployment.yaml"
      ]
    },
    {
      name: "go build",
      image: if protoc then "jyblsq/golang:protoc" else "golang:1.18.1",
      volumes: [
        {
            name: "gopath",
            path: "/go/"
        }
      ],
      environment: {
         GOPROXY: "https://goproxy.io,https://goproxy.cn,direct"
      },
      commands: [
        "cd " + workdir,
        "go mod download",
        local buildfile = "/drone/src/"+workdir+"/protobuf/build";
         if protoc then "if [ ! -f "+buildfile+" ]; then go run ./protobuf; fi" else "echo",
        "go mod tidy",
        "go build -ldflags '-linkmode \"external\" -extldflags \"-static\"' -o  /drone/src/"+name+" "+sourceFile
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
       purge: true,
       pull_image: false
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
  Pipeline("timepill","timepill","app",false,"tools/server","./timepill/cmd/record.go",["-t"]),
  Pipeline("hoper","hoper","app",true,"server/go/mod","",[])
]