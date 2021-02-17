# ltester
A command-line tool for load testing

Run the tool with:
```sh
make run
```

Run with flags:
```sh
make run ARGS="-url=http://example.com/ -method=GET -numRequest=100 -duration=5000 -warmUp=0 -change=0 -period=0 -respFile=respFile.txt"
```

In success Ltester generates HTML page that summurizes the results of the load testing: `result.html`.
There is also file, that holds all of the received responses: `respFile.txt`.
