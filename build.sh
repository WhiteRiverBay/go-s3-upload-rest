docker build -t go-s3-upload-rest .
docker tag go-s3-upload-rest whiteriverbay/go-s3-upload-rest:0.1
docker push whiteriverbay/go-s3-upload-rest:0.1