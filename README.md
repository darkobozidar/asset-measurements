# README

## Setup
* Copy the ".env.example" to ".env" (run `cp .env.example .env`)
* Run `docker compose up --build`
    * Simulation starts automatically on `simulator` service start.
    * Postgres and MongoDB databases are created automatically, Postgres tables are populated with fixtures.

## Assignment

Parts of the assignment that weren't done: "Optional Requirements".

## Endpoints

Asset CRUD:
* Get Asset by ID
    * Pattern: `GET "/asset/:id"`
    * Example: `GET "http://{host}:8080/asset/1"`
* Get All Assets
    * Pattern: `GET "/asset?isEnabled={isEnabled?}&type={type?}"`
    * Example: `GET "http://{host}:8080/asset?isEnabled=true&type=Type"`
* Create Asset
    * Pattern: `POST "/asset/"`
    * Example:
        * URL: `POST http://{host}:8080/asset`
        * Body: `{ "name": "Name", "description": "Description", "type": "Type", "isEnabled": true }`
* Update Asset
    * Pattern: `PUT "/asset/:id"`
    * Example:
        * URL: `PUT "http://{host}:8080/asset/1"`
        * Body: `{ "id": 1, "name": "New name" }`
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
    * Pattern: `GET "/measurement/:id?from={from}&to={to}&sort={asc|desc?}"`
    * Example: `GET "http://{host}:8080/measurement/1?from=2025-05-14T10:27:37Z&to=2025-05-14T15:52:18Z&sort=desc"`
* Get Grouped Measurements Average
    * Pattern  `GET "/measurement/:id/average?from={from}&to={to}&groupBy={1minute|15minute|1hour?}&sort={asc|desc?}"`
    * Example: `GET "http://{host}:8080/measurement/1/average?from=2025-05-14T05:40:00Z&to=2027-05-14T14:15:59Z&groupBy=15minute&sort=asc"`
