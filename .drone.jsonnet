// local mode(mode="app") = if mode == "app" then "app" else "node";
local deployrepo = 'https://github.com/hopeio/deploy';
local workspace = '/src';
local deploydir = '/deploy';
local deploytpl = '/code/deploy/';

local compileHost = {
    localhost : {
        dirprefix : '/mnt/d',
        codedir: self.dirprefix + '/code/hoper/',
        gopath: self.dirprefix +'/SDK/gopath',
    },
     tot: {
         dirprefix : '/home/new/data',
         codedir: self.dirprefix +'/code/hoper',
         gopath: self.dirprefix + '/gopath',
     }
};

local targetHost = {
    tx : {
       datadir:'/data',
       confdir:'/root/config',
    },
    tot: {
     dirprefix : '/home/new',
     datadir: self.dirprefix + '/data',
     confdir: self.dirprefix + '/config',
    }
};

local kubectl(compile, target, cmd) = if compile == target then {
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
    CACRT: {
      from_secret: 'ca_crt',
    },
    DEVCRT: {
      from_secret: 'dev_crt',
    },
    DEVKEY: {
      from_secret: 'dev_key',
    },
  },
  commands: [
    'cd deploy/cmd && chmod +x account.sh && ./account.sh ' + target,
    'cd ' + workspace,
  ] + cmd,
};


local Pipeline(group, name='', mode='app', type='bin' , buildDir='', sourceFile='', protopath='', opts=[], compile='localhost',target = 'tx', schedule='') = {

  local cconfig = compileHost[compile],
  local tconfig = targetHost[target],

  local fullname = if name == '' then group else group + '-' + name,
  local committag = fullname + '-v',
  local tag = '${DRONE_TAG##' + committag + '}',
  local datadir = tconfig.datadir,
  local dockerfilepath = deploydir + '/tpl/Dockerfile-' + type,
  local deppath = deploydir + '/tpl/deploy-' + mode +'.yaml',
  local protoGenpath = workspace + '/' + protopath,
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
      name: 'tailmon',
      host: {
        path: cconfig.dirprefix + '/code/tailmon/',
      },
    },
     {
      name: 'deploy',
      host: {
        path: cconfig.dirprefix + deploytpl,
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
      image: if protopath != '' then 'jybl/goprotoc' else 'golang:1.20',
      volumes: [
        {
            name: 'codedir',
            path: '/code/',
        },
        {
          name: 'tailmon',
          path: '/tailmon/',
        },
         {
          name: 'deploy',
          path: '/deploy/',
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
      //  'git tag -l | xargs git tag -d',
      //'git fetch --all && git reset --hard origin/master && git pull',
      'git clone /code ' + workspace,
      'git checkout -b deploy $DRONE_COMMIT_REF',
      'if [ ! -d /deploy/.git ]; then git clone '+ deployrepo + ' ' + deploydir + '; fi',
      'mkdir deploy',
      'cp -r '+ deploytpl +'certs deploy/certs',
      'cp -r /deploy/cmd deploy/cmd',
      'cp -r /deploy/tz deploy/tz',

       // edit Dockerfile && deploy file
      local buildfile =  '/code/' + protopath + '/build';
      if protopath != '' then 'if [ -f ' + buildfile + ' ]; then cp -r /code/' + protopath + '/* '+ protoGenpath + '; fi' else 'echo',

      'if [ ! -d ' + deploytpl + ' ];then mkdir -p ' + deploytpl + '; fi',

      local cmd = ['./' + fullname , '-c','./config/'+group+'.toml'] + opts;
      "sed -e 's/$${app}/" + fullname + "/g;s#$${cmd}#" + std.join('", "', [opt for opt in cmd]) + "#g' " + dockerfilepath + ' > '+ deploytpl + fullname + '-Dockerfile',
      "sed -e 's/$${app}/" + fullname + "/g;s/$${group}/" + group + "/g;s#$${datadir}#" + datadir + "#g;s#$${confdir}#" + tconfig.confdir + "#g;s#$${image}#jybl/" + fullname + ':' + tag + "#g;s#$${schedule}#" + schedule + "#g' " + deppath + ' > '+ deploytpl + fullname + '.yaml',
      'cp ' + deploytpl + fullname + '-Dockerfile deploy/Dockerfile',
      'cp ' + deploytpl + fullname + '.yaml deploy/deploy.yaml',
      // go build
      'cd ' + buildDir,

      local buildfile = protoGenpath + '/build';
      if protopath != '' then 'if [ ! -f ' + buildfile + ' ]; then protogen go -e -w -q -p '+ workspace+'/proto -g '+protoGenpath+'; fi' else 'echo',
      //'go mod tidy',
      'go build -trimpath -o  '+ workspace +'/deploy/'+ fullname + ' ' + sourceFile,
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
        'docker build -f deploy/Dockerfile -t $USERNAME/' + fullname+':'+tag+' deploy',
        if compile != target then 'docker push $USERNAME/'+ fullname+':'+ tag,
    ],
    },
    kubectl(compile,target, [
      if mode == 'job' || mode == 'cronjob' then 'kubectl --kubeconfig=/root/.kube/config delete --ignore-not-found -f deploy/deploy.yaml' else 'echo',
      'kubectl --kubeconfig=/root/.kube/config apply -f deploy/deploy.yaml',
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
  Pipeline('hoper', buildDir='server/go', protopath='server/go/protobuf'),
]
