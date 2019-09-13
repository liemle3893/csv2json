# CSV to JSON converter

There are times we stuck at import data from CSV to JSON storage (eg: elasticsearch)

Problem is CSV has no correct format. So I decide to write my own tools to convert csv and import as JSON.

### Get it!

```shell
$ go get github.com/liemle3893/csv2json
```

### Install it!

```shell
$ go install github.com/liemle3893/csv2json
```

### How to use

1. We need a definition about how CSV was organized. Here the sample config `config.hcl`

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
    // Include all file but except *.exe
    include = [ ".*" ]
    exclude = [ "*.exe" ]
    // Set file separator to tab
	separator = "\t"
    column "exclude" {
        type = "String"
        // Record at this index with value 'test' will be skipped
        excludes = ["test"]
    }
    column "Indexed" {
        type = "Indexed"
        default = "undefied"
        // Map
        // web => Browser
        // ios, android => Mobile
        // other => Unknown
        indices = {"undefied":"Unknown" "web"="Browser", "ios"="Mobile","android"="Mobile"}
    } 
}
```

2. And now we can have our result

```shell
$ csv2json -c config.hcl
```

For more information. Take a look at [samples](https://github.com/liemle3893/csv2json/blob/master/samples) package.
