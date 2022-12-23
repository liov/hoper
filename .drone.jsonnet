// local mode(mode="app") = if mode == "app" then "app" else "node";
local tpldir = 'build/k8s/app/';
local workspace = '/src';
local srcdir = workspace + '/';

local compileHost = {
    localhost : {
        dirprefix : '/mnt/d/',
        codedir: self.dirprefix + 'code/hoper/',
        gopath: self.dirprefix +'SDK/gopath',
    },
     tot: {
         dirprefix : '/home/new/data/',
         codedir: self.dirprefix +'code/hoper',
         gopath: self.dirprefix + 'gopath',
     }
};

local targetHost = {
    tx : {
       datadir:'/data',
       confdir:'/root/config',
    },
    tot: {
     dirprefix : '/home/new/',
     datadir: self.dirprefix + 'data',
     confdir: self.dirprefix + 'config',
    }
};

local kubectl(compile,target, cmd) = if compile == target then {
  name: 'deploy',
  image: 'bitnami/kubectl',
  user: 0,  //文档说是string类型，结果"root"不行 k8s runAsUser: 0
  volumes: [
    {
      name: 'kube',
      path: '/root/.kube/',
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
    'cd '+ tpldir + ' && chmod +x account.sh && ./account.sh ' + target,
    'cd ' + workspace,
  ] + cmd,
};


local Pipeline(group, name='', mode='app', type='bin' , workdir='tools/server/clawer', sourceFile='', protopath='', opts=[], compile='localhost',target = 'tx', schedule='') = {

  local cconfig = compileHost[compile],
  local tconfig = targetHost[target],


  local fullname = if name == '' then group else group + '-' + name,
  local committag = fullname + '-v',
  local tag = '${DRONE_TAG##' + committag + '}',
  local datadir = tconfig.datadir,
  local dockerfilepath = tpldir + 'Dockerfile-' + type,
  local deppath = tpldir + 'deploy-' + mode +'.yaml',
  kind: 'pipeline',
  type: 'docker',
  name: fullname + '-' + target,
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
        path: cconfig.codedir,
      },
    },
    {
      name: 'gopath',
      host: {
        path: cconfig.gopath,
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
  ],
  clone: {
    disable: true,
  },
  steps: [
    {
      name: 'clone && build',
      image: if protopath != '' then 'jybl/goprotoc' else 'golang:1.19.4',
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
      // "git config --global http.proxy 'socks5://proxy.tools:1080'",
      //"git config --global https.proxy 'socks5://proxy.tools:1080'",
      //"git clone ${DRONE_GIT_HTTP_URL} .",
      'cd /code',
      //  'git tag -l | xargs git tag -d',
      //'git fetch --all && git reset --hard origin/master && git pull',
      'cd ' + srcdir,
      'git clone /code .',
      'git checkout -b deploy $DRONE_COMMIT_REF',
      'cp -r /code/'+tpldir + 'certs '+ srcdir +tpldir,
       // edit Dockerfile && deploy file
      local buildfile =  '/code/' + workdir + protopath + '/build';
      if protopath != '' then 'if [ -f ' + buildfile + ' ]; then cp -r ' + protopath + ' '+ srcdir + workdir + '; fi' else 'echo',
      "sed -i 's/$${app}/" + fullname + "/g' " + dockerfilepath,
      local cmd = ['./' + fullname , '-c','./config/'+group+'.toml'] + opts;
      "sed -i 's#$${cmd}#" + std.join('", "', [opt for opt in cmd]) + "#g' " + dockerfilepath,
      "sed -i 's/$${app}/" + fullname + "/g' " + deppath,
      "sed -i 's/$${group}/" + group + "/g' " + deppath,
      "sed -i 's#$${datadir}#" + datadir + "#g' " + deppath,
      "sed -i 's#$${confdir}#" + tconfig.confdir + "#g' " + deppath,
      "sed -i 's#$${image}#jybl/" + fullname + ':' + tag + "#g' " + deppath,
      if mode == 'cronjob' then "sed -i 's#$${schedule}#" + schedule + "#g' " + deppath else 'echo',
      local bakdir = '/code/'+ tpldir + 'deploy/';
      'if [ ! -d ' + bakdir + ' ];then mkdir -p ' + bakdir + '; fi && cp -r ' + deppath + ' ' + bakdir + fullname + '.yaml && cp -r ' + dockerfilepath + ' ' + bakdir + fullname  + '-Dockerfile',
      // go build
      'cd ' + workdir,
      'go mod download',
      local genpath = srcdir + workdir + protopath;
      local buildfile = genpath + '/build';
      if protopath != '' then 'if [ ! -f ' + buildfile + ' ]; then generate go --proto='+srcdir+'proto --genpath='+genpath+'; fi' else 'echo',
      'go mod tidy',
      'go build -trimpath -o  '+ srcdir + fullname + ' ' + sourceFile,
      ],
    },
    {
      name: 'docker build',
      image: 'docker:20.10.19-cli-alpine3.16',
      privileged: true,
      volumes: [
        {
          name: 'dockersock',
          path: '/var/run/',
        },
      ],
      environment: {
          USERNAME: {
            from_secret: 'docker_username',
          },
          PASSWORD: {
            from_secret: 'docker_password',
          },

      },
    commands: [
        //'docker version',
        'docker login -u $USERNAME -p $PASSWORD',
        'docker build --rm=true -f build/k8s/app/Dockerfile-bin -t $USERNAME/' + fullname+':'+tag+' .',
        if compile != target then 'docker push $USERNAME/'+ fullname+':'+ tag,
    ],

    },
    kubectl(compile,target, [
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
  Pipeline('hoper', workdir='server/go/mod', protopath='/protobuf'),
  Pipeline('timepill', 'rbyorderid', mode='job',sourceFile='./timepill/cmd/recordby_orderid.go'),
  Pipeline('timepill', 'esload', mode='cronjob', sourceFile='./timepill/cmd/search_es.go',target='tot', schedule='00 10 * * *'),
  Pipeline('pro', sourceFile='./pro/cmd/record.go'),
  Pipeline('bilibili',  sourceFile='./bilibili/cmd/record_fav.go'),
  Pipeline('weibo',  sourceFile='./weibo/cmd/record_timer.go'),
]
