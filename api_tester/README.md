
## API Tester

Using */classes/url_checker.class.php*


### Purpose

For an API: test selected GET endpoints and response times.

Uses cURL multi for concurrent URL requests, and batches of cURL requests to minimise server dropped connections.


### Set-up

In *api_tester.php* &ndash; configure the API URLs, API URL increment sequence, and API query strings to be scanned.

*api_tester.php* currently defines a server running locally on port 8000 and endpoint URL sequence numbers.


### Execute

        php api_tester.php


## License

API Tester is released under the [GPL v.3](https://www.gnu.org/licenses/gpl-3.0.html).
