# Realworld Fiber Clean Architecture

> ### Golang codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

My first time working on this project [here](github.com/minhhoccode111/realworldgo)

[Go Clean Template](https://github.com/evrone/go-clean-template)

## Learned Concepts

- Go + Fiber
- Clean Architecture
- Dependency Injection
- Embedded struct like `Pagination` and `Timestamps` in other struct
- `Validator` for input validation
- In `repo` layer, we don't return `pgx.ErrNoRows` to outer layer to handle
  because it will make the router layer depend on postgresql implementation,
  instead we return our custom not found error `entity.ErrNoRows`
- Go debugger with `.vscode/launch.json`

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
DELETE /articles/{slug}/comments/{commentId}

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

## Todo

- [ ] Add unit testing, mock, integration testing
- [ ] Add caching Redis

## Resources

## Contributing

Contributions are **welcome and highly appreciated**!\
This project follows the [RealWorld API
spec](https://github.com/gothinkster/realworld) — please make sure your changes
remain compliant.
