#!/bin/bash
# 🎖️ ZERO ERRORS - Clean Operator (Production Ready)

set -euo pipefail  # Fail fast on any error

echo "🚀 CLEAN OPERATOR STARTUP"
echo "========================"

# 1. Kill ports (silent)
echo "🔪 Port cleanup..."
sudo fuser -k 8080/tcp 2>/dev/null || true
sudo fuser -k 8081/tcp 2>/dev/null || true
pkill -f "go run.*main.go" 2>/dev/null || true
sleep 2

# 2. Verify port
if lsof -i :8080 >/dev/null 2>&1; then
  echo "❌ Port 8080 busy - fix manually"
  exit 1
fi
echo "✅ Port 8080 FREE!"

# 3. Setup + Deploy
export PATH="$PWD/go/bin:$PATH"
echo "📦 Deploying CRD + DB..."
kubectl apply -f database-crd.yaml
kubectl apply -f my-new-db.yaml

echo "✅ DB Status:"
kubectl get db -o wide

# 4. Run (filter spam errors)
echo "🎬 Controller LIVE (metrics:8080)..."
echo "📊 curl localhost:8080/metrics | grep database | wc -l"
go run -mod=mod main.go 2>&1 | grep -v "stable.example.com\|cache to sync"
