<h1>About:</h1>
This tool is created to see your domain's SSL certificates expiry date. It allows you configuring email alerting system, so you don't miss your certificate's expiry date. You can check domains one by one or provide a file with all of your domains so the app can check it automatically.
<h1>Usage:</h1>
<h3>Command examples:</h3>

This will display the expiration date of the certificate on www.foobar.com:
```shell
./ssl-check -d www.foobar.com
```
Providing a file with domains instead of passing them one by one:
```shell
./ssl-check -f myDomains.txt 
```
Passing -l flag will work the same way as the previous ones but it will also display the number of days left before the certificate expires:
```shell
./ssl-check -d www.foobar.com -l
```
Or:
```shell
./ssl-check -f myDomains.txt -l
```
