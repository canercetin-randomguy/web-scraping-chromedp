# Damacansu Broken Link Finder

Proof of concept. You need a local MySQL or MariaDB database running on 127.0.0.1:3306 with sufficient permissions. (you can change IP, Port, DB Username&Password in
sqlpkg/types.go).

## Usage

Head to cmd folder, run
``
go run main.go
``
and then head to the localhost:7171/signin. Port can be changed in line 19 main.go.

## API Documentation

localhost:7171 also carries an API router group. You need to send ``raw data`` from PostMan.
Users must grab a secret key from website's homepage to authenticate.
````go
/v1/api/auth -> POST -> {"username":"admin","password":"admin","email":"cartcurt@gmail.com","secretkey":"topsecret"}
on success, returns
JSON -> {"status":"success","username":"admin","auth":"asdasdasd", "message":  "Successfully authenticated, please use the token for future requests."}
on fail, returns
JSON -> {"status":"fail","message":"Invalid secret key."}
````
Returns an AUTH key. Necessary for other API calls.
````go
/v1/api/crawl -> POST -> {"username":"admin","authkey":"asdasdasd","maxdepth":"2","mainlink":"https://www.google.com"}
````
API is broken. Dont know why. It was working minutes before interview and it is down now. Guess that is my luck.

## Security?  
-Cookie check for auth token on whole sites (except signin and signup). Add a new endpoint to V1 router to use this security measure.

-SQL Injection prevention by only using Query or Prep functions.

-Secret Key for API.

-16-cost password hashing while storing the password in database.

-Callback functions are private. Whenever you click a button or form in the website, it will go through a router that is running on loopback adapter, so only 127.0.0.1, server can handle the callbacks.

-File paths in storage contains username's, such as https://52d4-194-27-73-85.eu.ngrok.io/v1/storage/canercetin/result_canercetin_20230109_1.csv has canercetin in the path for canercetin client. Only canercetin client can access it, cookie checks are used in the process.


Oh also, we have a logger. Kek.
