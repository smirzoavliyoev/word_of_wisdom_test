# word_of_wisdom_test

# Architecture

## Send RequestChallenge request.

When user tries to do some operation he has to request challenge. Solve That challenge and return solution.
1. Client sends request to server
2. Request will be verified
3. Server will query challenge from manual storage.

So architecture looks like:

<img width="835" alt="Screen Shot 2022-10-23 at 22 57 02" src="https://user-images.githubusercontent.com/76482536/197407957-a14b0adf-4453-477c-bc18-6c71d7568093.png">

## Next step is: 
1. after solving the challenge client have to provide solution. 
2. Solution will be verified.
3. If solution os correct the quote will be provided to client.

<img width="832" alt="Screen Shot 2022-10-23 at 23 05 31" src="https://user-images.githubusercontent.com/76482536/197408312-63b1181c-96b3-4f19-adbe-7ffd8ab77c48.png">


# How Storage works?

There is concurrency safe method where also architecture was carried about data consistency.


# Manual installation of server

# start server
``` go run server/main.go ```

# Docker
```Docker build --tag wow:1.0.0 server -f Dockerfile.server```

# run client
``` ./build.sh && ./client_pack ```
