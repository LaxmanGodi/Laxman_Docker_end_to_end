# 🚀 Django DevOps & Automation Project
**Author:** Laxman | **Role:** SDET / Automation Test Engineer | **Location:** Bengaluru, India

This repository contains a containerized Django application deployed on AWS EC2. It serves as a professional portfolio piece to demonstrate proficiency in **CI/CD pipelines**, **Dockerization**, and **Environment Management** specifically for SDET workflows.

---

## 🛠 Command Reference Guide
Follow these steps to set up the environment, handle dependencies, and run the application.

### 1. System & Dependencies Setup
Before initializing the project, ensure the system handles timezones and Python environments correctly.

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

### 3. Docker Operations
| Task | Command |
| :--- | :--- |
| **Build Docker Image** | `docker build -t django-devops-app .` |
| **Run Container** | `docker run -p 8000:8000 django-devops-app` |
| **Stop All Containers** | `docker stop $(docker ps -q)` |

---

## 🏗 Project Architecture

1. **Frontend:** Custom HTML/CSS landing page highlighting SDET skills and Robotics Software Validation.
2. **Backend:** Django 5.x Framework with `tzdata` configured for **Asia/Kolkata** timezone.
3. **Database:** SQLite (Development) / Migration-ready for PostgreSQL.
4. **Environment:** Isolated `.venv` to prevent "externally-managed-environment" conflicts.

## 📁 Key File Descriptions
* `manage.py`: Entry point for Django administrative commands.
* `Dockerfile`: Instructions for building the automated Ubuntu-based environment.
* `requirements.txt`: List of all Python dependencies (Django, etc.).
* `.gitignore`: Configured to ignore `.venv/`, `__pycache__/`, and `db.sqlite3`.

---

## ⚠️ Troubleshooting
* **404 Page Not Found:** Check `devops/urls.py` for the root path: `path('', include('demo.urls'))`.
* **Required File Not Found:** If `.venv` is corrupted, run `rm -rf .venv` and recreate it.
* **Port Already in Use:** Use `fuser -k 8000/tcp` to kill stuck processes on Port 8000.
* **Timezone Issues:** Ensure `tzdata` is installed to avoid UTC/IST offset errors in logs.

---
*Developed for demonstrating DevOps automation and Robotics Software Validation.*
