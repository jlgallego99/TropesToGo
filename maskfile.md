# Tasks for TropesToGo
> After installing mask utility, run command:  
> **mask \<task\>**  
> where \<task\> is any of the following

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
else
    echo "There's been an error building the project"
fi
~~~

## test
> Executes all the tests of the project
~~~sh
cd tropestogo
go test -v ./...
~~~