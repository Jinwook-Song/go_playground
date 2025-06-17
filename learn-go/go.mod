module github.com/jinwook-song/learn-go

go 1.24.3

replace github.com/jinwook-song/learn-go/banking => ./banking

replace github.com/jinwook-song/learn-go/my_dict => ./my_dict

replace github.com/jinwook-song/learn-go/urlchecker => ./urlchecker

require (
	github.com/PuerkitoBio/goquery v1.10.3 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	golang.org/x/net v0.39.0 // indirect
)
