# Use rtltcp as a source of random numbers and serve them over http/json

Based on excellent [rtltcp](https://github.com/bemasher/rtltcp) and inspired by [rtl-entropy](https://github.com/pwarren/rtl-entropy)   

## Examples:
### Get a random 10 byte number
HTTP Get: http://127.0.0.1:8080/?l=10   
Returns: {"hash": "74437db7eea99d484788"}   

