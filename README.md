<h1>About:</h1>
This tool is created to see your domain's SSL certificates expiry date. It allows you configuring email alerting system, so you don't miss your certificate's expiry date. You can check domains one by one or provide a file with all of your domains so the app can check it automatically.
<h1>Usage:</h1>
<h3>Flags:</h3>
<ul>
	<li>d -> provide one single domain.</li>
	<li>f -> provide a file with domains.</li>
	<li>l -> will additionally display a number of days until the certificate expires.</li>
	<li>e -> provide a JSON file with email config.</li>
</ul>
<h3>Examples:</h3>
<code>
	./ssl-check -d www.foobar.com
</code>
<code>
	./ssl-check -d www.foobar.com -l
</code>
<code>
	./ssl-check -f myDomains.txt 
</code>
<code>
	./ssl-check -f myDomains.txt -l
</code>
