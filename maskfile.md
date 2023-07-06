# Tasks for TropesToGo
> After installing mask utility, run command:  
> **mask \<task\>**  
> where \<task\> is any of the following

## run
> Command for running TropesToGo CLI and be able to scrape data on TvTropes
~~~sh
cd tropestogo
go run ./main.go
~~~

### scrape
> Command for scraping data with the TropesToGo CLI

**OPTIONS**
* format
  * flags: -f --format
  * type: string
  * desc: The format of the generated dataset, CSV or JSON\
* limit
  * flags: -l --limit
  * type: number
  * desc: Limit the crawled works; is ignored if the -a flag is passed
* media
  * flags: -m --media
  * type: string
  * desc: The media type from which to scrape works
* output
  * flags: -o --output
  * type: string
  * desc: The name of the output dataset
* all
  * flags: -a -all
  * desc: Scrape all works in the media type

~~~sh
cd tropestogo
if [[ ! -z "$format" ]]; then 
    format="-f ${format}" 
else
    format=""
fi

    
if [[ ! -z "$media" ]]; then 
    media="-m ${media}"
else
    media=""
fi


if [[ ! -z "$output" ]]; then 
    output="-o ${output}"
else
    otuput=""
fi

if [[ ! -z "$limit" ]]; then
    limit="-l ${limit}"
else
    limit=""
fi

if [[ $all == "true" ]]; then
    go run ./main.go scrape -a $format $media $output $limit
else
    go run ./main.go scrape $format $media $output $limit
fi
~~~

### update
> Command for updating an scraped dataset with the TropesToGo CLI

**OPTIONS**
* dataset
  * flags: -d --dataset
  * type: string
  * desc: Dataset name to update

~~~sh
cd tropestogo
go run ./main.go update -d $dataset
~~~

## build
> Command for building the project
~~~sh
echo "Building project..."
mask build doc
mask build code
~~~

### build doc
> Builds a pdf with the project's documentation
~~~sh
echo "Building project's documentation..."
cd doc
make
~~~

### build code
> Builds the project's code without installing packages
> checking that the packages can be built
~~~sh
echo "Building project's code..."
cd tropestogo
go build -v ./...

if [ $? -eq 0 ]; then
    echo "Code built successfully"
    exit 0
else
    echo "There's been an error building the project"
    exit 1
fi
~~~

## test
> Executes all the tests of the project
~~~sh
cd tropestogo
go test -v ./...
~~~