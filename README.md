
# Go Song Recognition 

This project was done for my final year Enterprise Computing module Continuous Assessment. The project was inspired by the following video - https://www.youtube.com/watch?v=wbj0s1B_O4o and consists of three microservices that work together to form a music recognition application using the audd.io API.



## Usage

To run each microservice, open a terminal in their respective folder and run ```go run``` then the filename. eg:

```
go run Tracks.go
```

Test scripts 1-4 are for the 'tracks' microservice, script 5 is for the 'search' microservice and script 6 is for the 'cooltown' microservice. 

To run the test scripts, open a terminal in the 'test_scripts' folder and run the script you'd like to run using ```sh```.

Please note that:
- For 'cooltown' to work, 'tracks' and 'search' must already be running.
- For 'search' to work, you should replace the audd.io API key on line 14 of 'search.go' with your own.



## Microservices

#### Tracks

The 'tracks' microservice acts as a database for music files, providing Create, List, Read and Delete functions. The database is reset everytime this microservice is run. It listens to port 3000.

#### Search

The 'search' microservice takes in a Base64 encoded audio file and sends it to audd.io to attempt to recognise the song within the provided audio. This listens to port 3001.

#### Cooltown

The 'cooltown' microservice links 'tracks' and 'search' by taking in an imput audio file, sending it to 'search' so a song within it can be recognised, then returns the song stored in the 'tracks' database. This listens to port 3002.
## License

[MIT](https://choosealicense.com/licenses/mit/)

