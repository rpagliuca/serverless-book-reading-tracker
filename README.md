# About

* Book Reading Tracker 
* Serverless (AWS Lambda + Golang)

# API

## Entries

```json
(POST) /entries
{
    "book_id": "85167e72-65a3-4bce-ba46-80126ecbc3ab",
    "start_time": "2021-03-13T11:15:00Z"
}
// Optional properties are book_id, start_time , end_time, start_location and end_location
// No required properties

200:
{
    "success": true,
    "message": "OK"
}

422:
{
    "success": false,
    "message": "Invalid input"
}
```

```json
(GET) /entries

200:
{
    "ok": true,
    "entries": [
        {
            "id": "38d6fc1d-abd0-4487-8536-14297ab6abd8",
            "book_id": "776719c4-fdba-4510-9494-b457c684f496",
            "start_time": "2021-03-13T11:15:00Z",
            "end_time": "2021-03-13T12:15:00Z",
            "start_location": "150",
            "end_location": "200",
            "date_created": "2021-03-13T11:12:00Z",
            "date_modified": "2021-03-13T11:16:00Z",
            "version": 3
        },
        {
            "id": "0b1aa72b-134e-4fce-9dec-e3e1b827b363",
            "book_id": "776719c4-fdba-4510-9494-b457c684f496",
            "start_time": "2021-03-13T11:15:00Z",
            "end_time": "2021-03-13T12:15:00Z",
            "start_location": "150",
            "end_location": "200",
            "date_created": "2021-03-13T11:12:00Z",
            "date_modified": "2021-03-13T11:16:00Z",
            "version": 1
        }
    ]
}
```

```json
(GET) /entry/38d6fc1d-abd0-4487-8536-14297ab6abd8

200:
{
    "success": true,
    "message": "OK",
    "entry": {
        "id": "38d6fc1d-abd0-4487-8536-14297ab6abd8",
        "book_id": "776719c4-fdba-4510-9494-b457c684f496",
        "start_time": "2021-03-13T11:15:00Z",
        "end_time": "2021-03-13T12:15:00Z",
        "start_location": "150",
        "end_location": "200",
        "date_created": "2021-03-13T11:12:00Z",
        "date_modified": "2021-03-13T11:16:00Z",
        "version": 1
    }
}

404:
{
    "success": false,
    "message": "Not found"
}
```

```json
(PATCH) /entry/38d6fc1d-abd0-4487-8536-14297ab6abd8
{
    "end_location": "205",
    "version": 2
}
// Optional properties are book_id, start_time , end_time, start_location and end_location
// Required properties are version

200:
{
    "success": true,
    "message": "OK"
}

404:
{
    "success": false,
    "message": "Not found"
}

422:
{
    "success": false,
    "message": "Invalid input"
}
```

## Generic responses
```json
401:
{
    "success": false,
    "message": "Invalid authentication"
}

500:
{
    "success": false,
    "message": "Internal server error"
}
```
