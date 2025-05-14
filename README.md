# README

## Setup
* Copy the ".env.example" to ".env" (run `cp .env.example .env`)
* Run `docker compose up --build`
* Postgres and MongoDB should be created automatically, data is filled with fixtures.

## Assignment

Parts of the assignment that weren't done: only "Optional Requirements"

## Endpoints

Asset Service:
* Get Asset by ID
    * Pattern: `"/asset/:id"`
    * Example: `"http://{host}:8080/asset/2"`
* Get All Assets
    * Pattern: `"/asset?isEnabled={val1?}&type={val2}"`
    * Example: `"http://{host}:8080/asset?isEnabled=true&type=Type"`
* Create Asset
    * Pattern: `"/asset/"`
    * Example: `http://{host}:8080/asset`

TODO
