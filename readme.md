#Github client
Very basic client using channels and goroutines.
The repositories and branches are fetched using the github api. The program parses the json responses and converts it to the desired data structure. Finally the percentages of branch usage on those repositories are also printed out.

## Running the program
on the root folder (same as this readme.md file) execute

```bash
go run main.go
```
or
```bash
go build
./github-client
```

## Limitations
As this program is not using authenticated requests there's a limitation of 60 requests per hour, according to [https://developer.github.com/v3/#rate-limiting](https://developer.github.com/v3/#rate-limiting)