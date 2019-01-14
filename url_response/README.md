
# URL Response


#### Parse a URL list checking HTTP codes and response times, facilitated by Go's concurrency.


## Build

```bash
    git clone https://github.com/Tinram/URL-Response.git
    cd URL-Response/url_response
```

```bash
    go get github.com/asaskevich/govalidator

    go build url_response.go
```

For a slightly smaller executable:

```bash
    go build -ldflags="-s -w" url_response.go
```


## Run

The program by default expects a file of URLs called *urls.txt* to be present in the same directory:

```bash
    ./url_response
```

or for an alternative filename:

```bash
    ./url_response -f urls_test.txt
```

The URLs file requires one URL per line (<kbd>LF</kbd> or <kbd>CR</kbd><kbd>LF</kbd> line endings).

Invalid URLs will be stated in the terminal output.


## Results

Output is to the terminal and added/appended to a logfile called *url_response.log*.

Invalid URLs will display in the terminal but not in the logfile.

URL | HTTP code | HTTP msg | response time |
---- | ---- | ---- | ---- |
https://www.bbc.co.uk/ | 200 | OK | 0.11741 s |


## Credits

+ Mike Schilli: inspiration for revised and more effective channel pattern.
+ Alex Saskevich: govalidator URL check.


## License

URL Response is released under the [GPL v.3](https://www.gnu.org/licenses/gpl-3.0.html).
