// local mode(mode="app") = if mode == "app" then "app" else "node";
local tpldir = './build/k8s/app/';
local codedir = '/home/new/code/hoper/';
local workspace = '/src';
local srcdir = workspace + '/';

local kubectl(deplocal, cmd) = if deplocal then {
  name: 'deploy',
  image: 'bitnami/kubectl',
  user: 0,  //文档说是string类型，结果"root"不行 k8s runAsUser: 0
  volumes: [
    {
      name: 'kube',
      path: '/root/.kube/',
    },
    {
      name: 'minikube',
      path: '/root/.minikube/',
    },
  ],
  commands: cmd,
} else {
  name: 'deploy',
  image: 'bitnami/kubectl',
  user: 0,  //文档说是string类型，结果"root"不行 k8s runAsUser: 0
  environment: {
    CA: {
      from_secret: 'ca',
    },
    CACRT: {
      from_secret: 'ca_crt',
    },
    CAKEY: {
      from_secret: 'ca_key',
    },
  },
  commands: [
    'chmod +x ' + tpldir + 'account.sh && ' + tpldir + 'account.sh',
  ] + cmd,
};


local Pipeline(group, name='', mode='app', type='bin' , workdir='tools/server', sourceFile='', protoc=false, opts=[], deplocal=false, schedule='') = {
  local fullname = if name == '' then group else group + '-' + name,
  local committag = fullname + '-v',
  local tag = '${DRONE_TAG##' + committag + '}',
  local datadir = if deplocal then '/home/new/data' else '/data',
  local dockerfilepath = tpldir + 'Dockerfile-' + type,
  local deppath = tpldir + 'deploy-' + mode +'.yaml',
  kind: 'pipeline',
  type: 'kubernetes',
  name: fullname + if deplocal then '-local' else '',
  metadata: {
    namespace: 'default',
  },
  platform: {
    os: 'linux',
    arch: 'amd64',
  },
  workspace: {
    path: workspace,
  },
  trigger: {
    ref: [
      'refs/tags/' + committag + '*',
    ],
  },
  volumes: [
    {
      name: 'codedir',
      host: {
        path: codedir,
      },
    },
    {
      name: 'gopath',
      host: {
        path: datadir + '/deps/gopath/',
      },
    },
    {
      name: 'dockersock',
      host: {
        path: '/var/run/',
      },
    },
    {
      name: 'kube',
      host: {
        path: '/root/.kube/',
      },
    },
    {
      name: 'minikube',
      host: {
        path: '/root/.minikube/',
      },
    },
  ],
  clone: {
    disable: true,
  },
  steps: [
    {
      name: 'clone && build',
      image: if protoc then 'jybl/goprotoc' else 'golang:1.18.1',
      volumes: [
        {
            name: 'codedir',
            path: '/code/',
        },
        {
          name: 'gopath',
          path: '/go/',
        },
      ],
      environment: {
        GOPROXY: 'https://goproxy.io,https://goproxy.cn,direct',
      },
      commands: [
      // git clone
       "git config --global http.proxy 'socks5://proxy.tools:1080'",
      "git config --global https.proxy 'socks5://proxy.tools:1080'",
      //"git clone ${DRONE_GIT_HTTP_URL} .",
      'cd /code',
      'git tag -l | xargs git tag -d',
      'git fetch --all && git reset --hard origin/master && git pull',
      'cd ' + srcdir,
      'git clone /code .',
      'git checkout -b deploy $DRONE_COMMIT_REF',
       // edit Dockerfile && deploy file
      local buildfile = '/code/' + workdir + '/protobuf/build';
      if protoc then 'if [ -f ' + buildfile + ' ]; then cp -r /code/' + workdir + '/protobuf  '+ srcdir + workdir + '; fi' else 'echo',
      "sed -i 's/$${app}/" + fullname + "/g' " + dockerfilepath,
      local cmd = ['./' + fullname] + opts;
      "sed -i 's#$${cmd}#" + std.join('", "', [opt for opt in cmd]) + "#g' " + dockerfilepath,
      "sed -i 's/$${app}/" + fullname + "/g' " + deppath,
      "sed -i 's/$${group}/" + group + "/g' " + deppath,
      "sed -i 's#$${datadir}#" + datadir + "#g' " + deppath,
      "sed -i 's#$${image}#jybl/" + fullname + ':' + tag + "#g' " + deppath,
      if mode == 'cronjob' then "sed -i 's#$${schedule}#" + schedule + "#g' " + deppath else 'echo',
      local bakdir = '/code/deploy/';
      'if [ ! -d ' + bakdir + ' ];then mkdir -p ' + bakdir + '; fi && cp -r ' + deppath + ' ' + bakdir + fullname + '-' + tag + '.yaml',
      // go build
      'cd ' + workdir,
      'go mod download',
      local genpath = srcdir + workdir + '/protobuf';
      local buildfile = srcdir + workdir + '/protobuf/build';
      if protoc then 'if [ ! -f ' + buildfile + ' ]; then generate go --proto='+srcdir+'/proto --genpath='+genpath+'; fi' else 'echo',
      'go mod tidy',
      'go build -trimpath -o  '+ srcdir + fullname + ' ' + sourceFile,
      ],
    },
    {
      name: 'docker build',
      image: 'plugins/docker',
      privileged: true,
      volumes: [
        {
          name: 'dockersock',
          path: '/var/run/',
        },
      ],
      settings: {
        username: {
          from_secret: 'docker_username',
        },
        password: {
          from_secret: 'docker_password',
        },
        repo: 'jybl/' + fullname,
        tags: tag,
        dockerfile: dockerfilepath,
        force_tag: true,
        auto_tag: false,
        daemon_off: true,
        purge: true,
        pull_image: false,
        dry_run: deplocal,
      },
    },
    kubectl(deplocal, [
      if mode == 'job' || mode == 'cronjob' then 'kubectl --kubeconfig=/root/.kube/config delete --ignore-not-found -f ' + deppath else 'echo',
      'kubectl --kubeconfig=/root/.kube/config apply -f ' + deppath,
    ]),
    {
      name: 'dingtalk',
      image: 'jybl/notify',
      settings: {
        ding_token: {
          from_secret: 'ding_token',
        },
        ding_secret: {
          from_secret: 'ding_secret',
         },
      },
    },
  ],
};

[
  Pipeline('timepill', sourceFile='./timepill/cmd/record.go',opts=['-t']),
  Pipeline('hoper', workdir='server/go/mod', protoc=true,),
  Pipeline('timepill', 'rbyorderid', mode='job',sourceFile='./timepill/cmd/recordby_orderid.go'),
  Pipeline('timepill', 'esload', mode='cronjob', sourceFile='./timepill/cmd/search_es.go', deplocal=true, schedule='00 10 * * *'),
  Pipeline('pro', sourceFile='./pro/cmd/record.go'),
  //Pipeline('timepill', sourceFile='./timepill/cmd/record.go',deplocal=true,opts=['-t']),
  Pipeline('bilibili',  sourceFile='./bilibili/cmd/record_fav.go'),
]
