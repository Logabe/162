# 162 Web Server
This is a web server that I wrote in Go for a personal project. It uses a central index file to locate and serve files and has support for serving dynamic content via dedicated binary files.

## Index.json
An index.json file follows this format:
```
{
    //The main array
    "content" [
        {
            //File location
            "location": "text/hello_world.txt",
            //MIME type
            "content-type": "text/plain",
            //Should the file be displayed or downloaded?
            "content-disposition": "inline"
        }, {
            "location": "data/downloadme.zip",
            "content-type": "application/zip",
            //Specifying download name
            "content-disposition": "attachment; filename=\"resource.zip\""
        },{
            "location": "bin/executablefile",
            //What type will the returned content be?
            "content-type": "text/html",
            "content-disposition": "inline",
            "execution-method": "execute",
            //Any information you want to pass to the executable
            //At present only 'addr' works, but more can be added.
            "execution-data": ["addr"]
        },
    ]
}```