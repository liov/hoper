FROM bitnami/kubectl

USER root
COPY --chown=root config $HOME/.kube/config
ENV KUBECONFIG=$HOME/.kube/config