用户账户和用户组
Kubernetes 并不会存储由认证插件从客户端请求中提取出的用户及所属组的信息，它们仅仅用于检验用户是否有权限执行其所请求的操作。

客户端访问API服务的途径通常有三种：kubectl、客户端库或者直接使用 REST接口进行请求。

而可以执行此类请求的主体也被 Kubernetes 分为两类：现实中的“人”和 Pod 对象，

它们的用户身份分别对应于常规用户 (User Account ）和服务账号 （ Service Account） 。

Use Account（用户账号）：一般是指由独立于Kubernetes之外的其他服务管理的用 户账号，例如由管理员分发的密钥、Keystone一类的用户存储（账号库）、甚至是包 含有用户名和密码列表的文件等。Kubernetes中不存在表示此类用户账号的对象， 因此不能被直接添加进 Kubernetes 系统中 。

Service Account（服务账号）：是指由Kubernetes API 管理的账号，用于为Pod 之中的服务进程在访问Kubernetes API时提供身份标识（ identity ） 。Service Account通常要绑定于特定的命名空间，它们由 API Server 创建，或者通过 API 调用于动创建 ，附带着一组存储为Secret的用于访问API Server的凭据。

Kubernetes 有着以下几个内建的用于特殊目的的组 。system:unauthenticated ：未能通过任何一个授权插件检验的账号，即未通过认证测 试的用户所属的组 。system :authenticated ：认证成功后的用户自动加入的一个组，用于快捷引用所有正常通过认证的用户账号。system : serviceaccounts ：当前系统上的所有 Service Account 对象。system :serviceaccounts :<namespace＞：特定命名空间内所有的 Service Account 对象。

用户验证
尽管K8S认知用户靠的只是用户的名字，但是只需要一个名字就能请求K8S的API显然是不合理的，所以依然需要验证此用户的身份

下面我们来创建一个User Account，测试访问某些我们授权的资源：

在K8S中，有以下几种验证方式：

X509客户端证书
客户端证书验证通过为API Server指定--client-ca-file=xxx选项启用，API Server通过此ca文件来验证API请求携带的客户端证书的有效性，一旦验证成功，API Server就会将客户端证书Subject里的CN属性作为此次请求的用户名
静态token文件
通过指定--token-auth-file=SOMEFILE选项来启用bearer token验证方式，引用的文件是一个包含了 token,用户名,用户ID 的csv文件 请求时，带上Authorization: Bearer 31ada4fd-adec-460c-809a-9e56ceb75269头信息即可通过bearer token验证
静态密码文件
通过指定--basic-auth-file=SOMEFILE选项启用密码验证，类似的，引用的文件时一个包含 密码,用户名,用户ID 的csv文件 请求时需要将Authorization头设置为Basic BASE64ENCODED(USER:PASSWORD)
这里只介绍客户端验证

创建k8s User Account

1、创建证书
创建user私钥
-- 202210117 验证可行
```bash
cp /var/lib/minikube/certs/ca.crt  .
cp /var/lib/minikube/certs/ca.key .

# 创建根证书

(umask 077;openssl genrsa -out dev.key 2048)
openssl req -new -key dev.key -out dev.csr -subj "/O=k8s/CN=dev"
openssl  x509 -req -in dev.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out dev.crt -days 365
```
1、创建配置文件
创建配置文件主要有以下几个步骤：
```bash
kubectl config set-cluster --kubeconfig=/PATH/TO/SOMEFILE      #集群配置
kubectl config set-credentials NAME --kubeconfig=/PATH/TO/SOMEFILE #用户配置
kubectl config set-context    #context配置
kubectl config use-context    #切换context
```
* --embed-certs=true的作用是不在配置文件中显示证书信息。
* --kubeconfig=/root/cbmljs.conf用于创建新的配置文件，如果不加此选项,则内容会添加到家目录下.kube/config文件中，可以使用use-context来切换不同的用户管理k8s集群。
* context简单的理解就是用什么用户来管理哪个集群，即用户和集群的结合。
* 创建集群配置
```Bash
# 202210117 验证可行
kubectl config set-cluster k8s --server=https://192.168.253.136:6443 --certificate-authority=ca.crt --embed-certs=true --kubeconfig=/root/.kube/dev.conf
```
创建用户配置
```Bash
kubectl config set-credentials dev --client-certificate=dev.crt --client-key=dev.key --embed-certs=true --kubeconfig=/root/.kube/dev.conf
```
创建context配置
```Bash
kubectl config set-context dev --cluster=k8s --user=dev --kubeconfig=/root/.kube/dev.conf
```
切换context
```Bash
kubectl config use-context dev --kubeconfig=/root/.kube/dev.conf
```
kubectl config view

或者username password方式也可以、
```Bash
# 创建cluster
kubectl config set-cluster set-cluster k8s --server=http://10.10.3.127 --insecure-skip-tls-verify
# 创建user
kubectl config set-credentials  admin-credentials --username=admin --password=123456
# 创建context
kubectl config set-context dev  --cluster=k8s --namespace=kube-system --user=dev
# 指定当前使用的context
kubectl config use-context dev
```

创建Role
```bash
cat >role.yaml <<EOF
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: default
  name: dev
rules:
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - create
      - delete
  - apiGroups:
      - ""
      - "apps"
      - "batch"
    resources:
      - configmaps
      - services
      - apisixroutes
      - apisixtlses
      - deployments
      - jobs
      - cronjobs
      - pods
      - pods/log
    verbs:
      - get
      - create
      - delete
      - list
      - watch
      - update
      - patch
EOF
```
创建Rolebinding
```bash
cat >rolebinding.yaml <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: dev
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: dev
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: User
    name: dev
EOF
```
到这一步就可以进行验证了

kubectl get pod

我们是可以查看查看default命名空间的pod，但是其他空间的pod是无法查看的。

创建ClusterRole
cat cluster-reader.yaml
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cluster-reader
rules:
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
  - list
  - watch
```
创建ClusterRoleBinding
cat cluster-role.yaml
```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cluster-role
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-role
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: dev
```
验证结果
kubectl get pod --all-namespaces

就可以看到所有命名空间的pod了.
权限绑定指定的namespace

也可以使用下面方法进行绑定

kubectl get clusterrole  查看系统自带角色
kubectl create rolebinding devuser-admin-rolebinding（rolebinding的名字） --clusterrole=admin（clusterrole的名字，admin在k8s所有namespace下都有最高权限） --user=devuser（将admin的权限赋予devuser用户） --namespace=dev（范围是dev这个namespace下） 即dev


扩展:

kubectl api-resources 可以查看apiGroups

创建集群角色

cat  clusterrole.yaml
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: clusterrole
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: ["extensions", "apps"]
    resources: ["deployments"]
    verbs: ["get", "watch", "list"]
  - apiGroups: [""]
    resources: ["pods/exec"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: [""]
    resources: ["pods/log"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: [""]
    resources: ["namespaces","namespaces/status"]
    verbs: ["*"]   # 也可以使用['*']
  - apiGroups: ["","apps","extensions","apiextensions.k8s.io"]
    resources: ["role","replicasets","deployments","customresourcedefinitions","configmaps"]
```
集群绑定
cat  clusterrole.yaml
```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cluster-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: clusterrole
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: test
```