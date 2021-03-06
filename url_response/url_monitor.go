/**
	* URL Monitor
	*
	* Parse a file of URLs and continuously check HTTP codes and response times.
	*
	* Usage:
	*                ./url_monitor [-f <filename>] [-d <delay_secs>] [-t <timeout_secs>]
	*
	* @author        Martin Latter
	* @copyright     Martin Latter 01/05/2019
	* @version       0.06
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
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gookit/color"
)

type urlResults struct {
	Error        error
	URL          string
	ResponseCode int
	ResponseMsg  string
	Time         float64
}

func main() {

	const channelLimit = 100

	var urls []string
	var output string

	mainrun := false
	filename := "urls.txt"
	checktime := 30
	timeout := 6
	red := color.FgRed.Render
	green := color.FgGreen.Render

	flag.Usage = func() {
		usageText := "  url_response [-f] [-d] [-t]\n\tdefault urls file is urls.txt with one URL per line\n\t-f to use alternative filename\n\t-d to change check time duration (default: 30s)\n\t-t to change request timeout (default: 6s)\n"
		fmt.Fprintf(os.Stderr, "%s", usageText)
	}

	f := flag.String("f", "", "specify alternative filename")
	d := flag.String("d", "", "specify check time duration")
	t := flag.String("t", "", "specify request timeout")
	flag.Parse()

	if *f != "" {
		filename = *f
	}

	if *d != "" {
		checktime, _ = strconv.Atoi(*d)
	}

	if *t != "" {
		timeout, _ = strconv.Atoi(*t)
	}

	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()

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
			fmt.Printf(" %s is an invalid URL\n", red(tmpURL))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for {

		if mainrun {
			c = exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			c = exec.Command("tput", "civis")
			c.Stdout = os.Stdout
			c.Run()
		}

		mainrun = true

		results := fetch(urls, channelLimit, timeout)

		fmt.Println()

		for i := 0; i < len(urls); i++ {

			result := <-results

			if result.Error == nil {

				switch result.ResponseCode {

					case 200, 203, 206, 300, 301, 302, 303, 304, 307, 308:
						output = fmt.Sprintf(" %s  %s   %.3fs   %s", green(result.ResponseCode), result.ResponseMsg, result.Time, result.URL)

					default:
						output = fmt.Sprintf(" %s  %s   %.3fs   %s", red(result.ResponseCode), result.ResponseMsg, result.Time, result.URL)
				}
			} else {

				switch result.ResponseCode {

					case 901:
						output = fmt.Sprintf(" %s %s        %s", red("---"), red(result.ResponseMsg), result.URL)
					case 902:
						output = fmt.Sprintf(" %s %s    %s", red("---"), red(result.ResponseMsg), result.URL)
					case 903:
						output = fmt.Sprintf(" %s %s   %s", red("---"), red(result.ResponseMsg), result.URL)
				}

			}

			fmt.Println(output)

		}

		time.Sleep(time.Second * time.Duration(checktime))
	}

}

func fetch(urls []string, channelLimit int, timeout int) <-chan urlResults {

	results := make(chan urlResults, channelLimit)

	for _, url := range urls {

		go func(url string) {

			/* avoid default http client */
			ht := &http.Transport{
				IdleConnTimeout: time.Duration(timeout) * time.Second,
			}
			client := &http.Client{
				Transport: ht,
				Timeout:   time.Duration(timeout) * time.Second,
			}

			start := time.Now()

			resp, err := client.Get(url)

			if err != nil {
				s := err.Error()
				if strings.Index(s, "no such host") > -1 {
					results <- urlResults{Error: err, URL: url, ResponseCode: 901, ResponseMsg: "no host", Time: 0}
				} else if strings.Index(s, "request canceled") > -1 {
					results <- urlResults{Error: err, URL: url, ResponseCode: 902, ResponseMsg: "unreachable", Time: 0}
				} else if strings.Index(s, "connection refused") > -1 {
					results <- urlResults{Error: err, URL: url, ResponseCode: 903, ResponseMsg: "conn refused", Time: 0}
				} else {
					fmt.Println(err)
				}
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
