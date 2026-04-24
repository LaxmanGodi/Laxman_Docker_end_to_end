🔐 EKS RBAC Multi-User Mastery Guide
Production IAM → Kubernetes RBAC with Access Entries (eksctl + AWS CLI)

For DevOps/SRE Engineers - Complete multi-tenant EKS setup: Dev/UAT/Prod namespaces with granular permissions. 100% production-ready.

📖 Table of Contents
Architecture Overview

Prerequisites

Step-by-Step Implementation

Verification & Testing

Troubleshooting

Cleanup

Interview Questions

#1 Architecture Overview {#architecture}
text
IAM Users (dev-user, staging-user) 
    ↓ aws eks create-access-entry
EKS Access Entries (username: dev-user)
    ↓ ClusterRole + RoleBinding
Kubernetes RBAC (namespace-scoped permissions)
Flow: AWS IAM → EKS Access Entry → K8s Username → RBAC RoleBin

Production IAM → Kubernetes RBAC with Access Entries (eksctl + AWS CLI)

For DevOps/SRE Engineers - Complete multi-tenant EKS setup: Dev/UAT/Prod namespaces with granular permissions. 100% production-ready.

📖 Table of Contents
Architecture Overview

Prerequisites

Step-by-Step Implementation

Verification & Testing

Troubleshooting

Cleanup

Interview Questions

#1 Architecture Overview {#architecture}
text
IAM Users (dev-user, staging-user) 
    ↓ aws eks create-access-entry
EKS Access Entries (username: dev-user)
    ↓ ClusterRole + RoleBinding
Kubernetes RBAC (namespace-scoped permissions)
Flow: AWS IAM → EKS Access Entry → K8s Username → RBAC RoleBinding

#2 Prerequisites {#prerequisites}
bash
# Local Ubuntu 22.04+
aws --version          # v2.15+
kubectl --version      # v1.29+
eksctl version         # v0.175+
ansible --version      # v2.16+

# AWS Permissions
aws sts get-caller-identity  # Admin or EKSFullAccess
Project Structure:

text
eks-rbac-multiuser/
├── ansible/
│   └── setup-iam-users.yml
├── k8s/
│   ├── master-role.yaml
│   └── master-rolebinding.yaml
└── README.md
#3 Step-by-Step Implementation {#implementation}
Step 1: Create IAM Users (Ansible)
text
# ansible/setup-iam-users.yml
---
- name: Create EKS IAM Users
  hosts: localhost
  tasks:
    - name: Create dev-user
      iam_user:
        name: dev-user
        state: present
        access_key_state: present
    
    - name: Attach EKS read policy to dev-user
      iam_role_policy:
        role: dev-user
        policy_arn: arn:aws:iam::aws:policy/AmazonEKSClusterPolicy
        state: present
bash
ansible-playbook ansible/setup-iam-users.yml
aws iam list-users | grep dev-user  # Verify
Step 2: Provision EKS Cluster (eksctl)
bash
eksctl create cluster \
  --name laxman-cluster \
  --region us-east-1 \
  --version 1.29 \
  --nodegroup-name standard-nodes \
  --node-type t3.medium \
  --nodes 2 \
  --nodes-min 1 \
  --nodes-max 4 \
  --managed

# Update kubeconfig
aws eks update-kubeconfig --region us-east-1 --name laxman-cluster
kubectl get nodes  # Verify cluster ready (~12 mins)
Step 3: Define Master RBAC (Reusable ClusterRole)
text
# k8s/master-role.yaml - REUSABLE across namespaces
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: environment-developer-role
rules:
# Core resources (all namespaces)
- apiGroups: [""]
  resources: ["pods", "pods/log"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "replicasets"]
  verbs: ["get", "list", "create", "update", "patch", "delete"]
- apiGroups: [""]
  resources: ["services", "configmaps"]
  verbs: ["get", "list", "create"]

# Read-only for secrets (security)
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get"]
Step 4: Namespace-Scoped RoleBinding
text
# k8s/master-rolebinding.yaml - BIND USER TO NAMESPACE
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: dev-user-python-binding
  namespace: python-app-dev  # Dev namespace
subjects:
- kind: User
  name: dev-user            # EKS Access Entry username
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: environment-developer-role
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: staging-user-python-binding
  namespace: python-app-uat    # UAT namespace
subjects:
- kind: User
  name: staging-user
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: environment-developer-role
  apiGroup: rbac.authorization.k8s.io
Step 5: IAM → EKS Handshake (Access Entries)
bash
# Bridge AWS IAM → Kubernetes username
aws eks create-access-entry \
  --cluster-name laxman-cluster \
  --principal-arn arn:aws:iam::YOUR-ACCOUNT: user/dev-user \
  --username dev-user

aws eks create-access-entry \
  --cluster-name laxman-cluster \
  --principal-arn arn:aws:iam::YOUR-ACCOUNT: user/staging-user \
  --username staging-user

# Associate cluster-wide policy (optional)
aws eks associate-access-policy \
  --cluster-name laxman-cluster \
  --principal-arn arn:aws:iam::YOUR-ACCOUNT:user/dev-user \
  --policy-arn arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy \
  --access-scope type=cluster
Step 6: Deploy & Apply RBAC
bash
# Create namespaces
kubectl create namespace python-app-dev python-app-uat

# Apply RBAC
kubectl apply -f k8s/master-role.yaml
kubectl apply -f k8s/master-rolebinding.yaml

# Test permissions (SDET Style)
kubectl auth can-i create deployment --namespace python-app-dev --as dev-user
# Returns: yes

kubectl auth can-i create deployment --namespace python-app-uat --as dev-user
# Returns: no (cross-namespace blocked)
#4 Verification & Testing {#verification}
User-Specific Kubeconfig
bash
# dev-user context
kubectl config set-context dev-user-context \
  --cluster=laxman-cluster.us-east-1.eksctl.io \
  --user=arn:aws:iam::YOUR-ACCOUNT:user/dev-user

# Test as dev-user
kubectl --context dev-user-context get pods -n python-app-dev
kubectl --context dev-user-context create deployment nginx --image=nginx -n python-app-dev
Expected:

text
✅ dev-user: Can create deployments in python-app-dev ONLY
✅ staging-user: Can create deployments in python-app-uat ONLY
❌ Cross-namespace: Permission denied
#5 Troubleshooting {#troubleshooting}
Error	Cause	Fix
User "arn:..." cannot create	No Access Entry	aws eks create-access-entry
Unable to connect	Wrong kubeconfig	aws eks update-kubeconfig
No resources found	Namespace mismatch	Check kubectl get ns
RoleBinding denied	Wrong username	--username dev-user in access entry
Debug Script:

bash
kubectl auth can-i --list --as dev-user
aws eks list-access-entries --cluster-name laxman-cluster
#6 Cleanup {#cleanup}
bash
# 1. Delete K8s resources
kubectl delete -f k8s/master-rolebinding.yaml
kubectl delete -f k8s/master-role.yaml
kubectl delete namespace python-app-dev python-app-uat

# 2. Remove Access Entries
aws eks delete-access-entry \
  --cluster-name laxman-cluster \
  --principal-arn arn:aws:iam::YOUR-ACCT:user/dev-user

# 3. Delete EKS (~15 mins, $0.10/hour savings)
eksctl delete cluster --name laxman-cluster --region us-east-1

# 4. IAM Cleanup
ansible-playbook ansible/setup-iam-users.yml -e state=absent
#7 Interview Preparation {#interview}
text
**Q: How does EKS Access Entries work?**
A: IAM User → aws eks create-access-entry --username dev-user → 
K8s sees "dev-user" → ClusterRoleBinding maps to namespace permissions.

**Q: ClusterRole vs RoleBinding?**
A: ClusterRole = cluster-wide definitions
RoleBinding = namespace-scoped assignment
"developer-role" reusable across Dev/UAT/Prod.

**Q: Map IAM to K8s username?**
A: aws eks create-access-entry --principal-arn IAM-ARN --username k8s-username

**Q: Test RBAC without switching contexts?**
A: kubectl auth can-i create deployment --namespace dev --as dev-user
🚀 Quick Start (20 Minutes)
bash
# Clone & Run
git clone <your-repo>
cd eks-rbac-multiuser

ansible-playbook ansible/setup-iam-users.yml
eksctl create cluster --name laxman-cluster --region us-east-1 --nodes 2
kubectl apply -f k8s/

# Test
kubectl auth can-i create pods --namespace python-app-dev --as dev-user
# yes ✅

# Cleanup
eksctl delete cluster laxman-cluster
Production-Ready | Multi-Tenant | Interview Gold 🎯

Laxman Kumar | DevOps Engineer
EKS RBAC that actually works 🚀