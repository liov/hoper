helm repo add openkruise https://openkruise.github.io/charts/
helm repo update
helm install kruise openkruise/kruise --version 1.1.0

helm upgrade kruise openkruise/kruise --version 1.1.0 --reset-values --force