# Damacansu Broken Link Finder

Proof of concept. You need a local database running on 127.0.0.1:3306 (you can change IP and Port in
sqlpkg/types.go).

## Usage

Head to cmd folder, run
``
go run main.go
``
and then head to the localhost:7171. Port can be changed in line 19 main.go.

## TODO
No order. Do whatever is easiest. Start from logging tbh.

~~1- Integrate the broken link finder with the website.~~

2- After integrating, let users download .xlsx, .csv, .txt of the results, etc. This is probably hardest part. Yeah this is hardest part. UUUUUUUUUUUURGH.

2.5 - Make a small dashboard for the results, slap a refresh button. Go-app may be used for WASM, 
but it would make things far, far more and unnecessarily complicated. So just refresh the page whatever,
I am not a frontend dude.

~~2.75 - dont use js~~ too late, there is small scripts in page HTMLs. sorry.

~~3- Make a logger package, too many errors fiddling around. Make a new folder under cmd called logs and 
drop everyting there. Zap may be used, check this out. https://github.com/uber-go/zap/issues/294#issuecomment-280064854~~

~~3.25- Just dont forget to log.~~

4- Make more website pages, such as upgrade, ~~login~~, ~~register~~, etc.

~~5- go play one more match at dota.~~ this was the most mentally challenging todo.

6- Slap the whole thing to Docker.

7- Make an API based thing, so I can use it from Postman. Like make a POST endpoint for retrieving auth token with
user credentials.Then make a POST endpoint for submitting a request.Then make a GET endpoint for retrieving the results. 

8- For the love of god, fix JSONs, let the sign-up, or any page return sensible JSONs. Do this especially if you
want to do Step 7. Currently only something like {status:failure} is returned, I can do better than that.

9- Write unit tests for database functions with some sort of mock library. 
