# Use Ubuntu as base
FROM ubuntu:22.04

# Prevent interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

# Set working directory inside the container
WORKDIR /app

# Install system dependencies, python, and tzdata
RUN apt-get update && apt-get install -y \
    apt-utils \
    python3 \
    python3-pip \
    python3-venv \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# Set Timezone to Bengaluru
RUN ln -fs /usr/share/zoneinfo/Asia/Kolkata /etc/localtime && \
    dpkg-reconfigure --frontend noninteractive tzdata

# Copy requirements first to leverage Docker cache
COPY requirements.txt /app/

# Create virtual environment and install dependencies
RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
RUN pip install --upgrade pip && pip install -r requirements.txt

# Copy the actual project files
# Since we are running build from the same folder as manage.py
COPY manage.py /app/
COPY devops /app/devops
COPY demo /app/demo 

# Expose the Django port
EXPOSE 8000

# Run the application
CMD ["python3", "manage.py", "runserver", "0.0.0.0:8000"]
