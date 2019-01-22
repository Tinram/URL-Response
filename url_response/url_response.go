/**
	* URL Response
	*
	* Parse a file of URLs and check HTTP codes and response times.
	*
	* Usage:
	*                ./url_response [-f <filename>]
	*
	* @author        Martin Latter
	* @copyright     Martin Latter 10/12/2018
	* @version       0.05
	* @license       GNU GPL version 3.0 (GPL v3); http://www.gnu.org/licenses/gpl.html
	* @link          https://github.com/Tinram/URL-Response.git
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
)

type urlResults struct {
	Error        error
	URL          string
	ResponseCode int
	ResponseMsg  string
	Time         float64
}

func main() {

	const LOGNAME string = "url_response.log"
	var urls []string

	filename := "urls.txt"
	channelLimit := 100 // 100 is good for ~1000 URLs

	flag.Usage = func() {
		usageText := "  url_response [-f]\n\tdefault urls file is urls.txt with one URL per line\n\tuse -f to specify alternative filename\n"
		fmt.Fprintf(os.Stderr, "%s", usageText)
	}

	f := flag.String("f", "", "specify alternative filename")
	flag.Parse()

	if *f != "" {
		filename = *f
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		tmpURL := scanner.Text()

		/* validate URL */
		validURL := govalidator.IsRequestURL(tmpURL)

		if validURL != false {
			urls = append(urls, tmpURL)
		} else {
			fmt.Printf("%s is an invalid URL\n", tmpURL)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	results := fetch(urls, channelLimit)

	/* prepare logging to file in loop */
	flog, errLog := os.OpenFile(LOGNAME, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0640)
	if errLog != nil {
		log.Fatal(errLog)
	}
	defer flog.Close()

	for i := 0; i < len(urls); i++ {

		result := <-results

		if result.Error == nil {

			output := fmt.Sprintf("%s | %d | %s | %.5f s", result.URL, result.ResponseCode, result.ResponseMsg, result.Time)
			fmt.Println(output)

			/* log */
			log.SetOutput(flog)
			log.Printf("| " + output)
		}
	}
}

func fetch(urls []string, channelLimit int) <-chan urlResults {

	results := make(chan urlResults, channelLimit)

	for _, url := range urls {

		go func(url string) {

			/* avoid default http client */
			ht := &http.Transport{
				IdleConnTimeout: 5 * time.Second,
			}
			client := &http.Client{
				Transport: ht,
				Timeout:   10 * time.Second,
			}

			start := time.Now()

			resp, err := client.Get(url)

			if err != nil {
				fmt.Println(err)
			} else {
				defer resp.Body.Close()
				resp.Close = true
				elapsed := time.Since(start).Seconds()
				results <- urlResults{Error: err, URL: url, ResponseCode: resp.StatusCode, ResponseMsg: http.StatusText(resp.StatusCode), Time: elapsed}
			}

		}(url)
	}

	return results
}
