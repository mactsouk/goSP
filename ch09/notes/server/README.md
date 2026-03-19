# Testing the API with `curl(1)`

## Create a note

```
    curl -X POST http://localhost:8080/notes \
     -H "Content-Type: application/json" \
     -d '{"title":"Shopping","content":"Buy milk"}'
```

## List all notes

`curl http://localhost:8080/notes`

## Filter notes by title

`curl "http://localhost:8080/notes?title=shop"`

## Get a note by ID

`curl http://localhost:8080/notes/1`

## Update a note by ID

```
    curl -X PUT http://localhost:8080/notes/1 \
     -H "Content-Type: application/json" \
     -d '{"title":"Shopping List","content":"Buy milk and bread"}'
```

## Delete a note by ID

curl -X DELETE http://localhost:8080/notes/1
