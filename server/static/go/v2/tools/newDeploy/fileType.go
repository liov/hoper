package main

type FileModel struct {
}

const (
	Deployment_Go string = `---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: [name]
  namespace: [namespace]
spec:
  minReadySeconds: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: [replicas]
  template:
    metadata:
      namespace: [namespace]
      labels:
        app: [name]
    spec:
      containers:
      - name: [name]
        image: [url]/[author]/[name]:[version]
        imagePullPolicy: Always
        resources:
          # keep request = limit to keep this container in guaranteed class
          limits:
            cpu: [cpuLimit]
            memory: [memoryLimit]
          requests:
            cpu: [cpuRequest]
            memory: [memoryRequest]
        volumeMounts:
        - name: logs
          mountPath: [logPath]
        - name: config
          mountPath: /config
        command: ["/bin/sh", "-c"]
        args: ["[cmdArgs]"]
      volumes:
      - name: logs
        hostPath:
          path: [logTargetPath]
      - name: config
        configMap:
          name: [configMapName]
          items:
          - key: [configName]
            path: default.json
`
	Deployment_Web string = `---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: [name]
  namespace: [namespace]
spec:
  minReadySeconds: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: [replicas]
  template:
    metadata:
      namespace: [namespace]
      labels:
        app: [name]
    spec:
      containers:
      - name: [name]
        image: [url]/[author]/[name]:[version]
        imagePullPolicy: Always
        resources:
          # keep request = limit to keep this container in guaranteed class
          limits:
            cpu: [cpuLimit]
            memory: [memoryLimit]
          requests:
            cpu: [cpuRequest]
            memory: [memoryRequest]
        volumeMounts:
        - name: logs
          mountPath: [logPath]
        - name: config
          mountPath: /config
        command: ["/bin/sh", "-c"]
        args: ["[cmdArgs]"]
      volumes:
      - name: logs
        hostPath:
          path: [logTargetPath]
      - name: config
        configMap:
          name: [configMapName]
          items:
          - key: [configName]
            path: default.json
`
	Deployment_Java string = `---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: [name]
  namespace: [namespace]
spec:
  minReadySeconds: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: [replicas]
  template:
    metadata:
      namespace: [namespace]
      labels:
        app: [name]
    spec:
      containers:
      - name: [name]
        image: [url]/[author]/[name]:[version]
        imagePullPolicy: Always
        resources:
          # keep request = limit to keep this container in guaranteed class
          limits:
            cpu: [cpuLimit]
            memory: [memoryLimit]
          requests:
            cpu: [cpuRequest]
            memory: [memoryRequest]
        volumeMounts:
          - name: [name]
            mountPath: [logPath]
        command: ["/bin/sh","-c"]
        args: ["[cmdArgs]"]
      volumes:
        - name: [name]
          hostPath:
            path: [logTargetPath]
`
	Deployment_Go_Dev string = `---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: [name]
  namespace: [namespace]
spec:
  minReadySeconds: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: [replicas]
  template:
    metadata:
      namespace: [namespace]
      labels:
        app: [name]
    spec:
      nodeName: 192.168.1.210
      containers:
      - name: [name]
        image: [url]/[author]/[name]:[version]
        imagePullPolicy: Always
        volumeMounts:
        - name: logs
          mountPath: [logPath]
        - name: config
          mountPath: /config
        command: ["/bin/sh", "-c"]
        args: ["[cmdArgs]"]
      volumes:
      - name: logs
        hostPath:
          path: [logTargetPath]
      - name: config
        configMap:
          name: [configMapName]
          items:
          - key: [configName]
            path: default.json
`
	Deployment_Web_Dev string = `---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: [name]
  namespace: [namespace]
spec:
  minReadySeconds: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: [replicas]
  template:
    metadata:
      namespace: [namespace]
      labels:
        app: [name]
    spec:
      nodeName: 192.168.1.210
      containers:
      - name: [name]
        image: [url]/[author]/[name]:[version]
        imagePullPolicy: Always
        volumeMounts:
        - name: logs
          mountPath: [logPath]
        - name: config
          mountPath: /config
        command: ["/bin/sh", "-c"]
        args: ["[cmdArgs]"]
      volumes:
      - name: logs
        hostPath:
          path: [logTargetPath]
      - name: config
        configMap:
          name: [configMapName]
          items:
          - key: [configName]
            path: default.json
`
	Deployment_Java_Dev string = `---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: [name]
  namespace: [namespace]
spec:
  minReadySeconds: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: [replicas]
  template:
    metadata:
      namespace: [namespace]
      labels:
        app: [name]
    spec:
      nodeName: 192.168.1.210
      containers:
      - name: [name]
        image: [url]/[author]/[name]:[version]
        imagePullPolicy: Always
        volumeMounts:
          - name: [name]
            mountPath: [logPath]
        command: ["/bin/sh","-c"]
        args: ["[cmdArgs]"]
      volumes:
        - name: [name]
          hostPath:
            path: [logTargetPath]
`

	Ingress string = `---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name:  [name]
  namespace: [namespace]
  annotations:
    kubernetes.io/ingress.class: [proxyStatus]
spec:
  rules:
  - host:  [domain]
    http:
      paths:
      - path: /
        backend:
          serviceName: [name]
          servicePort: [port]
`

	Service string = `---
apiVersion: v1
kind: Service
metadata:
  name: [name]
  namespace: [namespace]
  annotations:
    prometheus.io/scrape: 'true'
    prometheus.io/path:   /metrics
    prometheus.io/port:   '8081'
  labels:
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: "[name]"
spec:
  type: ClusterIP
  ports:
    - name: [name]
      port: [port]
      targetPort: [servicePort]
      protocol: TCP
  selector:
    app: [name]
  sessionAffinity: ClientIP
`

	ServiceDev string = `---
apiVersion: v1
kind: Service
metadata:
  name: [name]
  namespace: [namespace]
  annotations:
    prometheus.io/scrape: 'true'
    prometheus.io/path:   /metrics
    prometheus.io/port:   '8081'
  labels:
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: "[name]"
spec:
  type: ClusterIP
  ports:
    - name: [name]
      port: [port]
      targetPort: [servicePort]
      protocol: TCP
      nodePort: [exportPort]
  selector:
    app: [name]
  sessionAffinity: ClientIP
`

	ServiceJar string = `---
apiVersion: v1
kind: Service
metadata:
  name: [name]
  namespace: [namespace]
  [annotations]
  labels:
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: "[name]"
spec:
  type: ClusterIP
  ports:
    - name: [name]
      port: [port]
      targetPort: [servicePort]
  selector:
    app: [name]
`

	Dockerfile_Go string = `FROM reg.hoper.xyz/library/centos:7
MAINTAINER "Hoper Team <erp@hyx.com>"

#修改容器时区
RUN rm /etc/localtime && \
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /[name]/build

ADD ./build/main /[name]/build/[name]

RUN mkdir /data && cd /data && mkdir logs

CMD ["./[name]", "-conf", "config"]
`

	Dockerfile_Java string = `FROM reg.hoper.xyz/library/tomcat:7.0.69-apm
MAINTAINER "Hoper Team <erp@hyx.com>"

#修改容器时区
RUN rm /etc/localtime && \
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ENV JAVA_OPTS $JAVA_OPTS -Duser.timezone=GMT+08

WORKDIR /usr/local/tomcat/

#将更换镜像中的tomcat配置文件
ADD ./build/config/[env]/server.xml /usr/local/tomcat/conf/
ADD ./build/config/[env]/context.xml /usr/local/tomcat/conf/

#将项目放置在镜像中的tomcat目录下
ADD  ./build/[name]/ /usr/local/tomcat/[name]/ROOT/

#将当前环境的项目配置文件放置在镜像的项目目录下
ADD ./build/config/[env]/ /usr/local/tomcat/[name]/ROOT/WEB-INF/classes/`

	Dockerfile_War string = `FROM reg.hoper.xyz/erp-library/tomcat:7.0.69-jre8
MAINTAINER "Hoper Team <erp@hyx.com>"

#修改容器时区
RUN rm /etc/localtime && \
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ENV JAVA_OPTS $JAVA_OPTS -Duser.timezone=GMT+08

WORKDIR /usr/local/tomcat/

RUN mkdir -p $JAVA_HOME/lib/fonts/fallback

ADD simsun.ttf $JAVA_HOME/lib/fonts/fallback/

RUN mv /usr/local/tomcat/conf/server.xml /usr/local/tomcat/conf/server_1.xml
ADD ./build/server.xml /usr/local/tomcat/conf/

#将项目放置在镜像中的tomcat目录下
ADD  ./build/[name].war /usr/local/tomcat/webapps/
`

	Dockerfile_Web string = `FROM reg.hoper.xyz/library/centos:7
MAINTAINER "Hoper Team <erp@hyx.com>"

#修改容器时区
RUN rm /etc/localtime && \
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /[name]/build

ADD ./build/main /[name]/build/[name]

ADD ./build/ /[name]/build/build

RUN mkdir /data && cd /data && mkdir logs

CMD ["./[name]", "-conf", "config"]
`

	Dockerfile_Jar string = `FROM reg.hoper.xyz/library/centos7-jre8:0.1.0
#FROM reg.hoper.xyz/yuwanchat/jdk:v1.8
MAINTAINER "Hoper Aura Dts <dts@hyx.com>"

#修改容器时区
RUN rm /etc/localtime && \
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /[name]

ADD ./build /[name]

#RUN java -jar [name].jar
`
)
