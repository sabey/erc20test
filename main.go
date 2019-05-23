package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*

Using only the Golang standard library and the etherscan.io/apis create a program that:

1. Accepts an ethereum address as a flag `./betest -address="0xea38eaa3c86c8f9b751533ba2e562deb9acded40"`
2. Retrieves "Get a list of "ERC20 - Token Transfer Events" by Address" in ascending order for the FUEL contract: 0xea38eaa3c86c8f9b751533ba2e562deb9acded40
3. Print to the stdour/stderr TWO JSON Transaction objects: one containing the largest `value` and another for the largest `gasPrice`
4. Print to a file `results.json` a JSON array of all of the Transactions

*/

const (
	env      = `api` // `api-ropsten`
	tokenURL = `http://%s.etherscan.io/api?module=account&action=tokentx&contractaddress=%s&startblock=0&endblock=999999999&sort=asc&page=%d&offset=10000`
)
const (
	BinanceBNB = `0xB8c77482e45F1F44dE1745F52C74426C631bDD52`
	Fuel       = `0xea38eaa3c86c8f9b751533ba2e562deb9acded40`
)

var (
	address = flag.String("address", "", "Contract Address")
)

func init() {
	flag.Parse()
}
func main() {
	log.Println("what it do baby boo")
	if *address == "" {
		log.Fatalln("address was empty!?")
		return
	}
	// get all transactions
	results, err := getAll(*address)
	if err != nil {
		log.Fatalf("failed to getAll \"%s\"\n", err)
		return
	}
	// found results?!
	if len(results) == 0 {
		log.Fatalln("found not result!?")
		return
	}
	log.Printf("results: %d\n", len(results))
	// find interesting transactions
	var largestValue *Transaction
	var largestGasPrice *Transaction
	for _, tx := range results {
		if largestValue == nil {
			largestValue = tx
		} else {
			if largestValue.Value.Cmp(tx.Value) < 0 {
				largestValue = tx
			}
		}
		if largestGasPrice == nil {
			largestGasPrice = tx
		} else {
			if largestGasPrice.GasPrice < tx.GasPrice {
				largestGasPrice = tx
			}
		}
	}
	// print interesting transactions
	bs, _ := json.Marshal(largestValue)
	log.Printf("largestValue: \"%s\" tx: \"%s\"\n", largestValue.Value, bs)
	bs, _ = json.Marshal(largestGasPrice)
	log.Printf("largestGasPrice: \"%d\" tx: \"%s\"\n", largestGasPrice.GasPrice, bs)
	// save results
	f, err := os.OpenFile("results.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("failed to open `results.json`: \"%s\"\n", err)
		return
	}
	// encode json to file
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(results); err != nil {
		log.Fatalf("failed to encode `results.json`: \"%s\"\n", err)
		return
	}
	// pz
	log.Println("all done!")
}
func getAll(
	address string,
) (
	[]*Transaction,
	error,
) {
	results := []*Transaction{}
	failures := 0
	for i := 1; i <= 10; i++ {
		// get transactions
		response, err := get(address, i)
		if err != nil {
			log.Fatalf("failed to get: \"%s\"\n", err)
			return nil, err
		}
		log.Printf("got(%d)\n", i)
		// read response
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("failed to read body: \"%s\"\n", err)
			return nil, err
		}
		// parse result
		result := &Result{}
		if err := json.Unmarshal(body, result); err != nil {
			log.Printf("original body: \"%s\"\n", body)
			log.Printf("failed to unmarshal result: \"%s\"\n", err)
			/*
				<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
				<html xmlns="http://www.w3.org/1999/xhtml">
				<head>
				<meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1"/>
				<title>403 - Forbidden: Access is denied.</title>
				<style type="text/css">
				<!--
				body{margin:0;font-size:.7em;font-family:Verdana, Arial, Helvetica, sans-serif;background:#EEEEEE;}
				fieldset{padding:0 15px 10px 15px;}
				h1{font-size:2.4em;margin:0;color:#FFF;}
				h2{font-size:1.7em;margin:0;color:#CC0000;}
				h3{font-size:1.2em;margin:10px 0 0 0;color:#000000;}
				#header{width:96%;margin:0 0 0 0;padding:6px 2% 6px 2%;font-family:"trebuchet MS", Verdana, sans-serif;color:#FFF;
				background-color:#555555;}
				#content{margin:0 0 0 2%;position:relative;}
				.content-container{background:#FFF;width:96%;margin-top:8px;padding:10px;position:relative;}
				-->
				</style>
				</head>
				<body>
				<div id="header"><h1>Server Error</h1></div>
				<div id="content">
				 <div class="content-container"><fieldset>
				  <h2>403 - Forbidden: Access is denied.</h2>
				  <h3>You do not have permission to view this directory or page using the credentials that you supplied.</h3>
				 </fieldset></div>
				</div>
				</body>
				</html>
			*/
			if failures < 5 {
				// try again
				log.Printf("ignoring failure %d - trying again!", failures+1)
				failures++
				i--
				continue
			}
			log.Printf("failed too many times!!!")
			return nil, err
		}
		log.Printf("got(%d) - Results: %d Status: \"%s\" Message: \"%s\"\n", i, len(result.Result), result.Status, result.Message)
		if result.Status == "0" {
			// no transactions left
			log.Println("no transactions remaining - we're done")
			break
		}
		results = append(results, result.Result...)
		failures = 0
	}
	return results, nil
}
func get(
	address string,
	page int,
) (
	*http.Response,
	error,
) {
	if page < 1 {
		page = 1
	}
	url := fmt.Sprintf(tokenURL,
		env,
		address,
		page,
	)
	log.Printf("GET \"%s\"\n", url)
	return http.Get(url)
}
