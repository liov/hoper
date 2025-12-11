// local mode(mode="app") = if mode == "app" then "app" else "node";
const deployrepo = 'https://github.com/hopeio/deploy';
const workspace = '/src';

const compileHost = {
    localhost:{
        dirprefix: '/mnt/d',
        codedir: '/mnt/d/code/hopeio/hoper',
        gopath: '/mnt/d/sdk/gopath',
    },
    mint: {
        dirprefix: '/var',
        codedir: '/var/code/hopeio/hoper',
        gopath: '/var/sdk/gopath',
    },
};

const deploytHost = {
    tx: {
        datadir: '/data',
        confdir: '/root/config',
    },
    mint: {
        dirprefix: '/var',
        datadir: '/var/data',
        confdir: '/var/config',
        },
};


function Pipeline(group, buildDir = '',name = '', deploy_kind = 'deployment', build_type = 'bin',  sourceFile = '', protopath = '', opts = [], compile = 'localhost', target = 'tx', schedule = '') {

    let cconfig = compileHost[compile];
    let tconfig = deploytHost[target];

    let fullname = name === '' ? group : group + '-' + name
    let committag = fullname + '-v'
    let tag = '${DRONE_TAG##' + committag + '}'
    let protoGenpath = workspace + '/' + protopath
    let buildfile = protoGenpath + '/build';
    return  {
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
                image: `${protopath ? 'jybl/goprotoc' : 'golang:1.20'}`,
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
                    `git clone /code ${workspace}`,
                    'git checkout -b deploy $DRONE_COMMIT_REF',
                    // edit Dockerfile && deploy file
                    `cp -r /code/${protopath}/* ${protoGenpath}`,
                    `${protopath !== '' ? 'if [ ! -f ' + buildfile + ' ]; then protogen go -d -e -w -v -i ' + workspace + '/proto -o ' + protoGenpath + '; fi' : 'echo'}`,
                    // go build
                    `cd ${buildDir}`,
                    //'go mod tidy',
                    `go build -trimpath -o  /code/build/${fullname} ${sourceFile}`,
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
                    step: "all",
                    deploy_dir: "/code/build",
                    deploy_kind: deploy_kind,
                    build_type: build_type,
                    docker_cmd: `./${fullname},-c,./config/${group}.toml${opts?"":","+opts.join(",")}`,
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

            }
        ]


    }
}

const pipelines = [Pipeline(group='hoper', buildDir = 'server/go', protopath = 'server/go/protobuf')]

