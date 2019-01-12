#!/usr/bin/env php
<?php

/**
    * Command-line example usage of using url_checker class to test GET API endpoint.
    *
    * Usage:
    *          php api_tester.php
    *
    * PHP provides a slow-responding server: php -S localhost:8000
    * else something more robust (nginx, Go) can be used as a testing endpoint.
    *
    * @author         Martin Latter
    * @copyright      Martin Latter 07/01/2019
    * @version        0.02
    * @license        GNU GPL version 3.0 (GPL v3); http://www.gnu.org/licenses/gpl.html
    * @link           https://github.com/Tinram/URL-Response.git
*/


require('classes/url_checker.class.php');

use Tinram\URLChecker\URLChecker;

date_default_timezone_set('Europe/London');


define('URL', 'http://localhost:8000/');
define('LOG_FILE', 'api_tester.log');
define('BATCH_SIZE', 200); # size of each cURL request batch


$aURLs = [];
$qs = '?style=red';


# generate API endpoint iterations + query string
for ($i = 1000; $i < 21000; $i++)
{
    $aURLs[] = URL . $i . $qs;
}

new URLChecker($aURLs);
