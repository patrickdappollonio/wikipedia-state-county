## Wikipedia State & County List parsing

This mini projects reads all the states and counties [from the Wikipedia list](https://en.wikipedia.org/wiki/List_of_counties_by_U.S._state) and prints to `stdout` a MySQL code ready to be imported into any MySQL Database.

The sql file will contain all the 50 states in the US with counties normalized to remove from the names values such as `Municipality and County of` from the names. 

Additionally, there's also a full output file ready for you to use in the same repository under the name [output.sql](output.sql). Just import it into any database and you're good to go.

### Getting the program and executing it

Just issue a `go get` to this repo, and then install it or run it, you can pipe the output to `pbcopy` on Mac OS X to get it on your clipboard, like this:

```
# fetches the program
go get github.com/patrickdappollonio/wikipedia-state-county 

# move to the app's directory
cd $GOPATH/src/github.com/patrickdappollonio/wikipedia-state-county

# execute the program and send the output to the clipboard
go run main.go | pbcopy
```

### Forks? Issues? Improvements?

Always welcome. Just send them through!
