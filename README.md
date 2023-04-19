# Forum

## 📚 Introduction
This project consists in creating a web forum that allows :

- communication between users
- associating categories to posts
- liking and disliking posts and comments
- filtering posts

Users need to register in order to create posts, add comments and like or dislike posts and comments.

## 👟 Requirements to run

- Docker engine must be installed
- Bash terminal window
- Web browser

## 🏃‍♂️ Running the program

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

The application will start on port 5050. Go to localhost:5050 on your browser.

### Terminal
Make sure you have all the necessary third-party packages installed.

```go
go run .
```
The application will start on port 8080. Go to localhost:8080 on your browser.


## 🧪 Testing the program
Audit can be found [here](https://github.com/01-edu/public/tree/master/subjects/forum/audit)

## ✏️ Notes
The server is written in Go. HTML, CSS and JavaScript are used for frontend. SQLite database is used to store data.

## 🤴 Authors
@Brooklyn_95 \
@kretesaak \
@margus.aid \
@GhanBuriGhan
