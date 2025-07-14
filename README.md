# Forum

## 📚 Introduction
This project consists in creating a web forum that allows :

- communication between users
- associating categories to posts
- liking and disliking posts and comments
- filtering posts

Users need to register in order to create posts, add comments and like or dislike posts and comments.

We started working on the project at the beginning of March and we use Bootstrap in the frontend (not for everything, but for some parts). Using frontend libararies and frameworks in this task has since been forbidden. We found out about this when we had everything already done and we were ready for the audit and Karl said that it is fine.

## 👟 Requirements to run

- Docker engine must be installed
- Bash terminal window
- Web browser

## 🏃‍♂️ Running the program

### ✈️ Terminal
Make sure you have all the necessary third-party packages installed.

You can run the code with:
```go
go run .
```
The application will start on **port 8080**. Go to localhost:8080 on your browser.

To open the database, use:
```bash
 sqlite3 ./db/forum.db
```

To see all the users, open the database and use:
```sqlite3
 SELECT * FROM user;

```
To see all the posts, use:
```sqlite3
 SELECT * FROM post;
```

To see all the comments, use:
```sqlite3
 SELECT * FROM comment;
```
To exit the terminal, use `Ctrl + D`


### 🐋 Docker

1. Build the image:

```bash
docker build -t forum -f Dockerfile .
```

2. Run the image:

```bash
docker container run -p 5050:8080 --detach --name forum-container forum
```

Or use a bash script to build and run the image:

```bash
sh dockersetup.sh
```

The application will start on **port 5050**. Go to localhost:5050 on your browser.

**The local database and the docker container one are different - they won't include anything you add to the other one!**

## 🧪 Testing the program
Audit can be found [here](https://github.com/01-edu/public/tree/master/subjects/forum/audit)

## ✏️ Notes
The server is written in Go. HTML, CSS and JavaScript are used for frontend. SQLite database is used to store data.

## 🤴 Authors
@Brooklyn_95 \
@kretesaak \
@margus.aid \
@GhanBuriGhan




# Project101 – CI/CD Deployment to K3s with GitHub Actions

This repository contains a Go application that is automatically built, tested, containerized, and deployed to a K3s Kubernetes cluster using GitHub Actions.

## Workflow Overview

On every push to the `master` branch, GitHub Actions performs the following steps:

1. Clone the repository
2. Set up the Go and Java environments
3. Run tests and generate a code coverage report
4. Build the Go binary
5. Build and push a Docker image to Docker Hub
6. Deploy the application to K3s:
   - Creates or updates a Kubernetes Deployment
   - Uses a dynamically assigned NodePort for the service
   - Sets up an Ingress resource for external access

## Docker Image

The Docker image is tagged `latest` and pushed to:

docker.io/maxtech470/project101:latest

markdown
Copy
Edit

You can update the `COMMIT_TAG` variable in the workflow to use versioned tags if desired.

## Kubernetes Configuration

- Namespace: `forum`
- Replicas: 6
- Container port: 8080
- Service port: 80
- NodePort: dynamically assigned between 30000–32767
- Ingress path: `/`
- Ingress controller: expects `nginx` class

## Required GitHub Secrets

To enable this workflow, the following secrets must be configured in your GitHub repository:

| Secret Name           | Description                                       |
|-----------------------|---------------------------------------------------|
| `DOCKER_USERNAME`     | Docker Hub username                               |
| `DOCKER_PASSWORD`     | Docker Hub password or personal access token      |
| `KUBECONFIG`          | Base64-encoded kubeconfig for accessing the K3s cluster |
| `MY_SECRET_CODE`      | Application-level secret passed as an environment variable |

To encode your kubeconfig:
base64 -w 0 ~/.kube/config

shell
Copy
Edit

## Accessing the Application

After deployment, the workflow prints the dynamically selected NodePort. Use that to access the application:

http://<node-ip>:<nodePort>

vbnet
Copy
Edit

If you're using Ingress and DNS, point a domain (or your local `/etc/hosts`) to your K3s node IP and access the app at:

http://project101.local/

markdown
Copy
Edit

## Monitoring Integration (Prometheus and Grafana)

1. Ensure Prometheus is scraping your application's `/metrics` endpoint.
2. Add the service name to Prometheus' `scrape_configs`.
3. In Grafana:
   - Add Prometheus as a data source using this URL:
     ```
     http://prometheus-server.monitoring.svc.cluster.local:80
     ```
   - Create dashboards using Prometheus queries.

## Example API Usage

Test the application after deployment using curl:

curl http://<node-ip>:<nodePort>/

css
Copy
Edit

To hit a specific endpoint with a header:

curl -H "Authorization: Bearer <token>" http://<node-ip>:<nodePort>/api/health

vbnet
Copy
Edit

## Migrating to Helm (Optional)

To modularize the deployment using Helm:

1. Create a chart:
helm create project101

sql
Copy
Edit
2. Move your Kubernetes manifests into the `templates/` folder and replace hardcoded values with variables.
3. Update your GitHub Action to install using Helm:
helm upgrade --install project101 ./chart/project101 -n forum --set image.tag=latest

shell
Copy
Edit

## Manual Debugging

You can apply the manifests locally for troubleshooting:

kubectl apply -f <generated-manifest>.yaml

css
Copy
Edit

Or forward a local port to the service:

kubectl port-forward svc/project101-service -n forum 8080:80

Author:
@maxtech001 (Jude Ifeanyi Eze)
