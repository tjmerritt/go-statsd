# go-statsd
Yet another Statsd Client for GoLang

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/tjmerritt/go-statsd)
[![license](http://img.shields.io/badge/license-BSD3-red.svg?style=flat)](http://opensource.org/licenses/BSD-3-Clause)

https://godoc.org/github.com/tjmerritt/go-statsd

# Features
* Support for all standard statsd types
* Metric aggregation to minimize UDP traffic
* Thread specific stats collector for optimum performance within a single thread
* Thread safe stats collector for ease of use and peformance with rapid thread generation
