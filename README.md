# BrowserScrape

Neat little tool to use your very own browser for webscraping. **Early beta**

[![IMAGE ALT TEXT HERE](https://img.youtube.com/vi/RvDoHGD06vY/0.jpg)](https://www.youtube.com/watch?v=RvDoHGD06vY)


## How to setup

1. Clone Repo
2. Open "chrome://extensions"
3. Load unpacked extension from "extension" folder of this repo
4. Copy uuid and you are good

## Disclaimers

- This is a very early beta
- Default extension runs with MY (Simons) Server. You HAVE TO TRUST ME. As I could go rouge and just let your browser open whatever url I want. Self-host the server if you want to be ultra-secure
- It only works when the extension is connected to the server. If your browser is closed, NO SCRAPY SCRAPY!

## FAQ

> I'm not a programmer or github user, can I simply download the extension folder from Github and follow the rest of the instructions or does 'cloning' do something more?

Yes, you can also just download the folder as zip from github.

> Is it running code on your server?

Yes, I run an instance of the "server/" (which code also is in the folder). The server actually connects the domain https://simon.red/browserscrape/ with your browser. Without that, it would not be possible to connect to your browser instance.

> How often does the UUID change

It is stored in the local storage of your browser, so it should only change when you reinstall the extension or click the "regenerate UUID" button in the extension. No guarantee though, it might get changed with you deleting browser history as well

>  is it related to this browser instance?

Yes, the uuid is randomly generated on the first time you open the extension and then every time you click the "regenerate uuid" button.

> I guess that is why so many of the Chrome store extensions that do something similar have been removed.

Yeah, I think there is no chance to get this into the chrome store.

> possible expensive to run if you had a lot of users

I think that should be doable, as not much processing happens on the server. It basically just pipes text from one end to another.  Would say for thousand of users, one server is enough (<1$/month)

> Is there a chance that a site would see it as logging in to an account in another location and get one banned for sharing account logins? ie: could it look like the user is logged in from two different locations?

There is NO chance that the site detects this as "logged in from two different locations". As it is not, it logs in from the same location.It MIGHT be possible that they detect it as bot traffic as you are directly opening a page, that you usually would navigate to, or if you e.g. open the page 9 times per minute it also might be a non-natural usage pattern. But the risk is that same as if you would use local browser automation. 
