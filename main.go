package main

import (
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"os"
 	"golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/gmail/v1"
)

func check(e error) {
    if e != nil {
        log.Fatal(e)
    }
}

func reportFormatter(symbol string, message string) string {
    return fmt.Sprintf("%s|%s", symbol, message)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
        // The file token.json stores the user's access and refresh tokens, and is
        // created automatically when the authorization flow completes for the first
        // time.
        tokFile := "gmail_api/token.json"
        tok, err := tokenFromFile(tokFile)
        if err != nil {
                tok = getTokenFromWeb(config)
                saveToken(tokFile, tok)
        }
        return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code: %v", err)
        }

        tok, err := config.Exchange(context.TODO(), authCode)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web: %v", err)
        }
        return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
        f, err := os.Open(file)
        if err != nil {
                return nil, err
        }
        defer f.Close()
        tok := &oauth2.Token{}
        err = json.NewDecoder(f).Decode(tok)
        return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
        fmt.Printf("Saving credential file to: %s\n", path)
        f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
        if err != nil {
                log.Fatalf("Unable to cache oauth token: %v", err)
        }
        defer f.Close()
        json.NewEncoder(f).Encode(token)
}

func sendEmail(){
	b, err := ioutil.ReadFile("gmail_api/credentials.json")
    check(err)

    // If modifying these scopes, delete your previously saved token.json.
    config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
    if err != nil {
            log.Fatalf("Unable to parse client secret file to config: %v", err)
    }
    client := getClient(config)

    srv, err := gmail.New(client)
    if err != nil {
            log.Fatalf("Unable to retrieve Gmail client: %v", err)
    }

    user := "me"
        r, err := srv.Users.Labels.List(user).Do()
        if err != nil {
                log.Fatalf("Unable to retrieve labels: %v", err)
        }
        if len(r.Labels) == 0 {
                fmt.Println("No labels found.")
                return
        }
        fmt.Println("Labels:")
        for _, l := range r.Labels {
                fmt.Printf("- %s\n", l.Name)
        }
}

type Kline struct {
	OpenTime int
	Open float64
	High float64
	Low float64
	Close float64
	Volume float64
	CloseTime int
	QuoteAssetVolume float64
	NumberOfTrades int
	TakerBuyBaseAssetVolume float64
	TakerBuyQuoteAssetVolume float64
	Data1 string
}

func (r *Kline) UnmarshalJSON(p []byte) error {
    var tmp []interface{}
    if err := json.Unmarshal(p, &tmp); err != nil {
        return err
    }
  
   	r.OpenTime = int(tmp[0].(float64))
	r.Open, _ = strconv.ParseFloat(tmp[1].(string), 64)
	r.High, _ = strconv.ParseFloat(tmp[2].(string), 64)
	r.Low, _ = strconv.ParseFloat(tmp[3].(string), 64)
	r.Close, _ = strconv.ParseFloat(tmp[4].(string), 64)
	r.Volume, _ = strconv.ParseFloat(tmp[5].(string), 64)
	r.CloseTime = int(tmp[6].(float64))
	r.QuoteAssetVolume, _ = strconv.ParseFloat(tmp[7].(string), 64)
	r.NumberOfTrades = int(tmp[8].(float64))
	r.TakerBuyBaseAssetVolume, _ = strconv.ParseFloat(tmp[9].(string), 64)
	r.TakerBuyQuoteAssetVolume, _ = strconv.ParseFloat(tmp[10].(string), 64)
	r.Data1 = tmp[11].(string)

    return nil
}

func main() {
	symbols := []string{
		"ADAEUR", 
		"LINKEUR",
		"LTCEUR",
		"BNBEUR",
		"DOGEEUR", 
		"XLMEUR", 
		"XRPEUR", 
		"DOTEUR",
		"BTCEUR",
		"ETHEUR"}

	var report []string

	for _, symbol := range symbols {
		url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=1h&limit=500", symbol)
		resp, err := http.Get(url)
		check(err)

	 	body, err := ioutil.ReadAll(resp.Body)
	 	check(err)

   		var klines []Kline
   		json.Unmarshal(body, &klines)

   		var sma7 []float64
   		var sma25 []float64
   		var sma99 []float64

   		for idx, _ := range klines {
   			idxPlusOne := idx + 1
   			sum7 := 0.0
   			sum25 := 0.0
   			sum99 := 0.0
   			
   			if (idxPlusOne - 7) >= 1 {
   				for _, k := range klines[idx-7:idx] {
   					sum7 += k.Close
   				}

   				sma7 = append(sma7, (sum7 / 7))
   			}

   			if (idxPlusOne - 25) >= 1 {
   				for _, k := range klines[idx-25:idx] {
   					sum25 += k.Close
   				}

   				sma25 = append(sma25, (sum25 / 25))
   			}

   			if (idxPlusOne - 99) >= 1 {
   				for _, k := range klines[idx-99:idx] {
   					sum99 += k.Close
   				}

   				sma99 = append(sma99, (sum99 / 99))
   			}
   		}

   		lKlines := klines[len(klines)-1]
   		lSMA7 := sma7[len(sma7)-1]
   		lSMA25 := sma25[len(sma25)-1]
   		lSMA99 := sma99[len(sma99)-1]

   		if (lKlines.Close <= lSMA7)  && (lKlines.Close >= lSMA25) {
   			report = append(report, reportFormatter(symbol, "BETWEEN MA(7) AND MA(25)"))
   		}

   		if (lKlines.Close <= lSMA25)  && (lKlines.Close >= lSMA99) {
   			report = append(report, reportFormatter(symbol, "BETWEEN MA(25) AND MA(99)"))
   		}

   		red := 0
   		for _, k  := range klines[len(klines)-7:len(klines)-2] {
   			if k.Open >= k.Close {
   				red += 1
   			} else {
   				red = 0
   			}
   		}

   		if(red >= 3){
   			report = append(report, reportFormatter(symbol, "RED >= 3"))
   		}

   		oneDayPriceChangePercentage := ((klines[len(klines)-1].Close - klines[len(klines)-26].Close) / klines[len(klines)-26].Close) * 100

   		if oneDayPriceChangePercentage < 0.0 {
   			report = append(report, reportFormatter(symbol, fmt.Sprintf("PRICE CHANGE 24H %f%%", oneDayPriceChangePercentage)))
   		}

   		maxClose := 0.0
   		minClose := 0.0
   		for _, k  := range klines[len(klines)-30:len(klines)-3] {
   			if minClose == 0.0 {
   				minClose = k.Close
   			}

   			if k.Close > maxClose {
   				maxClose = k.Close
   			}

   			if k.Close < minClose {
   				minClose = k.Close
   			}

   		}

   		priceMinMaxPercentage := (lKlines.Close-minClose) / (maxClose-minClose) * 100

   		if priceMinMaxPercentage <= 65 {
   			report = append(report, reportFormatter(symbol, fmt.Sprintf("PRICE MIN MAX 24H %f%%", priceMinMaxPercentage)))
   		}
	}

	for _, r := range report{
		fmt.Println(r)
	}
}