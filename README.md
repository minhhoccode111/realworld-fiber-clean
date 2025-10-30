# Realworld Fiber Clean Architecture

My first time working on this project [here](github.com/minhhoccode111/realworldgo)

[Go Clean Template](https://github.com/evrone/go-clean-template)

## Learned Concepts

- Fiber and Clean Architecture
- Embedded struct like `Pagination` in other struct
- Struct Tags + Validator
- In `repo` layer, we don't return `pgx.ErrNoRows` to outer layer to handle
  because it will make the router layer depend on postgresql implementation,
  instead we return our custom not found error `entity.ErrNoRows`

## How to start

## Resources
