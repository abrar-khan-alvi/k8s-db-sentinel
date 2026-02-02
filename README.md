# DB Sentinel

DB Sentinel is a Kubernetes-native utility designed to ensure high availability for your critical PostgreSQL database pods. It acts as a dedicated watchdog that continuously monitors the health and presence of a specific pod (by default `my-postgres`) and automatically initiates recovery procedures if the pod goes missing.

## Features

- **Continuous Monitoring:** Monitors the status of the `my-postgres` pod in the default namespace.
- **Auto-Healing:** Automatically creates a new PostgreSQL pod if the existing one is not found.
- **Dual-Mode Operation:**
  - **In-Cluster:** Runs seamlessly as a Kubernetes Pod using a ServiceAccount.
  - **Local Development:** Can be run locally using your `~/.kube/config` for testing and development.

## Prerequisites

- Go 1.x (for local development)
- Docker
- Kubernetes Cluster (Minikube, Kind, etc.)
- `kubectl` configured

## Getting Started

### Local Development

You can run the sentinel locally to monitor a cluster you have access to via your local kubeconfig.

```bash
# Clone the repository
git clone <repository-url>
cd db-sentinel

# Run the application
go run main.go
```

The application will detect it is running locally and use your `~/.kube/config`.

### Docker Build

To build the Docker image:

```bash
docker build -t db-sentinel:v1 .
```

### Kubernetes Deployment

The project includes a `deploy.yaml` file to deploy the Sentinel into your cluster with the necessary permissions.

1. **Build the image (if using Kind/Minikube, you might need to load it):**
   ```bash
   # For Kind
   kind load docker-image db-sentinel:v1
   ```

2. **Deploy resources:**
   ```bash
   kubectl apply -f deploy.yaml
   ```

   This will create:
   - `ServiceAccount`: `sentinel-sa`
   - `ClusterRole`: `sentinel-role` (Permissions: get, list, watch, create, delete pods)
   - `ClusterRoleBinding`: `sentinel-binding`
   - `Deployment`: `db-sentinel`

3. **Verify:**
   Check the logs of the sentinel pod to see it monitoring:
   ```bash
   kubectl logs -l app=db-sentinel
   ```

## Configuration

The application requires no configuration files. It defaults to monitoring a pod named `my-postgres` in the `default` namespace.

- **Environment Variables**: None currently used for configuration.
- **Command Line Flags**:
  - `-kubeconfig`: Absolute path to the kubeconfig file (optional, defaults to `~/.kube/config` when running locally).
