# Heresy

Those who can’t get used to the absence of generics will commit heresy unless Go2 is released.

JSON sanitizer removes personally identifiable information from a dataset:

```sh
go run sanitizer.go input.json → output.json with encrypted (removed) id, names, phones
```


A TCP proxy server ~~for some deviant Chinese frameworks~~ that listens for TCP pockets on localhost, 

looks up packet IP destination, reads from it and sends the reply:
```sh
sudo iptables -t nat -A OUTPUT -p tcp -m tcp --dport 443 -j REDIRECT --to-ports 1111
sudo iptables -t nat -A OUTPUT -p tcp -m tcp --dport 80 -j REDIRECT --to-ports 1111
sudo useradd tcpprunner
go build tcpproxy.go
sudo iptables -t nat -A OUTPUT -m tcp -p tcp --dport 80 -m owner --uid-owner tcpprunner -j RETURN
sudo iptables -t nat -A OUTPUT -m tcp -p tcp --dport 443 -m owner --uid-owner tcpprunner -j RETURN
sudo -u tcpprunner ./tcpproxy
```

![](https://github.com/alissiawells/Heresy/blob/master/anonymization.jpeg)
