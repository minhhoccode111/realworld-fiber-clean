# Realworld Fiber Clean Architecture

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

## How to start

## Resources
