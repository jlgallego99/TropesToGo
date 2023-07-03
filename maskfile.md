# Tasks for TropesToGo
> After installing mask utility, run command:  
> **mask \<task\>**  
> where \<task\> is any of the following

## run
> Command for running TropesToGo CLI and be able to scrape data on TvTropes
~~~sh
cd tropestogo
go run ./...
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