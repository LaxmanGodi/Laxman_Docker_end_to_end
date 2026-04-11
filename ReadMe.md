To make your GitHub documentation look professional and highly readable, I’ve formatted your content into **GitHub Flavored Markdown (GFM)**. This includes a clean header, bold tables, and clear "Callout" blocks for troubleshooting.

Copy and paste the block below into your `README.md` file:

---

# 🚀 Django DevOps & Automation Project

**Author:** Laxman  
**Role:** SDET / Automation Test Engineer  
**Location:** Bengaluru, India  

This repository contains a containerized Django application deployed on AWS EC2. It serves as a professional portfolio piece to demonstrate proficiency in **CI/CD pipelines**, **Dockerization**, and **Environment Management** specifically for SDET workflows.

---

## 🛠 Command Reference Guide
Follow these steps to set up the environment, handle dependencies, and run the application.

### 1. System & Dependencies Setup
| Task | Command |
| :--- | :--- |
| **Install tzdata** | `sudo apt update && sudo apt install -y tzdata` |
| **Create Virtual Env** | `python3 -m venv .venv` |
| **Activate Env** | `source .venv/bin/activate` |
| **Upgrade Pip** | `pip install --upgrade pip` |
| **Install Django** | `pip install django` |

### 2. Django Development
| Task | Command |
| :--- | :--- |
| **Apply Migrations** | `python3 manage.py migrate` |
| **Create Admin User** | `python3 manage.py createsuperuser` |
| **Run Server (Local)** | `python3 manage.py runserver` |
| **Run Server (Cloud)** | `python3 manage.py runserver 0.0.0.0:8000` |

### 3. Docker Operations (The "Golden Path")
> **Note:** To avoid Port Conflicts and Build Errors, always follow this sequence inside the `devops/` folder.

| Step | Action | Command |
| :--- | :--- | :--- |
| **1** | **Clean Port 8000** | `sudo fuser -k 8000/tcp` |
| **2** | **Clear Containers** | `docker stop $(docker ps -aq) && docker rm $(docker ps -aq)` |
| **3** | **Build & Tag** | `docker build -t laxmangodi/django-devops-app:v1 .` |
| **4** | **Push to Hub** | `docker push laxmangodi/django-devops-app:v1` |
| **5** | **Run (Detached)** | `docker run -d -p 8000:8000 --name web-app laxmangodi/django-devops-app:v1` |

---

## 🏗 Project Architecture

* **Frontend:** Custom HTML/CSS landing page highlighting SDET skills and Robotics Software Validation.
* **Backend:** Django 5.x Framework with `tzdata` configured for **Asia/Kolkata** timezone.
* **Database:** SQLite (Development) / Migration-ready for PostgreSQL.
* **Environment:** Isolated `.venv` to prevent "externally-managed-environment" conflicts.

---

## 📁 Key File Descriptions
* `manage.py`: Entry point for Django administrative commands.
* `Dockerfile`: Instructions for building the automated Ubuntu-based environment.
* `requirements.txt`: List of all Python dependencies.

> [!IMPORTANT]
> `manage.py` and `Dockerfile` must remain in the **same directory** (`devops/`) for successful builds to ensure the Docker build context is correct.

---

## ⚠️ Troubleshooting & Common Mistakes

### **Address Already in Use**
This happens if a local Django server is running or a "zombie" container is stuck. 
* **Solution:** Always run `sudo fuser -k 8000/tcp` before starting Docker.

### **Dockerfile Path Error**
If you get `"No such file or directory"`, ensure you are inside the `devops/` folder. 
* **Solution:** Docker cannot see files outside its current "Build Context." Run `cd devops/` before building.

### **Broken Pipe / Session Timeout**
If your EC2 disconnects, the `-d` (detached) flag ensures the container keeps running. 
* **Solution:** Use `docker ps` to verify the status after logging back in.

### **404 Page Not Found**
Check `devops/urls.py` for the root path: `path('', include('demo.urls'))`.

---
*Developed for demonstrating DevOps automation and Robotics Software Validation.*
