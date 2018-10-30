# Heresy

Those who can’t get used to the absence of generics will commit heresy unless Go2 is released.

JSON sanitizer removes personally identifiable information from a dataset:

```sh
go run sanitizer.go input.json → output.json with encrypted (removed) id, names, phones
```
![](https://github.com/alissiawells/Heresy/blob/master/anonymization.jpeg)
