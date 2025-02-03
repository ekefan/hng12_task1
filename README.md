# HNG12: Stage_1 Task

In this task I developed a public api that returns the following information in JSON format:

1. A Number, passed to the api through uri params
2. Bool showing if it's a prime number
3. Bool showing if it's a perfect number
4. Array holding properties of the number such as odd, or even or armstrong
5. Fun Fact about the number, gotten from [numbersapi](http://numbersapi.com)

## API Specification

Endpoint `GET** <the deployed url domain>/api/classify-number?number=371`
Required JSON Response Format (200 OK):

```json
{
    "number": 371,
    "is_prime": false,
    "is_perfect": false,
    "properties": ["armstrong", "odd"],
    "digit_sum": 11,  // sum of its digits
    "fun_fact": "371 is an Armstrong number because 3^3 + 7^3 + 1^3 = 371" //gotten from the numbers API
}
```

Required JSON Response Format (400 Bad Request)

```json
{
    "number": "alphabet",
    "error": true
}
```
