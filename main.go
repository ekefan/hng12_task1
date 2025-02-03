package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const reqParam = "number"
const requestURL = "http://numbersapi.com"

type Resp struct {
	Number     int      `json:"number"`
	IsPrime    bool     `json:"is_prime"`
	IsPerfect  bool     `json:"is_perfect"`
	Properties []string `json:"properties"`
	DigitSum   int      `json:"digit_sum"`
	FunFact    string   `json:"fun_fact"`
}
type ErrResp struct {
	Number string `json:"number"`
	Error  bool   `json:"error"`
}

func main() {

	http.HandleFunc("/api/classify-number", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handleCors(w)
			return
		}

		// Parse request query
		requestQuery, err := url.Parse(r.RequestURI)
		if err != nil {
			slog.Error("can not parse request uri", "details", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		queryParam, err := url.ParseQuery(requestQuery.RawQuery)
		if err != nil {
			slog.Error("can not parse request query", "details", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		numStr := queryParam.Get(reqParam)
		if numStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp := ErrResp{
				Number: numStr,
				Error:  true,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Validate request params
		number, err := strconv.Atoi(numStr)
		if err != nil {
			slog.Error("Error converting reqParam to number", "details", err)
			w.WriteHeader(http.StatusBadRequest)
			resp := ErrResp{
				Number: numStr,
				Error:  true,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		// write response
		funFact, err := getFuncFact(number)
		if err != nil {
			slog.Error("could not get fun fact", "details", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := Resp{
			Number:     number,
			IsPrime:    isPrime(number),
			IsPerfect:  isPerfect(number),
			Properties: getProperties(number),
			DigitSum:   getDigitSum(number),
			FunFact:    funFact,
		}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusOK)
	})
	err := http.ListenAndServe(":8080", nil)
	slog.Error("Error starting service", "details", err)
}

func handleCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func isPerfect(n int) bool {
	if n <= 1 {
		return false
	}
	sum1 := 0
	for i := 1; i < n; i++ {
		if n%i == 0 {
			sum1 += i
		}
	}
	return sum1 == n
}

func getDigitSum(num int) int {
	if num < 0 {
		num *= -1
	}
	res := 0
	for num > 0 {
		res += num % 10
		num /= 10
	}
	return res
}

func isArmStrong(n int) bool {
	tmp := n
	result := 0
	count := 0
	for i := n; i > 0; i /= 10 {
		count++
	}
	for tmp > 0 {
		remainder := tmp % 10
		result += int(math.Pow(float64(remainder), float64(count)))
		tmp /= 10
	}
	return result == n
}

func getFuncFact(num int) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%v", requestURL, num))
	if err != nil {
		slog.Error("Error requesting resource from url", "details", err)
		return "", err
	}
	defer resp.Body.Close()

	reqBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading resp body", "details", err)
		return "", err
	}
	return string(reqBody), nil
}

func getProperties(n int) []string {
	retme := []string{}
	if isArmStrong(n) {
		retme = append(retme, "armstrong")
	}
	if n%2 == 0 {
		retme = append(retme, "even")
	} else {
		retme = append(retme, "odd")
	}
	return retme
}
