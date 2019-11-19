# thesaurus
[![CircleCI](https://circleci.com/gh/jakewarren/thesaurus.svg?style=shield)](https://circleci.com/gh/jakewarren/thesaurus)
[![GitHub release](http://img.shields.io/github/release/jakewarren/thesaurus.svg?style=flat-square)](https://github.com/jakewarren/thesaurus/releases])
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/jakewarren/thesaurus)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/thesaurus/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakewarren/thesaurus)](https://goreportcard.com/report/github.com/jakewarren/thesaurus)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)

A command-line thesaurus app for Oxford.

![](screenshot.jpg)
<img src="thesaurusrex.jpg" width="300"/>

## Installation

If you have a working Go environment, you can install via `go get`:
```shell
go get github.com/jakewarren/thesaurus/...
```

## Configuration

Uses the configuration scheme used by [Rican7/define](https://github.com/Rican7/define).

### Configuration file

A configuration file can be stored at `~/.define.conf.json` and **thesaurus** will automatically load the values specified there.

### Environment variables

Some configuration values can also be specified via environment variables. This is especially useful for API keys of different sources.

The following environment variables are read by **thesaurus**'s sources:

- `OXFORD_DICTIONARY_APP_ID`
- `OXFORD_DICTIONARY_APP_KEY`

## Acknowledgements

Very heavily based on:  
* [travis-g/dict](https://github.com/travis-g/dict)  
* [Rican7/define](https://github.com/Rican7/define)  


## Obtaining API keys

The following are links to register for API keys for the different sources:

- [Oxford Dictionaries API](https://developer.oxforddictionaries.com/?tag=#plans)
