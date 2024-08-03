tag=0.2
docker build -t go-s3-upload-rest .
docker tag go-s3-upload-rest whiteriverbay/go-s3-upload-rest:$tag
docker push whiteriverbay/go-s3-upload-rest:$tag