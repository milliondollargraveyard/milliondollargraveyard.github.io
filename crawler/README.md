# Crawler

Crawl sites from the milliondollarhomepage and annotate with response status.

## Usage

```
$ crawler --help
Usage of ./crawler:
  -c uint
    	concurrency (default 1)
  -limit int
    	abort after processing this many lines (default -1)
  -skip uint
    	skip this many lines before processin

$ crawler < input.jsons > output.jsons
```

Given an input that [looks like
this](https://gist.github.com/shazow/79a8129c8f2d91442c2f08165f0264c3#file-input-jsons):

```json
{"coords": "630,310,640,320", "href": "http://www.getpixel.net/", "title": "getpixel.net, stock photography"}
{"coords": "850,50,860,60", "href": "http://www.mynewbritain.com/", "title": "MyNewBritain.com"}
{"coords": "390,280,420,310", "href": "http://www.pandasoftware.com/", "title": "PC infected? Free Spyware Scan - PandaSoftware.com"}
{"coords": "690,560,700,570", "href": "http://www.frozenweb.co.uk/", "title": "FrozenWeb.co.uk UK Web Hosting Specialists"}
{"coords": "610,150,640,170", "href": "http://www.sillyant.com/?src=1M", "title": "SillyAnt - something for your cellphone and PDA"}
{"coords": "160,160,290,190", "href": "http://www.rentclicks.com/", "title": "Home Rentals, Homes for Rent, and Apartments"}
{"coords": "210,80,220,90", "href": "http://www.kickbuttideas.com/12.php?p=1000&a=extramoney", "title": "Kick Butt Ideas for Making Money"}
{"coords": "320,70,330,80", "href": "http://www.fastminimoto.co.uk/", "title": "Mini Motos / Dirt Bikes / Quads / Parts *BARGAIN*"}
{"coords": "730,0,740,10", "href": "http://www.hamsterland.com/million.asp", "title": "THE GREEN DOT Hamsterland, hamster, cage, feed"}
{"coords": "520,0,590,10", "href": "http://localtap.net/", "title": "Localtap.com - LOCAL BREWED BEER DELIVERED DIRECT TO YOU"}
```

Produce an output that [looks like
this](https://gist.github.com/shazow/79a8129c8f2d91442c2f08165f0264c3#file-output-jsons):

```json
{"href":"http://www.getpixel.net/","coords":"630,310,640,320","title":"getpixel.net, stock photography","response":{"status":200}}
{"href":"http://www.hamsterland.com/million.asp","coords":"730,0,740,10","title":"THE GREEN DOT Hamsterland, hamster, cage, feed","response":{"status":200,"title":"Hamster Land: Homepage, hamsters, hamster cage, feed"}}
{"href":"http://www.dreamwords.com/","coords":"420,300,440,330","title":"TOM CORVEN: Free Fun Fantastic Fiction","response":{"status":200,"size":5622,"title":"Dreamwords - Paul Story"}}
{"href":"http://localtap.net/","coords":"520,0,590,10","title":"Localtap.com - LOCAL BREWED BEER DELIVERED DIRECT TO YOU","response":{"status":200,"title":"localtap.net Is For Sale"}}
{"href":"http://www.eriqx.com/","coords":"310,90,320,100","title":"EriqX[.com]","response":{"status":200,"title":"EriqX[.com]"}}
{"href":"http://www.frozenweb.co.uk/","coords":"690,560,700,570","title":"FrozenWeb.co.uk UK Web Hosting Specialists","response":{"status":200,"title":"Portal Home - FrozenWeb"}}
{"href":"http://www.mynewbritain.com/","coords":"850,50,860,60","title":"MyNewBritain.com","response":{"status":200,"title":"My New Britain - The Magic Of Tourism"}}
{"href":"http://www.webgatehost.com/","coords":"470,180,490,190","title":"Web Hosting,Domain Name Registration,low cost host","response":{"status":200,"title":"Webgatehost.com"}}
{"href":"http://atyu.com/show/","coords":"450,90,460,100","title":"RESERVED FOR: DIE AND GO TO HELL!","response":{"error":"Get http://atyu.com/show/: read tcp 10.4.43.50:53752-\u003e184.168.221.96:80: read: connection reset by peer"}}
{"href":"http://www.homes-uk.co.uk/","coords":"0,110,10,120","title":"UK Property Server","response":{"status":200,"title":"Online Estate Agent"}}
```
