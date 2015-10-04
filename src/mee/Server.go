package main

import (
	"encoding/json"
	"github.com/gorilla/rpc"
	myjson "github.com/gorilla/rpc/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type StockArgs struct {
	StockSymbolandPercentage string
	Budget                   float32
}

type StockReply struct {
	Tradeid        int
	Stocks         string
	UnvestedAmount float32
}

type CheckArgs struct {
	Tradeid int
}
type CheckReply struct {
	Stocks           string
	UnvestedAmount   float32
	CurrentMarketVal float32
}
type Inputs struct {
	Company string
	Percent float64
	Stocks  int
	Rate    float64
}

type Record struct {
	Stocks         string
	UnvestedAmount float32
}

type FromApi1 struct {
	FromApi2 struct {
		C1       int    `json:"count"`
		C2       string `json:"created"`
		C3       string `json:"lang"`
		FromApi3 struct {
			FromApi4 struct {
				Value string `json:"Ask"`
			} `json:"quote"`
		} `json:"results"`
	} `json:"query"`
}

var M map[int]Record

var cntr int

type StockService struct{}
type CheckPortfolioService struct{}

func (h *StockService) Say(r *http.Request, args *StockArgs, reply *StockReply) error {
	var basic_url = "http://query.yahooapis.com/v1/public/yql?env=store://datatables.org/alltableswithkeys&format=json"
	var yql = "&q=select%20Ask%20from%20yahoo.finance.quotes%20where%20symbol%20=%20"
	s := strings.Split(args.StockSymbolandPercentage, ",")

	count := len(s)
	i := 0
	det := make([]Inputs, count)
	var v string
	var total float64
	total = 0.0

	for i < count {
		k := strings.Split(s[i], ":")
		v = k[1]
		v1 := v[0 : len(v)-1]
		det[i].Company = k[0]
		l, err := strconv.ParseFloat(v1, 32)
		if err != nil {
			// Invalid string
		}
		det[i].Percent = float64(l)
		var final = basic_url + yql + "%22" + det[i].Company + "%22"
		res, err := http.Get(final)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		type JsonResponse FromApi1
		var jsonReturn JsonResponse
		json.Unmarshal(body, &jsonReturn)

		if res.StatusCode != 200 {
			log.Fatal("Unexpected status code", res.StatusCode)
		}
		rs := jsonReturn.FromApi2.FromApi3.FromApi4.Value
		det[i].Rate, err = strconv.ParseFloat(rs, 64)
		if err != nil {
			// Invalid string
		}

		amt := (det[i].Percent / 100) * (float64(args.Budget))

		tmp := amt / det[i].Rate
		det[i].Stocks = int(tmp)
		total = total + (float64(det[i].Stocks) * det[i].Rate)
		i = i + 1
	}
	Remain := float64(args.Budget) - total
	reply.UnvestedAmount = float32(Remain)
	i = 0

	var out string
	out = ""
	for i < len(det) {
		p1 := strconv.Itoa(det[i].Stocks)
		p2 := strconv.FormatFloat(det[i].Rate, 'f', 2, 64)

		out = out + det[i].Company + ":" + p1 + ":$" + p2 + ","
		i = i + 1
	}

	out = out[0 : len(out)-2]

	reply.Stocks = out
	var re Record
	re.Stocks = out
	re.UnvestedAmount = float32(Remain)
	cntr = cntr + 1
	if len(M) == 0 {
		M = make(map[int]Record)

	}
	M[cntr] = re
	reply.Tradeid = cntr

	return nil
}

func (h *CheckPortfolioService) Che(r *http.Request, args *CheckArgs, reply *CheckReply) error {

	var re Record
	re = M[args.Tradeid]
    log.Println(re)
	reply.UnvestedAmount = re.UnvestedAmount

	var basic_url = "http://query.yahooapis.com/v1/public/yql?env=store://datatables.org/alltableswithkeys&format=json"
	var yql = "&q=select%20Ask%20from%20yahoo.finance.quotes%20where%20symbol%20=%20"
	s := strings.Split(re.Stocks, ",")
	var out = ""
	count := len(s)
	i := 0

	var total float64
	total = 0.0

	for i < count {
		k := strings.Split(s[i], ":")
		com := k[0]

		v := k[2]

		v1 := v[1:len(v)]

		rate, err := strconv.ParseFloat(v1, 64)
		if err != nil {
			// Invalid string
		}

		var final = basic_url + yql + "%22" + com + "%22"
		res, err := http.Get(final)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		type JsonResponse FromApi1
		var jsonReturn JsonResponse
		json.Unmarshal(body, &jsonReturn)

		if res.StatusCode != 200 {
			log.Fatal("Unexpected status code", res.StatusCode)
		}
		rs := jsonReturn.FromApi2.FromApi3.FromApi4.Value
		r, err := strconv.ParseFloat(rs, 64)
		if err != nil {
			// Invalid string
		}
		
		var sym string

		if rate > r {
			sym = "-"
		} else if rate < r {
			sym = "+"
		} else {
			sym = ""
		}
		currates := strconv.FormatFloat(r, 'f', 2, 64)
		out = out + com + ":" + k[1] + ":" + sym + "$" + currates + ","
		temp, err := strconv.ParseFloat(k[1], 64)
		total = total + (temp * r)
		i = i + 1
	}
	reply.Stocks = out
	total32 := float32(total)
	reply.CurrentMarketVal = total32
	return nil
}

func main() {
	jsonRPC := rpc.NewServer()
	jsonCodec := myjson.NewCodec()
	jsonRPC.RegisterCodec(jsonCodec, "application/json")
	jsonRPC.RegisterCodec(jsonCodec, "application/json; charset=UTF-8") // For firefox 11 and other browsers which append the charset=UTF-8
	jsonRPC.RegisterService(new(StockService), "")
	jsonRPC.RegisterService(new(CheckPortfolioService), "")
	log.Print("this is calling http")
	http.Handle("/rpc", jsonRPC)
	http.ListenAndServe(":8080", jsonRPC)
	log.Print("http : Listen And Service")

}
