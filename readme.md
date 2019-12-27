# offerzendata - A utility for grabbing Offerzen public profile data for analysis

`offerzendata` is a utility that allows you to download the Offerzen company 
public profiles for further data analysis.

# How to install and run

`offerzendata` is written in go, so you need go installed on your system. To 
install `offerzendata` is to simply run:

```sh
go get -u github.com/ShiraazMoollatjie/offerzendata
```

To run `offerzendata` you need to have an **authorization token**. You need to 
contact offerzen or the repo owner (me) directly regarding how to get a token
for using offerzendata. The token needs to be set as part of the `OFFERZEN_TOKEN`
environment variable.

```sh
 OFFERZEN_TOKEN=your_token ./offerzendata 
``` 