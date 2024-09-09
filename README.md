# Kubernetes System Monitoring with Helm and Autoscaling

This project implements a system monitoring tool using Docker and Kubernetes to track the health of web servers. The application allows users to add and monitor web servers via API endpoints, and stores server statuses in a database. Horizontal Pod Autoscaling (HPA) is configured for automatic scaling, while liveness, readiness, and startup probes ensure the application’s health.

## Key Contributions

- **Developed System Monitoring API**: Implemented endpoints for adding servers and retrieving health statuses.
- **Containerization & Kubernetes Deployment**: Used Docker to containerize the application and deployed it to a Kubernetes cluster.
- **Helm & HPA**: Configured Horizontal Pod Autoscaling for efficient scaling and used Helm charts for streamlined Kubernetes resource management.
- **Health Probes**: Implemented liveness, readiness, and startup probes to monitor and maintain application performance.
- **StatefulSet & Database**: Managed database services using Kubernetes StatefulSets for both reading and writing operations.

## Features

- Periodic server health checks with status tracking.
- Autoscaling based on resource usage.
- Deployment automation using Helm charts.
- Metrics and logging integration for monitoring.
