# index-replicate

[![Build Status](https://github.com/brogand93/index-replicate/actions/workflows/go-tests.yml/badge.svg)](https://github.com/brogand93/index-replicate/actions/workflows)

index-replicate allows you to replicate an index ETF. It allows you to create a custom index to use as the basis of a stock portfolio.

This form of do-it-yourself indexing is pointless for many people. It's going to be a lot simpler to simply invest in an ETF that tracks the index you wish to track. If however, you are a resident of EU this may not be possible without significant tax hurdles getting in your way.

This small CLI aims to get a close approximation of a chosen index (or percentage of an index) at the previous trading day's closing prices for a given amount of money. The source this data is [svcga.com/sc/index](https://svcga.com/sc/index).

## Installing

Download the correct release for your system from the [releases section](https://github.com/brogand93/index-replicate/releases)

## Running the Index Replicater

### Help Page

```bash
$ index-replicate -h

Index Replicate allows you to replicate an index

Usage:
  index-replicate [flags]

Flags:
  -c, --contribution float   amount you wish to invest in the index
  -h, --help                 help for index-replicate
  -i, --index string         index to replicate [sp500, dowjones, nasdaq100] (default "sp500")
      --output-to-csv        output to a local csv file
  -p, --percentage int       percentage of the index you wish to replicate (default 100)
  -r, --round                round share buy quantity to the nearest whole share
```

### Example For S&P 500

If you wanted to replicate the top 70% of the S&P 500, and you wanted to invest 10,000 you could run the following:

```bash
$ index-replicate -i sp500 --percentage 70 --contribution 10000
+---------------------------------------------+--------+----------+--------+-----------+
|                    NAME                     | SYMBOL | QUANTITY | WEIGHT |   VALUE   |
+---------------------------------------------+--------+----------+--------+-----------+
| Apple Inc.                                  | AAPL   |     6.91 | 8.37 % | 836.84 $  |
| Microsoft Corporation                       | MSFT   |     3.27 | 7.70 % | 769.73 $  |
| Amazon.com Inc.                             | AMZN   |     0.18 | 5.70 % | 570.13 $  |
| Facebook Inc. Class A                       | FB     |     1.05 | 2.83 % | 282.68 $  |
| Alphabet Inc. Class A                       | GOOGL  |     0.13 | 2.71 % | 271.14 $  |
| Alphabet Inc. Class C                       | GOOG   |     0.13 | 2.64 % | 263.56 $  |
...
```
