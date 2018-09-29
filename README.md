# My Tab API

This is the backend for the [My Tab](https://my-tab.herokuapp.com/) app.
Front end code can be seen [here](https://github.com/JonathanWThom/my-tab).

### Drink Endpoints

The base url is `https://my-tab-api.herokuapp.com/`.

`GET /drinks`

* Returns all drinks for the current user.
* Includes metadata about total number of standard drinks and average per day
across the time period.
* Can receive optional query parameters of start and end date

`POST /drinks`

* Creates a drink record.
* Must include `oz`, `percent`, and `imbibed_on` JSON in request body.

`DELETE /drinks/{id}`

* Deletes a drink record.
* Only returns server errors, if any.

### User Endpoints

`POST /signup`
* Creates a user record.
* Must include `username` and `password` JSON in request body.
* Returns JWT token for user.

`POST /login`
* Finds user record by `username` and `password`.
* Returns JWT token for user.

### JWT Authentication & Authorization

* To access resources, include the following as the `Authorization` header:
`Bearer your-json-web-token`.
* Drink resources are restricted to those that belong to the current user, as
determined by the JWT.

### Standard Drinks

* One standard drink is 0.6 oz of pure alcohol - 12 oz at 5% ABV (beer), 5 oz at 12% (wine),
or 1.5 oz at 40% (liquor) would all be one standard drink.
* The `stddrink` library is a helper for determining standard drinks across time.
The functions are model-agnostic, so can be used outside of this API.

* Public methods include:
  - `Calculate`, determine number of standard drinks for a given oz/percent
  - 'StddrinksPerDay', given a time range and group of standard drinks (`float64`s), determine
  drinks per day
  - `TotalStdDrinks`, given a group of standard drinks (`float64`s), determine the
  grand total.

### Setup

If you wanted to run this locally, you'd need the following:
- [Go](https://golang.org/)
- RSA [public & private keys](https://gist.github.com/ygotthilf/baa58da5c3dd1f69fae9).
- [Postgresql](https://www.postgresql.org/). Create the `my_tab` database and add
the tables included in `tables.sql`.

### Todo and Known Bugs

* More test coverage
* Dates are returned the client as the day before what was entered.

### License

MIT
