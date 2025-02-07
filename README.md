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

## FAQ (Frequently Asked Questions)

### Installation & Setup
**Q: Do I need programming knowledge to install this extension?**

A: No, you can simply download the ZIP folder from GitHub and follow the installation instructions. No programming or GitHub experience is required.

**Q: How is the UUID (unique identifier) generated and managed?**

A: The UUID is randomly generated when you first install the extension. It's stored in your browser's local storage and only changes if you:
- Reinstall the extension
- Click the "regenerate UUID" button
- Possibly when clearing browser history/data

### Technical Details
**Q: Does this extension run on a remote server?**

A: Yes, there is a server component that connects your browser to the domain https://simon.red/browserscrape/. This server acts as a bridge between your browser instance and the extension.

**Q: What are the server costs and scalability?**

A: The server requirements are minimal since it primarily functions as a message relay. One server can handle thousands of users for less than $1/month, as it only pipes text between endpoints without heavy processing.

### Security & Detection
**Q: Can websites detect this as a shared or duplicate login?**

A: No, websites won't detect this as logging in from multiple locations because all traffic comes from your original browser location. However, be aware that:
- Unusual navigation patterns might be flagged as bot behavior
- Rapid, repeated page loads may appear suspicious
- The risk level is similar to using local browser automation tools

### Distribution
**Q: Why isn't this extension available in the Chrome Web Store?**

A: Due to Chrome Web Store's strict policies regarding browser automation and scraping tools, this extension cannot be distributed through their platform. This is common for similar tools. would navigate to, or if you e.g. open the page 9 times per minute it also might be a non-natural usage pattern. But the risk is that same as if you would use local browser automation. 
