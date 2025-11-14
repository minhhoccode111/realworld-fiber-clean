# Realworld Fiber Clean Architecture

> ### Golang codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

## Learned Concepts

- Go + Fiber
- Clean Architecture
- Dependency Injection
- Embedded struct like `Pagination` and `Timestamps` in other struct
- Generate mocks with `mockgen` and `go generate ./...`
- Build Tags like '-tags migrate'
- Struct tags for input validation in `pkg/validatorx/validator.go`
- Struct tags for env vars loading in `config/config.go`
- Swagger support in Golang project
- DB Migration with `migrate`
- `squirrel` for query builder
- `lvh.me` domain
- Go debugger with `.vscode/launch.json`
- In `repo` layer, we don't return `pgx.ErrNoRows` to outer layer to handle
  because it will make the router layer depend on postgresql implementation,
  instead we return our custom not found error `entity.ErrNoRows`
- Generate TypeScript API types

## How to start developing

Start `db`, `rabbitmq`, `nats` services first

```bash
make compose-up
```

Then start the app with swagger

```bash
make run-swag
```

Or run with debugger (`F5`) using `.vscode/launch.json` config

Or run with `air` using `.air.toml` for live reload

```
air
```

And for many more `make` commands please checkout with

```bash
make help
```

Generate `types.ts`

## Endpoints

```http
       # Auth
POST   /users
POST   /users/login
GET    /user
PUT    /user

       # Article, Favorite, Comments
POST   /articles
GET    /articles
                ?tag={tag1,tag2}
                &author={username}
                &favorited={username}
                &limit={limit}
                &offset={offset}
GET    /articles/feed
GET    /articles/{slug}
PUT    /articles/{slug}
DELETE /articles/{slug}
POST   /articles/{slug}/favorite
DELETE /articles/{slug}/favorite
POST   /articles/{slug}/comments
GET    /articles/{slug}/comments
DELETE /articles/{slug}/comments/{commentID}

       # Profiles
GET    /profiles/{username}
POST   /profiles/{username}/follow
DELETE /profiles/{username}/follow

       # Tags
GET    /tags
```

## Database Design

```txt
- Users
       - id
idx    - email     - unique
idx    - username  - unique
       - password
       - image
       - bio
       - created_at
       - updated_at
- Articles
       - id
       - author_id
idx    - slug      - unique
       - body
       - title
       - description
       - created_at
       - updated_at
       - deleted_at
- Comments
       - id
       - article_id
       - author_id
       - body
       - created_at
       - deleted_at
- ArticleTags
       - article_id
       - tag_id
- Tags
       - id
idx    - name      - unique
- Favorites
       - user_id
       - article_id
- Follows
       - follower_id
       - following_id
```

- User vs. Article: One-to-Many
- Article vs. Comment: One-to-Many
- Article vs. Tag: Many-to-Many
- Follow User vs. User: Many-to-Many
- Favorite User vs. Article: Many-to-Many

## Test

Run the test script

```bash
./docs/api-test/run-api-tests.sh
```

Expected output

```bash
┌─────────────────────────┬──────────────────┬─────────────────┐
│                         │         executed │          failed │
├─────────────────────────┼──────────────────┼─────────────────┤
│              iterations │                1 │               0 │
├─────────────────────────┼──────────────────┼─────────────────┤
│                requests │               32 │               0 │
├─────────────────────────┼──────────────────┼─────────────────┤
│            test-scripts │               48 │               0 │
├─────────────────────────┼──────────────────┼─────────────────┤
│      prerequest-scripts │               18 │               0 │
├─────────────────────────┼──────────────────┼─────────────────┤
│              assertions │              335 │               0 │
├─────────────────────────┴──────────────────┴─────────────────┤
│ total run duration: 16.8s                                    │
├──────────────────────────────────────────────────────────────┤
│ total data received: 23kB (approx)                           │
├──────────────────────────────────────────────────────────────┤
│ average response time: 9ms [min: 1ms, max: 62ms, s.d.: 16ms] │
└──────────────────────────────────────────────────────────────┘
```

## Todo

- [ ] Add mocks, unit testing, integration testing
- [ ] Add caching Redis

## Resources

- My first Realworld implementation using Golang:
  [gorealworld](github.com/minhhoccode111/realworldgo)
- Starter Template using [Golang Clean
  Architecture](https://github.com/evrone/go-clean-template)

## Contributing

Contributions are **welcome and highly appreciated**!\
This project follows the [RealWorld API
spec](https://github.com/gothinkster/realworld) — please make sure your changes
remain compliant.
