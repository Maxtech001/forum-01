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

### Terminal
Make sure you have all the necessary third-party packages installed.

Before you can run the code in terminal without docker image, you need to create a db:
```bash
 mkdir db && sh dbcreate.sh
```
And then you can run the code with:
```go
go run .
```
The application will start on port 8080. Go to localhost:8080 on your browser.

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
