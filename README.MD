# CSV to JSON converter

There are times we stuck at import data from CSV to JSON storage (eg: elasticsearch)

Problem is CSV has no correct format. So I decide to write my own tools to convert csv and import as JSON.

### Get it!

```shell
$ go get github.com/liemlhd3893/csv2json
```

### Install it!

```shell
$ go install github.com/liemlhd3893/csv2json
```

### How to use

1. We need a definition about how CSV was organized. Here the sample config file

```hcl
# Where csv will be organized
root = "."
# Where we want ouput our data
out_directory = "./out"

# Relative directory with root path
directory "user_action" {
    # File match this will be included
    include = [ ".*" ]
    # File match this pattern will be excluded
    exclude = [  ]

    /* 
    Definition of first column in CSV file.
        Naming will be your friend as the document to share with you co-worker (and your future-self)
    */
	column "a" {
        // Type of data. Current only support Int, Float, Boolean, String. You can easily extends by add more Parser. See `parser` package.
		type = "Int"
        // The JSON Path
		path = "a.s.a"
	}
	column "b" {
		type = "String"
		path = "a.s.b"
        // If this `true`. Skip this column
		skip = true
	}	
	column "c" {
		type = "String"
		path = "a.s.c"
	}				
	column "d" {
		type = "String"
		path = "a.s.ip"
        // Default value for this file. This convinient to deal with expanded CSV over time.
		default = "127.0.0.1"
	}

    // This column will always be added into final JSON.
	additional_column "type" {
		type = "String"
		path = "a.type"
		default = "PING"
	}
}


directory "user_info" {
    include = [ ".*" ]
    exclude = [  ]
	separator = "\t"
	column "a" {
		type = "Int"
		path = "a.s.a"
	}
	column "b" {
		type = "String"
		path = "a.s.b"
		skip = true
	}	
	column "c" {
		type = "String"
		path = "a.s.c"
	}				
	column "d" {
		type = "Boolean"
		path = "a.s.d"
	}
	column "d" {
		type = "String"
		path = "a.s.ip"
		default = "127.0.0.1"
	}		
	additional_column "type" {
		type = "String"
		path = "a.type"
		default = "PING"
	}
}
```

For more information. Take a look at [samples](https://github.com/liemlhd/csv2json/sample) package.