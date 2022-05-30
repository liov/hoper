helm repo add gitlab-jh https://charts.gitlab.cn/
helm repo update
helm get values gitlab > gitlab.yaml
helm install gitlab gitlab-jh/gitlab -f gitlab.yaml -n tools
helm upgrade gitlab gitlab-jh/gitlab -f gitlab.yaml -n tools

helm uninstall gitlab
kubectl get pvc,secret -lrelease=gitlab