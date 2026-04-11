Project Documentation: DevOps & SDET Automation
Prepared by: Laxman, Automation Test Engineer

Location: Bengaluru, India

Date: April 2026

1. Project Overview
This project serves as a professional demonstration of integrating Django web development with modern DevOps practices. It covers the full lifecycle of an application, from local development in a virtual environment to containerization using Docker on a cloud-based EC2 instance.

Primary Focus: SDET, Automation Frameworks, and Robotics Software Validation.

Infrastructure: AWS EC2 (Ubuntu 22.04).

Key Technologies: Python, Django, Docker, Bash, and Linux Administration.

2. Command Reference Guide
Environment & Setup
Action	Command
Activate Virtual Env	source venv1/bin/activate
Install Django	pip install django
Upgrade Pip	pip install --upgrade pip
Django Operations
Action	Command
Database Migration	python3 manage.py migrate
Create Superuser	python3 manage.py createsuperuser
Run Server (Public)	python3 manage.py runserver 0.0.0.0:8000
Kill Port 8000	fuser -k 8000/tcp
Docker Operations
Action	Command
Build Image	docker build -t devops-app .
Run Container	docker run -p 8000:8000 devops-app
List Containers	docker ps -a
3. Dockerfile Architecture
The project utilizes a custom Ubuntu-based Dockerfile designed for maximum control over the environment.

Base Image: ubuntu:22.04

Layering: System tools are installed first, followed by the project code to optimize build caching.

Entrypoint: Uses a bash script to ensure venv1 is activated for every command execution.

Port Mapping: Exposed on 8000 for external accessibility.

4. Troubleshooting & Best Practices
Common Issues
Connection Timed Out: Ensure the AWS Security Group allows Inbound TCP traffic on Port 8000.

Permission Denied: Use ls -l to check file ownership. If files are owned by root, use sudo chown to regain control.

Vim Editing: If backspace/arrows fail in vi, install the full version using sudo apt install vim.

Maintenance Tips
Always use a virtual environment (venv) to avoid corrupting system-level Python packages.

Always keep ALLOWED_HOSTS updated in settings.py when changing instance IPs.

Always document new changes in the README.md for team collaboration.

5. Author Contact Information
Role: Automation Test Engineer / SDET

Expertise: Selenium, Java, REST Assured, NVIDIA Isaac Sim, Docker.

GitHub/LinkedIn: [Insert your profile links here]
