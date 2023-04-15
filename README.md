## 🐋 Docker
## 🏃‍♂️ Running and testing the program

1. Build the image:

```bash
docker build -t forum -f Dockerfile .
```

2. Check if image is built:

```bash
docker images
```

3. Run the image:

```bash
docker container run -p 5050:8080 --detach --name forum-container forum
```

4. Check for running containers:

```bash
docker ps -a
```

5. See inside the container:

```bash
docker exec -it forum-container /bin/bash
```

6. And to see the file system:

```bash
ls -l
```


Or use a bash script to build and run the image:

```bash
sh dockersetup.sh
```

The application will start on port 5050. Go to localhost:5050 on your browser.