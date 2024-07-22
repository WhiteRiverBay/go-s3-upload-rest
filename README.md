# HOW TO USE THIS COMPONENET

A REST to AWS S3 uploader


## Build & Run

### Build

```
go build cmd/main.go

```

### Start
```
# REMOVED
ACCESS_KEY=
SECRET_KEY=
BUCKET=
REGION=

LOG_FILE=logs/$(date +"%Y-%m-%d_%H-%M-%S").log

nohup ./main -access $ACCESS_KEY \
    -secret $SECRET_KEY -bucket $BUCKET -region $REGION 2>&1 > $LOG_FILE &

echo $! > pid
```

### Test

```
#!/bin/bash

# This script is used to post the data to the server
# The data is stored in the file data.txt

# The server address
server="http://localhost:8080/upload"

# The data file
datafile="bitcoin.png"

# Post the data to the server, multipart/form-data is used, field name is "file"
curl -X POST -F "file=@$datafile" $server
```
