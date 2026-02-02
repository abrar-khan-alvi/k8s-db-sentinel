# 🛡️ DB Sentinel: Kubernetes Self-Healing Operator

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Kubernetes](https://img.shields.io/badge/Kubernetes-v1.28-326CE5?style=flat&logo=kubernetes)
![Docker](https://img.shields.io/badge/Docker-v24-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green)

**DB Sentinel** is a custom Kubernetes Controller (Operator) designed to ensure high availability for PostgreSQL database pods. 

It implements the **Kubernetes Reconciliation Loop** pattern to continuously monitor the health of critical database infrastructure and automatically initiates recovery procedures—without human intervention—if a pod fails or is deleted.

---

## 📸 Demo & Architecture

### The "Self-Healing" in Action
*(The system detecting a failure and repairing it in <2 seconds)*

![Self Healing Demo](./http://github.com/abrar-khan-alvi/k8s-db-sentinel/blob/main/Self%20healing.png)


### How it Works (Logic Flow)
This operator mimics the logic used by enterprise tools like **KubeDB**. It runs a control loop that constantly compares the *Desired State* (Pod exists) with the *Actual State*.

```mermaid
sequenceDiagram
    participant Admin as DevOps Eng
    participant K8s as Kubernetes API
    participant DB as Postgres Pod
    participant Sentinel as DB Sentinel (Operator)

    Note over Sentinel: 1. Watch Loop (Active)
    Sentinel->>K8s: GET /pods/my-postgres
    K8s-->>Sentinel: Status: Running [OK]
    
    Admin->>K8s: kubectl delete pod my-postgres (Simulate Crash)
    K8s-->>DB: Terminate Signal
    
    Note over Sentinel: 2. Detect Drift
    Sentinel->>K8s: GET /pods/my-postgres
    K8s-->>Sentinel: Error: NotFound [ALERT]
    
    Note over Sentinel: 3. Reconciliation (Heal)
    Sentinel->>K8s: POST /pods (Create New Pod)
    K8s-->>DB: Launch New Instance
    Sentinel->>Sentinel: Log [SUCCESS] Recovery Complete

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
