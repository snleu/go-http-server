# go-http-server

## Challenge
`Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to the return the correct numbers after restarting it, by persisting data to a file.`

## Questions/Concerns
1. What kind of requests should be possible? Should anything be possible or only one?
2. On start up, load the requests from a json file. 
  a. However, when should the updated db be written to the file? When a request is received? 
  b. Is this overly complicated; is there a better way to write new data to file?
3. At the moment, several standard libraries are being used. How to get it down to one?
