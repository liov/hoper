// local mode(mode="app") = if mode == "app" then "app" else "node";
local deployrepo = 'https://github.com/hopeio/deploy';
local workspace = '/src';

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

local deploytHost = {
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


local Pipeline(group, name='', deploy_kind='deployment', build_type='bin' , buildDir='', sourceFile='', protopath='', opts=[], compile='localhost',target = 'tx', schedule='') = {

  local cconfig = compileHost[compile],
  local tconfig = deploytHost[target],

  local fullname = if name == '' then group else group + '-' + name,
  local committag = fullname + '-v',
  local tag = '${DRONE_TAG##' + committag + '}',
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
      name: 'cherry',
      host: {
        path: cconfig.dirprefix + '/code/cherry/',
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
          name: 'cherry',
          path: '/cherry/',
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
       // edit Dockerfile && deploy file
      'cp -r /code/' + protopath + '/* '+ protoGenpath,

      local buildfile = protoGenpath + '/build';
      if protopath != '' then 'if [ ! -f ' + buildfile + ' ]; then protogen go -e -w -q -p '+ workspace+'/proto -g '+protoGenpath+'; fi' else 'echo',
       // go build
      'cd ' + buildDir,
      //'go mod tidy',
      'go build -trimpath -o  /code/deploy/'+ fullname + ' ' + sourceFile,
      ],
    },
     {
      name: 'deploy',
      image: 'jybl/deploy',
      privileged: true,
      user: 0,
      volumes: [
        {
          name: 'codedir',
          path: '/code/',
        },
        {
         name: 'dockersock',
         path: '/var/run/',
       },
      ],
      settings: {
          group: group,
          name: name,
          setp: "all",
          deploy_dir: "/code/deploy",
          deploy_kind: deploy_kind,
          build_type: build_type,
          docker_cmd: './' + fullname + ',-c,./config/'+group+'.toml' + std.join(',', [opt for opt in opts]),
          docker_username: {
            from_secret: 'docker_username',
          },
          docker_password: {
            from_secret: 'docker_password',
          },
          config_dir: tconfig.confdir,
          data_dir: tconfig.datadir,
          schedule: schedule,
          cluster: target,
          ca_crt: {
            from_secret: 'ca_crt',
          },
          dev_crt: {
            from_secret: 'dev_crt',
          },
          dev_key: {
            from_secret: 'dev_key',
          },
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
