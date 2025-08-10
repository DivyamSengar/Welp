# Welp: A High-Performance Yelp-Like Microservices Platform

Welcome to **Welp**, a scalable, high-performance microservices application built to emulate the core functionalities of Yelp. This project showcases the design, deployment, and optimization of a cloud-native application using a modern tech stack. The system is architected with distinct services for restaurant details, reviews, and reservations, all communicating via gRPC.

This repository contains the full implementation, documentation, and scripts required to build, deploy, and test the application.

---

## 🚀 Key Features & Achievements

* **Microservices Architecture**: Built with independent **Detail**, **Review**, and **Reservation** services using **Golang** and **gRPC** for efficient, low-latency communication.
* **Scalable Deployment**: Deployed on **Google Cloud Platform (GCP)** using **Docker** containers and orchestrated with **Kubernetes**, featuring automated CI/CD pipelines.
* **Bottleneck Analysis & Optimization**: Conducted in-depth latency-throughput analysis to identify and resolve frontend bottlenecks, **improving maximum throughput by 4.2x** through strategic Kubernetes pod scaling.
* **Advanced Caching**: Integrated **LRU** and **LFU** caching policies, which **reduced storage-layer access latency by 85%** under Zipfian workloads and achieved **sub-10ms cache-hit response times**.
* **Low-Level Performance Tuning**: Engineered a custom C-based thread-scheduling benchmark to optimize for cache locality on GCP, **reducing cross-core memory latency by 10.8ms** and minimizing cache misses through same-core co-location.

---

## 🛠️ Technology Stack

* **Languages**: Golang, Python, C
* **Containerization & Orchestration**: Docker, Kubernetes
* **Cloud Platform**: Google Cloud Platform (GCP)
* **Communication**: gRPC, Protobuf
* **Benchmarking**: `wrk2`

---

## 🏛️ Architecture & Implementation

The Welp platform is built on a foundation of several key components that work together to deliver a robust and scalable service.

* **Core Microservices**: The application's backend is composed of three primary services:
    * **Detail Service**: Manages foundational information about restaurants, such as location, style, and capacity.
    * **Review Service**: Handles the creation and retrieval of user reviews and ratings for restaurants.
    * **Reservation Service**: Manages user reservations and tracks restaurant popularity based on reservation frequency.

* **API & Communication**: A **Frontend** service acts as the gateway, exposing an HTTP API to the outside world and routing requests to the appropriate backend services using **gRPC**. Data structures are defined using **Protocol Buffers** for efficient serialization.

* **Performance Analysis and Scaling**: The system was rigorously benchmarked using `wrk2` to model user load. This allowed for a detailed analysis of latency versus throughput, which led to the identification of the frontend as the primary bottleneck. By strategically scaling the frontend pods in Kubernetes, the system's maximum throughput was increased by a factor of 4.2.

* **Caching & Storage Layer**: To improve read performance and reduce database load, a look-aside caching layer was introduced. This layer supports multiple eviction policies, including **LRU (Least Recently Used)** and **LFU (Least Frequently Used)**. Under realistic Zipfian workloads, the cache dramatically reduced latency for frequently accessed data. The storage backend is emulated to simulate real-world database performance.

---

## 🔧 Getting Started

### Prerequisites

* A Google Cloud Platform project with a configured VM instance.
* Docker installed and configured.
* Go programming language environment.
* `kubectl` command-line tool.

### Installation & Deployment

1.  **Clone the repository:**
    ```bash
    git clone <your-repo-url>
    cd <repository-directory>
    ```

2.  **Set up the Kubernetes environment:**
    * Run the setup scripts provided in the `scripts/` directory to configure your Kubernetes cluster on the GCP VM.

3.  **Build and Push the Docker Image:**
    * Create a private repository on Docker Hub.
    * Use the provided build script, making sure to log in to Docker first.
    ```bash
    sudo docker login -u <your-docker-username>
    sudo bash scripts/build_image.sh -u <your-docker-username> -t <tag>
    ```

4.  **Deploy the Application:**
    * Update the Kubernetes manifest files in the `manifests/` directory to point to your Docker image.
    * Create a Kubernetes secret for your private Docker repository.
    ```bash
    kubectl create secret docker-registry regcred --docker-server=<your-docker-server> --docker-username=<your-docker-username> --docker-password=<your-docker-password>
    ```
    * Apply the manifests to deploy Welp:
    ```bash
    kubectl apply -f manifests/
    ```

Your Welp microservices platform should now be running on your Kubernetes cluster!
