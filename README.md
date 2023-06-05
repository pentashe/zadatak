# FER - GO 2023

/jmbag, /sum, /multiply, /fetch realizirani po specifikaciji

### /0036391234 (vas jmbag)
  - throws *404: File not yet created*
  - Write to file: Send json in POST body, e.g. {"ime": "Ivan", "prezime": "Horvat", "JMBAG": "0036391234"}
  - Read: response to GET is of type text/plain

### /saxpy
  - Vector size: **const N uint**
  - POST data to server in form of JSON:
```
        {
           "a":{
             "0": 2
           },
           "x": {
             "0": 1,
             "1": 2,
             "2": -3
           },
           "y": {
             "0": 1,
             "1": 1,
             "2": 1
           }
        }
```
  - throws **400: JSON Format Error** and **400: X or Y index greater than N**
  - Parallelism achieved using **WaitGroup** and **RWMutex**
  - SAXPY calculated with **gonum.org/v1/gonum** module
