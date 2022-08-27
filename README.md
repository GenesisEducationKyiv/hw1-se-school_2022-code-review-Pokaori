# Bitcoin Service
This is test task for Software Engineering School. This website allows you to subscribe your email to get information about the BTC to UAH exchange rate.
## Technologies
The following technologies were used for this project:
- Go 1.13
- gorilla/mux
- Gomail
- Docker
## Project Structure
The project includes the following packages (modules):
- Controllers - process requests.
- Config - store settings for our application.
- Models - store classes for storages (EmailJsonStorage).
- Routes - register new paths and use controllers to process requests.
- Utils - classes for some additional functionality:
  - BitcoinConverterCoingate - get Bitcoin rate from Coingate API.
  - EmailBTCtoUAHNotifier - send emails to all subscribed emails. It uses goroutines and channels to do it asynchronously. 
 Also project use logs to write error and some notifications for easier debugging. 
 Docker create lightweight image as it uses multistage build. Once we have built our project it moves binary to raw alpine image. 
 ## Run
1. Before running you need to create .env file. Copy the .env.example file and replace the <b>secret</b> variables with your own values.
2. Build project with docker
```console
$ docker build --no-cache -t bitcoin-service .
```
3. Run container
```console
$ docker  run -p 8000:8000 bitcoin-service
```
## Possible improvements
It is only MVP and I have some ideas how I can improve it:
- Store rate value in a cache.
- Read emails piece by piece not to load memory too much.
- Add rotate log files.
- If it is really important to send email to every subscriber every time, we can continue sending asynchronously while we have errors. 
- Better validation.
- Unit tests.
