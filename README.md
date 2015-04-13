# Address Autocomplete web service
This is a proof-of-concept address autocomplete web service written in Go and backed by Elasticsearch that can support "typeahead" for address search fields in applications. Read [my blog post] for more details and check out [the demo](http://davewalk.github.io/address-autocomplete).  
### Installation
You'll need [Go](http://golang.org/doc) and [ElasticSearch](http://www.elasticsearch.org).  

Get this repo (don't clone!):  

    go get -d davewalk.net/address-autocomplete

To create the index, custom analzyer and `address` mapping:  

    cd $GOPATH/src/github.com/davewalk/address-autocomplete
    ./index.sh

And to populate the database:  

    export AUTOCOMPLETE_FILE=/home/user/path/to/csv/file
    cd index/
    go build
    ./index

and make yourself a sandwich because that'll take a while.  

NOTE: If your Elasticsearch instance isn't running at `localhost:9200`, you'll have to change references to it in the `main.go` and `index/main.go` files.  

### Starting the app

From the main directory:  

    go build
    export AUTOCOMPLETE_PORT=8000 [or whatever]
    ./address-autocomplete

### Usage

The resource is at `/autocomplete` and it takes two query parameters: `q` and `num` (optional). `q` is a string that you are trying to get address suggestions for and `num` is the number of suggestions that you want (25 max, 10 default). For example:
    
    http://localhost:8000/autocomplete?q=1234 m&q=15

### Future
The API can be made more effective by:
1. Street name aliasing (ie. "NINTH" -> "9TH")
2. Matching on misspellings (ie. "MAKRET" -> "MARKET") (via n-grams?)
3. Weighting of returned results by some additional attribute of each address (street type? ie. 
"minor" vs. "major" streets). This requires GIS to do and I didn't feel like it at the time. If someone wanted to attach street centerline classes to each row in the data and send me a pull request that would be awesome! Pretty please?  

Pull requests accepted!

### License
MIT
