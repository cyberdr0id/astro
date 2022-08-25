## Service receives an APOD image and saves it to AWS Simple storage service.

Before start you must set up next environment variables.
```
DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME

AWS_BUCKET
AWS_REGION
AWS_ACCESS_KEY
AWS_ACCESS_KEY_ID  

APP_PORT
```

# Database
To create database with neccessary tables run `make db`  
  
**Make sure that *Docker* is installed**

# Application
To run application type `make run`.

# Docker image
To build Docker image run `make image`.