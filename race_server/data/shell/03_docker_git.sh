# Topic: Docker and Git
docker build -t myapp:latest .
docker run -d \
    -p 8080:8080 \
    -v ./data:/app/data \
    --name myapp \
    myapp:latest

git add -A
git commit -m "deploy: update config"
git push origin main
git tag -a v1.0.0 -m "release"
git push --tags
