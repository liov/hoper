# Execute the build command in the directory where the Dockerfile is located (by default, the project root), specifying the tag manually.
$ docker build -t apisix-dashboard:$tag .

# For users in mainland China, the `ENABLE_PROXY` parameter can be provided to speed up module downloads.
$ docker build -t apisix-dashboard:$tag . --build-arg ENABLE_PROXY=true

# If you want to use the latest codes to build, you can specify the `APISIX_DASHBOARD_VERSION` parameter to `master`.
# This parameter can also be specified as branch name of a specific version, such as `v2.1.1`.
$ docker build -t apisix-dashboard:$tag . --build-arg APISIX_DASHBOARD_VERSION=master

kubectl create configmap apisix-dashboard-conf.yaml --from-file=conf.yaml -n ingress-apisix