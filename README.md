# REST API TV Series Service

## Endpoints

| Endpoint                    | Method | Description |
| --------------------------- | :----: | :-------------------------: |
| `/v1/series`                | `POST` | `Create series`             |
| `/v1/series?q={{query}}`    | `GET`  | `Find series by title`      |
| `/v1/series/{{id}}`         | `GET`  | `Get series by id`          |
| `/v1/series/{{id}}/reviews` | `GET`  | `Get series' reviews by id` |
| `/v1/reviews`               | `POST` | `Create a review`           |

## Testing endpoints using cURL

- Create a series

**Request**

```
curl --request POST 'localhost:8000/v1/series' \
  --header 'Content-Type: "application/json"' \
  --data-raw '{
      "title": "Title",
      "description": "Description",
      "episodes": 20,
      "begin_year": 2022,
      "creator": "Creator"
  }'
```

**Response**
```
{
    "id": "635c4bb1-41b1-46e9-a87e-292a9d63b9e3",
    "title": "Title",
    "description": "Description",
    "episodes": 20,
    "begin_year": 2022,
    "end_year": 0,
    "creator": "Creator"
}
```

- Create a review for a series

**Request**

```
curl --request POST 'localhost:8000/v1/reviews' \
  --header 'Content-Type: "application/json"' \
  --data-raw '{
      "series_id": "635c4bb1-41b1-46e9-a87e-292a9d63b9e3",
      "author_id": "635c4bb1-41b1-46e9-a87e-292a9d63b9e3",
      "text": "This is a great show"
  }'
```

**Response**

```
{
    "id": "26efa50a-953e-4aeb-befb-ccc14058989b",
    "series_id": "635c4bb1-41b1-46e9-a87e-292a9d63b9e3",
    "author_id": "635c4bb1-41b1-46e9-a87e-292a9d63b9e3",
    "text": "This is a great show"
}
```

- Get series by ID

**Request**

`curl --request GET 'localhost:8000/v1/series/{{series_id}}'`

**Response**

```
{
    "id": "635c4bb1-41b1-46e9-a87e-292a9d63b9e3",
    "title": "Title",
    "description": "Description",
    "episodes": 20,
    "begin_year": 2022,
    "end_year": 0,
    "creator": "Creator",
    "reviews": [
        {
            "id": "26efa50a-953e-4aeb-befb-ccc14058989b",
            "author_id": "635c4bb1-41b1-46e9-a87e-292a9d63b9e3",
            "text": "This is a great show"
        }
    ]
}
```

- Get series's reviews

**Request**

`curl --request GET 'localhost:8000/v1/series/{{series_id}}/reviews'`

**Response**

```
{
    "reviews": [
        {
            "id": "26efa50a-953e-4aeb-befb-ccc14058989b",
            "author_id": "635c4bb1-41b1-46e9-a87e-292a9d63b9e3",
            "text": "This is a great show"
        }
    ]
}
```

- Find series by title

**Request**

`curl --request GET 'localhost:8000/v1/series?q={{query}}'`

**Response**

```
{
  "series": [
      {
          "id": "635c4bb1-41b1-46e9-a87e-292a9d63b9e3",
          "title": "Title",
          "begin_year": 2022,
          "end_year": 0,
          "creator": "Creator"
      }
  ]
}
```

## TODO

- Swagger documentation
- Logging 
- Testing database packages
- Use mocking frameworks
