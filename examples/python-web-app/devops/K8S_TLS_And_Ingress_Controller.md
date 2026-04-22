Kubernetes TLS & Ingress Controller Mastery Guide
Complete Hands-On: Self-Signed Wildcard Certs → Multi-Domain TLS → Production Ingress

For DevOps Engineers - From Local Development to Production TLS Termination with NGINX Ingress, Traefik, and AWS ALB.

📖 Table of Contents
Why TLS? Business + Technical Advantages

TLS vs Routes: Comparison

Popular Ingress Controllers

Project Structure

Step-by-Step TLS Implementation

Troubleshooting Matrix

Interview Prep Questions

#1 Why TLS? {#why-tls}
Business Reasons
text
🔒 Customer Trust: "https://portfolio.laxman.dev" = Professional
🛡️ Google SEO: HTTPS = +10% search ranking
💳 PCI Compliance: HTTPS required for payments
🌍 Global CDNs: Cloudflare/AWS require TLS
Technical Advantages
text
✅ Encryption: Man-in-middle protection
✅ Browser Trust: No "Not Secure" warnings  
✅ ALB Termination: AWS Load Balancer → HTTPS → HTTP (pods)
✅ Wildcard Certs: *.laxman.dev covers all subdomains
#2 TLS vs Routes Comparison {#tls-vs-routes}
Feature	TLS (HTTPS)	Routes (HTTP)
Security	🔒 Encrypted	⚠️ Plaintext
SEO	✅ Google loves	❌ Penalty
Browser	✅ Green lock	⚠️ "Not Secure"
Performance	📉 +5ms handshake	🚀 Zero overhead
Cert Cost	💰 Free (Let's Encrypt)	✅ Free
Wildcard	🌟 .domain.com	N/A
Winner: TLS for production, HTTP Routes for internal/dev.

#3 Popular Ingress Controllers {#popular-ingress}
Controller	Use Case	TLS Support	Complexity
NGINX	🥇 Most Popular	🌟 Wildcard + ACM	⭐⭐
Traefik	Auto SSL	Let's Encrypt Auto	⭐⭐⭐
AWS ALB	EKS Production	AWS ACM	⭐
HAProxy	High Performance	Manual Certs	⭐⭐⭐⭐
Contour	Envoy Proxy	Advanced	⭐⭐⭐⭐
#4 Project Structure {#project-structure}
text
Laxman-TLS-Ingress-Project/
├── 📁 certs/                    # TLS Certificates
│   ├── wildcard.laxman.dev.crt
│   ├── wildcard.laxman.dev.key
│   └── tls-secret.yaml
├── 📁 nginx-ingress/            # Custom TLS Image
│   ├── Dockerfile
│   ├── nginx.conf
│   └── build.sh
├── 📁 manifests/                # K8s Resources
│   ├── portfolio-complete.yaml
│   ├── ingress-tls.yaml
│   └── multi-domain-tls.yaml
├── 📁 troubleshooting/          # Debug Scripts
│   └── check-tls.sh
└── README.md                    # This file!
#5 Step-by-Step TLS Implementation {#tls-implementation}
Phase 1: Generate Self-Signed Wildcard Certificate
bash
# Step 1: Create CA (Certificate Authority)
openssl req -x509 -newkey rsa:4096 -keyout ca.key -out ca.crt -days 365 \
  -subj "/CN=ca.laxman.dev/O=Laxman DevOps" -nodes

# Step 2: Generate Wildcard Private Key
openssl genrsa -out certs/wildcard.laxman.dev.key 2048

# Step 3: Create Certificate Signing Request (CSR)
openssl req -new -key certs/wildcard.laxman.dev.key -out wildcard.csr \
  -subj "/CN=*.laxman.dev/O=Laxman DevOps"

# Step 4: Sign with CA (Wildcard Cert)
openssl x509 -req -in wildcard.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out certs/wildcard.laxman.dev.crt -days 365 \
  -extensions v3_req -extfile <(cat <<EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
x509_extensions = v3_req

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = *.laxman.dev
DNS.2 = portfolio.laxman.dev
DNS.3 = api.laxman.dev
EOF
)

# Verify
openssl x509 -in certs/wildcard.laxman.dev.crt -text -noout
Output: Subject: CN=*.laxman.dev + DNS:*.laxman.dev, portfolio.laxman.dev

Phase 2: Create Kubernetes TLS Secret
text
# manifests/tls-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: wildcard-tls-secret
type: kubernetes.io/tls
data:
  tls.crt: $(base64 -w0 certs/wildcard.laxman.dev.crt)
  tls.key: $(base64 -w0 certs/wildcard.laxman.dev.key)
bash
kubectl apply -f manifests/tls-secret.yaml
kubectl get secret wildcard-tls-secret -o yaml
Phase 3: Deploy Applications + Ingress
text
# manifests/portfolio-complete.yaml (Your existing app)
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: portfolio-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: portfolio
  template:
    metadata:
      labels:
        app: portfolio
    spec:
      containers:
      - name: nginx
        image: nginx:alpine
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: portfolio-service
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 80
  selector:
    app: portfolio
Phase 4: TLS Ingress (Single Domain)
text
# manifests/ingress-tls.yaml - BASIC TLS
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: portfolio-tls-ingress
  annotations:
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - portfolio.laxman.dev
    secretName: wildcard-tls-secret
  rules:
  - host: portfolio.laxman.dev
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: portfolio-service
            port:
              number: 80
Phase 5: Multi-Domain TLS Ingress (Production)
text
# manifests/multi-domain-tls.yaml - ENTERPRISE GRADE
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: multi-domain-tls
  annotations:
    # NGINX Ingress Specific
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
    # AWS ALB Compatible
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip
spec:
  tls:
  - hosts:
    - portfolio.laxman.dev
    - api.laxman.dev
    - poornima.laxman.dev
    secretName: wildcard-tls-secret  # Single cert for ALL domains
  rules:
  - host: portfolio.laxman.dev
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: portfolio-service
            port:
              number: 80
  - host: api.laxman.dev
    http:
      paths:
      - path: /health
        pathType: Prefix
        backend:
          service:
            name: portfolio-service
            port:
              number: 80
  - host: poornima.laxman.dev
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: portfolio-service
            port:
              number: 80
Phase 6: Deploy Everything
bash
# One-command production deployment
kubectl apply -f manifests/portfolio-complete.yaml
kubectl apply -f manifests/tls-secret.yaml
kubectl apply -f manifests/multi-domain-tls.yaml

# Verify
kubectl get ingress
kubectl get secret wildcard-tls-secret
Phase 7: Test TLS
bash
# Minikube
minikube addons enable ingress
minikube tunnel  # Expose LoadBalancer

# Test HTTPS
curl -k -H "Host: portfolio.laxman.dev" https://$(minikube ip)

# Browser: https://portfolio.laxman.dev (Accept self-signed warning)
#6 Custom NGINX Ingress with TLS {#custom-ngress}
nginx-ingress/Dockerfile
text
FROM nginx/nginx-ingress:3.2.0
# Custom TLS config for wildcard certs
COPY nginx.conf /etc/nginx/nginx.conf
USER root
RUN nginx -t && rm -rf /etc/nginx/conf.d/*
EXPOSE 80 443
USER 101
nginx-ingress/nginx.conf
text
events {
    worker_connections 1024;
}

http {
    upstream portfolio_backend {
        server portfolio-service.default.svc.cluster.local:80;
    }
    
    server {
        listen 443 ssl http2;
        server_name portfolio.laxman.dev;
        
        ssl_certificate /etc/nginx/ssl/wildcard.laxman.dev.crt;
        ssl_certificate_key /etc/nginx/ssl/wildcard.laxman.dev.key;
        
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
        
        location / {
            proxy_pass http://portfolio_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
    }
}
Build & Deploy Custom Ingress
bash
docker build -t laxman/nginx-ingress-tls:1.0 nginx-ingress/
kubectl apply -f nginx-ingress/ingress-deployment.yaml
#7 Troubleshooting Matrix {#troubleshooting}
Issue	Symptoms	Fix
404 on HTTPS	curl: (52) Empty reply	kubectl describe ingress → Check backend
TLS Handshake Fail	curl: (35) SSL connection	openssl x509 -in certs/*.crt -text → Verify SANs
Ingress NotReady	kubectl get pods -n ingress-nginx	minikube addons enable ingress
No External IP	Pending LoadBalancer	minikube tunnel (Minikube)
Cert Mismatch	Browser warning	Add DNS:*.laxman.dev to alt_names
Connection Reset	curl: (56) Recv failure	Check nginx.ingress.kubernetes.io/ssl-redirect
Debug Script troubleshooting/check-tls.sh
bash
#!/bin/bash
echo "🔍 TLS Ingress Debug"
echo "1. Ingress Status:"
kubectl get ingress -o wide
echo -e "\n2. TLS Secret:"
kubectl get secret wildcard-tls-secret -o yaml | grep -A5 tls.crt
echo -e "\n3. Pods:"
kubectl get pods -l app=portfolio
echo -e "\n4. Events:"
kubectl describe ingress multi-domain-tls | grep -A10 Events
#8 Interview Preparation {#interview-prep}
Common Questions + Answers
text
**Q: How does TLS termination work in Ingress?**
A: Ingress Controller (NGINX) terminates TLS using Secret (tls.crt/tls.key). 
Forwards HTTP to backend Service:80. Single wildcard cert (*.laxman.dev) 
handles portfolio.laxman.dev + api.laxman.dev.

**Q: Self-signed vs Let's Encrypt?**
A: Self-signed: Development (openssl req -x509), browser warnings
Let's Encrypt: Production (cert-manager auto-renewal every 90 days)

**Q: Why wildcard certs?**
A: *.laxman.dev covers portfolio.laxman.dev, api.laxman.dev, 
poornima.laxman.dev. Single Secret vs multiple certs.

**Q: Debug TLS 404?**
A: 1) kubectl describe ingress 2) Check Secret base64 decode
3) Verify Service selector matches pod labels 4) kubectl logs -n ingress-nginx
Live Demo Commands
bash
# Show your skills LIVE in interviews
kubectl get ingress -o yaml | grep -A10 tls
openssl x509 -in certs/wildcard.laxman.dev.crt -text | grep DNS
curl -k -H "Host: portfolio.laxman.dev" https://$(minikube ip)
🚀 Quick Start (5 Minutes)
bash
# 1. Generate certs
bash generate-wildcard-certs.sh

# 2. Deploy everything
kubectl apply -f manifests/

# 3. Test TLS
curl -k https://portfolio.laxman.dev/$(minikube ip) -H "Host: portfolio.laxman.dev"

# 4. Cleanup
kubectl delete ingress --all && kubectl delete secret wildcard-tls-secret
