# README

## Setup
* Copy the ".env.example" to ".env" (run `cp .env.example .env`)
* Run `docker compose up --build`
* Postgres and MongoDB should be created automatically, data is filled with fixtures.

## Assignment

* Parts of the assignment that weren't done: only "Optional Requirements".
* Some minor TODOs left in the code.

## Endpoints

TODO add POST / PUT body

Asset CRUD:
* Get Asset by ID
    * Pattern: `GET "/asset/:id"`
    * Example: `GET "http://{host}:8080/asset/1"`
* Get All Assets
    * Pattern: `GET "/asset?isEnabled={isEnabled?}&type={type?}"`
    * Example: `GET "http://{host}:8080/asset?isEnabled=true&type=Type"`
* Create Asset
    * Pattern: `POST "/asset/"`
    * Example: `POST http://{host}:8080/asset`
* Update Asset
    * Pattern: `PUT "/asset/:id"`
    * Example: `PUT "http://{host}:8080/asset/1"`
    * Note: only the filed values that need to be updated can be specified.
* Delete Asset
    * Pattern: `DELETE "/asset/:id"`
    * Example: `DELETE "http://{host}:8080/asset/1"`
    * Note: doesn't actually delete the record, just sets `isActive` to `false`.

Asset Measurements:
* Get Latest Measurement
    * Pattern: `GET "/measurement/:id/latest"`
    * Example: `GET "http://{host}:8080/measurement/1/latest"`
* Get Measurement Range
    * Pattern: `GET "/measurement/:id?from={from}Z&to={to}&sort={asc|desc?}"`
    * Example: `GET "http://{host}:8080/measurement/1?from=2025-05-14T10:27:37Z&to=2025-05-14T15:52:18Z&sort=desc"`
* Get Grouped Measurements Average
    * Pattern  `GET "/measurement/:id/average?from={from}&to={to}&groupBy={1minute|15minute|1hour}&sort={asc|desc?}"`
    * Example: `GET "http://{host}:8080/measurement/1/average?from=2025-05-14T05:40:00Z&to=2027-05-14T14:15:59Z&groupBy=15minute&sort=asc"`
