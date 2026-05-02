Kubernetes Database Operator 🚀
[ [

Production-ready CRD + Controller-runtime operator with AUTO PORT CLEANUP 🔧

🎯 What is This?
Custom Resource: Database - Declarative PostgreSQL
CRD: databases.laxmangodi.com/v1
Controller: Full lifecycle + Smart port management

🏭 Controller Lifecycle (Factory Analogy)
text
👥 Customer Orders (CR Create) → 📋 Order Book (API Server)
    ↓
🔔 Bell Rings (EventSource Watch) → 👷 Foreman (Controller)
    ↓
📦 Local Inventory (Cache Sync) → 🎫 Job Queue (Workqueue)
    ↓
🔧 Workers Build (Reconcile) → 📊 Dashboard (Metrics)
    ↓
✅ Status "Shipped" → 🔄 Watch Next Order
🚀 Smart Auto Port Cleanup 🔧
Zero port conflicts - EVER!

1. Production Script (run-operator.sh)
bash
#!/bin/bash
# 🚀 SMART STARTUP - Auto kills existing ports!

PORT=${1:-8081}  # ./run-operator.sh 9090 for custom port

echo "🔍 Checking port $PORT..."
if lsof -i :$PORT > /dev/null 2>&1; then
  echo "⚡ Port $PORT in use! Killing..."
  sudo kill -9 $(lsof -t -i:$PORT) 2>/dev/null || true
  sleep 2
fi

echo "✅ Port $PORT FREE! Starting..."

# Deploy + Run
kubectl apply -f config/crd/bases/
kubectl apply -f samples/my-db.yaml

go run ./main.go 2>&1 | grep -v "debug\|stable.example"
Usage:

bash
chmod +x run-operator.sh
./run-operator.sh          # Port 8081
./run-operator.sh 9090     # Custom port
2. Go Code Integration (main.go)
go
// 🔥 AUTO PORT CLEANUP before manager starts
func cleanupPort(port string) error {
    cmd := exec.Command("sh", "-c", 
        fmt.Sprintf("lsof -t -i:%s | xargs -r kill -9", port))
    return cmd.Run()
}

func main() {
    port := os.Getenv("METRICS_PORT")
    if port == "" { port = "8081" }
    
    fmt.Printf("⚡ Cleaning port %s...\n", port)
    cleanupPort(port)
    time.Sleep(2 * time.Second)
    
    // Start manager
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        MetricsBindAddress: ":" + port,
        Port: 9443,
    })
}
3. VS Code Task (.vscode/tasks.json)
json
{
  "version": "2.0.0",
  "tasks": [{
    "label": "🚀 Start Operator (Auto Kill)",
    "type": "shell",
    "command": "bash",
    "args": ["-c", "sudo kill -9 $(lsof -t -i:8081) 2>/dev/null || true && ./run-operator.sh"],
    "group": "build",
    "presentation": { "reveal": "always" }
  }]
}
Run: Ctrl+Shift+P → "Tasks: Run Task" → "Start Operator"

📋 Zero-to-Production Setup
Prerequisites
bash
minikube start
go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
Full Lifecycle
bash
# 1. Generate
make manifests generate

# 2. Smart Start (Auto port kill)
./run-operator.sh

# 3. Verify
curl localhost:8080/metrics | grep database  # 38+ metrics
kubectl get db
Expected Output:

text
🔍 Checking port 8081...
⚡ Port 8081 in use! Killing...
✅ Port 8081 FREE!
✅ DB: laxman-fresh-db postgres 20
INFO	Metrics server:8081 LIVE
📂 Files Structure
text
├── api/v1/database_types.go      # CR Schema
├── controllers/database_controller.go  # Reconcile logic
├── main.go                       # Manager + Port cleanup
├── config/crd/bases/             # CRD YAML
├── samples/my-db.yaml            # Sample CR
├── run-operator.sh               # 🚀 Smart launcher
├── Makefile                      # Kubebuilder
└── .vscode/tasks.json            # Auto-kill task
🎮 Sample CR (samples/my-db.yaml)
text
apiVersion: laxmangodi.com/v1
kind: Database
metadata:
  name: production-db
spec:
  engine: postgres
  storageGb: 50
  replicas: 3
status:
  phase: Provisioning
📊 Metrics (localhost:8081/metrics)
Metric	Meaning	Healthy Value
active_workers	Concurrent reconciles	0 (idle)
reconcile_total{result="success"}	Successful runs	>0
workqueue_depth	Backlog	0
🔧 Production Deployment
bash
# Docker + Deploy
make docker-build IMG=yourrepo/db-operator:v1.0
make deploy IMG=yourrepo/db-operator:v1.0

# Verify
kubectl get db,deployment -l app=database-operator
🧪 Troubleshooting
Error	Fix
Port refused	./run-operator.sh (auto-kills)
Cache sync timeout	Normal - filtered in script
CR not found	kubectl apply -f config/crd/bases/
🎓 Controller Components Deep Dive
Manager 🏢: Orchestrates all components

Scheme 🔑: Knows about your CRD types

Cache 📦: Local copy (fast reads)

Informer 👂: Watches API changes

Reconciler 🔧: Business logic

Health Probes ❤️: /healthz

📄 License
Apache 2.0 © Laxman Godi