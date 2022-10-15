# CFP API

## API Documentation

Get the API running:

```sh
go run main.go
```

### Speakers

Create a Speaker:

```bash
curl -sd '{"ID":"default/ScottRigby","name":"Scott Rigby","bio":"Scott is a rad dad","email":"scott@email.com"}' \
-H "Content-Type: application/json" \
-X POST localhost:8080/api/speakers | jq
{
  "ID": "default/ScottRigby",
  "Name": "Scott Rigby",
  "Bio": "Scott is a rad dad",
  "Email": "scott@email.com",
  "Timestamp": "0001-01-01T00:00:00Z"
}
```

Get all Speakers:

```bash
curl -sX GET localhost:8080/api/speakers | jq
[
  {
    "ID": "default/ScottRigby",
    "Name": "NewName",
    "Bio": "Scott is a rad dev",
    "Email": "scott@email.com",
    "Timestamp": "0001-01-01T00:00:00Z"
  }
]
```

Get a Speaker by ID:
```bash
curl -sX GET localhost:8080/api/speakers/default-ScottRigby | jq
[
  {
    "ID": "default/ScottRigby",
    "Name": "NewName",
    "Bio": "Scott is a rad dev",
    "Email": "scott@email.com",
    "Timestamp": "0001-01-01T00:00:00Z"
  }
]
```

Update a Speaker:

```bash
curl -sd '{"ID":"default/ScottRigby","name":"NewName","bio":"Scott is a rad dev","email":"scott@email.com"}' \
-H "Content-Type: application/json" \
-X PUT localhost:8080/api/speakers/default-ScottRigby | jq
{
  "ID": "default/ScottRigby",
  "Name": "NewName",
  "Bio": "Scott is a rad dev",
  "Email": "scott@email.com",
  "Timestamp": "0001-01-01T00:00:00Z"
}
```

Delete a Speaker:

```bash
curl -X DELETE localhost:8080/api/speakers/default-ScottRigby
```


### Proposals

Create a Proposal:
```bash
curl -sd '{"ID":"default/MyAwesomeTalk","Title":"my awesome talk","Abstract":"This is a rad talk","Type":"lightning talk","SpeakerID":"default/ScottRigby","Final":false,"Submission":{"Status":"draft"}}' \
-X POST localhost:8080/api/proposals | jq
{
  "ID": "default/MyAwesomeTalk",
  "Title": "my awesome talk",
  "Abstract": "This is a rad talk",
  "Type": "lightning talk",
  "SpeakerID": "default/ScottRigby",
  "Final": false,
  "Submission": {
    "LastUpdate": "0001-01-01T00:00:00Z",
    "Status": "draft"
  }
}
```

Get all Proposals:

```bash
curl -sX GET localhost:8080/api/proposals | jq
[
  {
    "ID": "default/MyAwesomeTalk",
    "Title": "my awesome talk",
    "Abstract": "This is a rad talk",
    "Type": "lightning talk",
    "SpeakerID": "default/ScottRigby",
    "Final": false,
    "Submission": {
      "LastUpdate": "0001-01-01T00:00:00Z",
      "Status": "draft"
    }
  },
  {
    "ID": "default/AnotherCoolTalk",
    "Title": "another cool talk",
    "Abstract": "This is a super rad talk",
    "Type": "lightning talk",
    "SpeakerID": "default/ScottRigby",
    "Final": false,
    "Submission": {
      "LastUpdate": "0001-01-01T00:00:00Z",
      "Status": "draft"
    }
  }
]
```

Get a Proposal by ID:

```bash
curl -sX GET localhost:8080/api/proposals/default-MyAwesomeTalk | jq
{
  "ID": "default/MyAwesomeTalk",
  "Title": "my awesome talk",
  "Abstract": "This is a rad talk",
  "Type": "lightning talk",
  "SpeakerID": "default/ScottRigby",
  "Final": false,
  "Submission": {
    "LastUpdate": "0001-01-01T00:00:00Z",
    "Status": "draft"
  }
}
```

Update a Proposal:

```bash
curl -sd '{"ID":"default/MyAwesomeTalk","Title":"NewTalkTitle","Abstract":"This is a rad talk","Type":"lightning talk","SpeakerID":"default/ScottRigby","Final":false,"Submission":{"Status":"draft"}}' \
-X PUT localhost:8080/api/proposals/default-MyAwesomeTalk | jq
{
  "ID": "default/MyAwesomeTalk",
  "Title": "my very awesome talk",
  "Abstract": "This is a rad talk",
  "Type": "lightning talk",
  "SpeakerID": "default/ScottRigby",
  "Final": false,
  "Submission": {
    "LastUpdate": "2022-10-13T15:23:17.978854+02:00",
    "Status": "draft"
  }
}
```

Delete a Proposal:

```bash
curl -X DELETE localhost:8080/api/proposals/default-MyAwesomeTalk
```
