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
> Builds the project's code
~~~sh
echo "Building project's code..."
go version
~~~

## test
> Executes all the tests of the project
~~~sh
go version
~~~

### test docker
> Executes all the tests of the project on the docker container
~~~bash
docker run jlgallego99/tropestogo
~~~