<h1>About:</h1>
This tool is created to see your domain's SSL certificates expiry date. It allows you configuring email alerting system, so you don't miss your certificate's expiry date. You can check domains one by one or provide a file with all of your domains so the app can check it automatically.
<h1>Usage:</h1>
<h3>Flags:</h3>

This will display the expiration date of the certificate on www.foobar.com:
```shell
./ssl-check -d www.foobar.com
```
Providing a file with domains instead of passing them one by one:
```shell
./ssl-check -f myDomains.txt 
```
Displaying number of days left before the certificate will expire:
```shell
./ssl-check -d www.foobar.com -l
```
```shell
./ssl-check -f myDomains.txt -l
```
Providing a config file for email notifications:
```shell
./ssl-check -f myDomains.txt -l -e config.json
```
```shell
./ssl-check -d www.foobar.com -l -e config.json
```
<h3>Sample file with domains:</h3>

```txt
//You have to pass domains line by line
www.foo.com
www.bar.com
www.foo.org
www.bar.org
```
<h3>Sample file with email config:</h3>

```js
{
	"sender_email": "foo@bar.com", // your email address
	"sender_password": "foobar", // your password for the email
	"smtp_host": "smtp.foobar.com", // smtp server from which you want to send the message
	"smtp_port": "587", // port on which the server is exposed
	"receivers": [ // list of addresses where you want to send the message
		"foo@bar.com", 
		"foo1@bar.com"
	],
	"alert_before": 14 // will send an email 14 days before the expiration
}
```
NOTE: If you are not using your local smtp server you have to make sure to allow the access for third party applications in your mail account.