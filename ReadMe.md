# 🚀 Django DevOps & Automation Project
**Author:** Laxman | **Role:** SDET / Automation Test Engineer | **Location:** Bengaluru, India

This repository contains a containerized Django application deployed on AWS EC2. It serves as a portfolio piece to demonstrate proficiency in **CI/CD pipelines**, **Dockerization**, and **Environment Management**.

---

## 🛠 Command Reference Guide
Use these commands to manage the environment and run the application.

### 1. Environment Setup (Local/EC2)
| Task | Command |
| :--- | :--- |
| **Create Virtual Env** | `python3 -m venv .venv` |
| **Activate Env** | `source .venv/bin/activate` |
| **Install Dependencies** | `pip install -r requirements.txt` |
| **Reset Environment** | `rm -rf .venv && python3 -m venv .venv` |

### 2. Django Development
| Task | Command |
| :--- | :--- |
| **Apply Database Changes** | `python3 manage.py migrate` |
| **Create Admin Account** | `python3 manage.py createsuperuser` |
| **Run Dev Server (Local)** | `python3 manage.py runserver` |
| **Run Dev Server (Cloud)** | `python3 manage.py runserver 0.0.0.0:8000` |

### 3. Docker Operations
| Task | Command |
| :--- | :--- |
| **Build Docker Image** | `docker build -t django-devops-app .` |
| **Run Container** | `docker run -p 8000:8000 django-devops-app` |
| **Stop All Containers** | `docker stop $(docker ps -q)` |

---

## 🏗 Project Architecture


1. **Frontend:** Custom HTML/CSS landing page highlighting SDET skills.
2. **Backend:** Django 5.x Framework.
3. **Database:** SQLite (Container-native).
4. **Deployment:** Dockerized Ubuntu 22.04 environment.

## 📁 Key File Descriptions
* `manage.py`: Entry point for Django administrative commands.
* `Dockerfile`: Instructions for building the automated environment.
* `requirements.txt`: List of all Python dependencies.
* `.gitignore`: Prevents sensitive files (like `.venv` and `db.sqlite3`) from being uploaded to GitHub.

---

## ⚠️ Troubleshooting
* **404 Page Not Found:** Ensure you have configured `path('', include('demo.urls'))` in `devops/urls.py`.
* **Port Already in Use:** If Port 8000 is stuck, run `fuser -k 8000/tcp` to clear it.
* **Externally Managed Environment:** This is a Linux safety feature; always activate the `.venv` before running `pip`.

---
*Developed for the purpose of demonstrating Robotics Software Validation and Automation workflows.*
