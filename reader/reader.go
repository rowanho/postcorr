package reader

import (
    "io/ioutil"
)


func check(e error) {
    if e != nil {
        panic(e)

    }
}

/**
* Function readfile
* parameter filename - the name of the file's contents to read
* returns string - the file's contents
*/

func ReadFile(filename string) string{
    dat, err := ioutil.ReadFile(filename)
    check(err)
    return string(dat)
}
