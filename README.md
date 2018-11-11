# Heresy

Those who can’t get used to the absence of generics will commit heresy unless Go2 is released.

### Multithreading scrapper

crawls from the start page parsing URLs which contain key words

* Installation
```sh
$ git clone git clone https://github.com/alissiawells/Heresy.git
$ cd Heresy
$ go run spider.go https://start_page depth language keyword1 keyword2 ... keywordN
```
* Dependences

[Stemmer](link) for key words (suports Russian, English and other languages)
```sh
$ go get github.com/kljensen/snowball
```

### TCP proxy server 

listens for TCP pockets, looks up packet IP destination, reads from it and sends the reply
```sh
$ sudo iptables -t nat -A OUTPUT -p tcp -m tcp --dport 443 -j REDIRECT --to-ports 1111
$ sudo iptables -t nat -A OUTPUT -p tcp -m tcp --dport 80 -j REDIRECT --to-ports 1111
$ sudo useradd tcpprunner
$ go build tcpproxy.go
$ sudo iptables -t nat -A OUTPUT -m tcp -p tcp --dport 80 -m owner --uid-owner tcpprunner -j RETURN
$ sudo iptables -t nat -A OUTPUT -m tcp -p tcp --dport 443 -m owner --uid-owner tcpprunner -j RETURN
$ sudo -u tcpprunner ./tcpproxy
```

### JSON sanitizer 

removes personally identifiable information from a dataset

```sh
$ go run sanitizer.go input.json → output.json with encrypted (removed) id, names, phones
```

![](https://github.com/alissiawells/Heresy/blob/master/anonymization.jpeg)
