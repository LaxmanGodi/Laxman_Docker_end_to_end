# Use a slim Python base for a smaller image size
FROM python:3.9-slim

# Set working directory
WORKDIR /app

# Install system-level dependencies for SSL/Cryptography
# These are required to build pyOpenSSL correctly
RUN apt-get update && apt-get install -y \
    gcc \
    libssl-dev \
    libffi-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
# We include both standard Django and the SSL-capable extensions
RUN pip install --no-cache-dir \
    django \
    django-extensions \
    pyOpenSSL \
    Werkzeug

# Copy your Django project files
COPY . .

# We leave the CMD empty so you can override it in your 
# deploy.yaml for each different test (Offload vs Bridge)
