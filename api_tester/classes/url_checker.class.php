<?php

declare(strict_types=1);

namespace Tinram\URLChecker;

final class URLChecker
{
    /**
        * URLChecker
        *
        * URL Tester and benchmarker using cURL multi for concurrency.
        *
        * Coded to PHP 7.1
        *
        * @author         Martin Latter
        * @copyright      Martin Latter, 07/01/2019
        * @version        0.05
        * @license        GNU GPL version 3.0 (GPL v3); http://www.gnu.org/licenses/gpl.html
        * @link           https://github.com/Tinram/URL-Response.git
    */


    /* @var string, user agent */
    private $sUserAgent = 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:64.0) Gecko/20100101 Firefox/64.0';

    /* @var array, cURL options holder */
    private $aOpts = [];

    /* @var array, URLs holder */
    private $aURLs = [];

    /* @var string, default logfile name */
    private $sLogfile = 'url_checker.log';

    /* @var int, default cURL request batch size */
    private $iBatchSize = 100;


    public function __construct(array $aURLs = null)
    {
        if (is_null($aURLs))
        {
            die(__METHOD__ . '() requires an array of URLs to be passed.' . PHP_EOL);
        }

        $this->sLogfile = (defined('LOG_FILE')) ? LOG_FILE : $this->sLogfile;
        $this->iBatchSize = (defined('BATCH_SIZE')) ? BATCH_SIZE : $this->iBatchSize;

        $this->aURLs = $aURLs;
        $this->setup();
        $aBatches = $this->createBatches();
        $this->logWrite('start', true);
        $sMessage = $this->runner($aBatches);
        $this->outputMessage($sMessage);
    }


    /**
        * Set-up cURL options
        *
        * @return  void
    */

    private function setup(): void    /* remove :void for PHP 7.0 */
    {
        $this->aOpts =
        [
            CURLOPT_HEADER => false,
            CURLOPT_TIMEOUT => 30,
            CURLOPT_NOBODY => true,
            CURLOPT_USERAGENT => $this->sUserAgent,
            CURLOPT_FAILONERROR => true,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_CONNECTTIMEOUT => 5,
            CURLOPT_IPRESOLVE => CURL_IPRESOLVE_V4
        ];
    }


    /**
        * Split array of URLs into chunks/batches for cURL processing.
        *
        * @return  array
    */

    private function createBatches(): array
    {
        return array_chunk($this->aURLs, $this->iBatchSize);
    }


    /**
        * Process URL arrays by cURL multi, deliberately enforcing a slight delay per batch.
        *
        * @return  string, time details
    */

    private function runner($aBatches): string
    {
        $sOutput = '';
        $fTS = microtime(true);

        # enforce delay on cURL multi to reduce server overload/dropped connections
        foreach ($aBatches as $aBatch)
        {
            $this->process($aBatch);
            usleep(1000);
        }

        $fTE = microtime(true);
        $sOutput .= 'URL checker run, see generated logfile ' . $this->sLogfile . PHP_EOL;
        $sOutput .= sprintf('Total time taken: %01.3f', $fTE - $fTS) . ' s' . PHP_EOL;

        return $sOutput;
    }


    /**
        * Process URL array with cURL multi.
        *
        * @return  void
    */

    private function process(array $aURLs): void
    {
        $aCurlHandles = [];
        $iRunning = 0;

        $rMh = curl_multi_init();

        foreach ($aURLs as $sUrl)
        {
            $rCh = curl_init($sUrl);
            curl_setopt_array($rCh, $this->aOpts);
            curl_multi_add_handle($rMh, $rCh);
            $aCurlHandles[$sUrl] = $rCh;
        }

        # execute cURL handles
        do
        {
            $iSH = curl_multi_exec($rMh, $iRunning);
        }
        while ($iSH === CURLM_CALL_MULTI_PERFORM);

        # for Windows cURL multi hanging, credit: xxavalanchexx@gmail.com
        if (curl_multi_select($rMh) === -1)
        {
            usleep(100);
        }

        while ($iRunning && $iSH === CURLM_OK)
        {
            do
            {
                $iSH2 = curl_multi_exec($rMh, $iRunning);
            }
            while ($iSH2 === CURLM_CALL_MULTI_PERFORM);
        }

        # grab URL content, remove handles
        foreach ($aCurlHandles as $rCh)
        {
            $aResults = curl_getinfo($rCh);
            # URL | HTTP code | time taken
            $this->logWrite($aResults['url'] . ' | ' . $aResults['http_code'] . ' | ' . $aResults['total_time']);
            curl_multi_remove_handle($rMh, $rCh);
        }

        curl_multi_close($rMh);
    }


    /**
        * Log messages to file.
        *
        * @param   string $sMessage, message to log
        * @param   boolean $bTimestamp, toggle to include/omit timestamp on line
        * @return  void
    */

    private function logWrite(string $sMessage = '', bool $bTimestamp = false): void
    {
        if (empty($sMessage))
        {
            return;
        }

        if ($bTimestamp)
        {
            $sMessage = $this->getTimestamp() . ' | ' . $sMessage;
        }

        $iLogWrite = file_put_contents($this->sLogfile, $sMessage . PHP_EOL, FILE_APPEND);

        if (!$iLogWrite)
        {
            die('could not write to logfile ' . $this->sLogfile . PHP_EOL);
        }
    }


    /**
        * Return a timestamp with a custom format.
        *
        * @return  string, custom date format
    */

    private function getTimestamp(): string
    {
        return date('Y-m-d H:i:s P T');
    }


    /**
        * Output message.
        *
        * @param   string $sMessage, message to print
        * @return  void
    */

    private function outputMessage(string $sMessage = ''): void
    {
        if (empty($sMessage))
        {
            return;
        }
        else
        {
            echo $sMessage;
        }
    }
}
