# REST API TV Series Service

## Endpoints

| Endpoint                    | Method | Description |
| --------------------------- | :----: | :-------------------------: |
| `/v1/series`                | `POST` | `Create series`             |
| `/v1/series?q={{query}}`    | `GET`  | `Find series by title`      |
| `/v1/series/{{id}}`         | `GET`  | `Get series by id`          |
| `/v1/series/{{id}}/reviews` | `GET`  | `Get series' reviews by id` |
| `/v1/reviews`               | `POST` | `Create a review`           |

## TODO

- Docker
- Logging 
- Testing database packages
- Use mocking frameworks
