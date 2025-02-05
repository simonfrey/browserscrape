const baseURL = "https://simon.red/browserscrape";

async function generateUUID() {
    const uuid = crypto.randomUUID();
    await chrome.storage.local.set({ uuid });
    return uuid;
}

async function getUUID() {
    const data = await chrome.storage.local.get('uuid');
    return data.uuid || await generateUUID();
}

async function capturePageContent(url) {
    try {
        const tab = await chrome.tabs.create({ url, active: false });
        await waitForPageLoad(tab.id);
        const [{ result }] = await chrome.scripting.executeScript({
            target: { tabId: tab.id },
            function: () => ({
                html: document.documentElement.outerHTML,
                title: document.title,
                url: window.location.href,
                isValid: document.body.textContent.length > 100 && document.title.length > 0
            })
        });
        await chrome.tabs.remove(tab.id);
        if (!result.isValid) throw new Error('Page content appears incomplete');
        return result;
    } catch (error) {
        console.error('Capture error:', error);
        throw error;
    }
}

async function waitForPageLoad(tabId) {
    return new Promise((resolve) => {
        const checkComplete = async () => {
            try {
                const [{ result }] = await chrome.scripting.executeScript({
                    target: { tabId },
                    function: () => document.readyState
                });
                if (result === 'complete') setTimeout(resolve, 2000);
                else setTimeout(checkComplete, 500);
            } catch {
                setTimeout(checkComplete, 500);
            }
        };
        checkComplete();
    });
}

async function handleIncomingMessage(message) {
    try {
        const { url, response_uuid } = JSON.parse(message);
        if (!url || !response_uuid) throw new Error('Invalid JSON structure');

        const pageContent = await capturePageContent(url);
        await sendCapturedData(response_uuid, pageContent);
    } catch (error) {
        console.error('Processing error:', error);
    }
}

async function sendCapturedData(response_uuid, pageContent, retries = 3, delay = 2000) {
    try {
        const responseUrl = `${baseURL}/response/${response_uuid}`;
        const response = await fetch(responseUrl, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                url: pageContent.url,
                title: pageContent.title,
                html: pageContent.html
            })
        });

        if (!response.ok) throw new Error(`Failed to send data: ${response.statusText}`);
        console.log(`Data sent successfully to ${responseUrl}`);
    } catch (error) {
        console.error('HTTP POST error:', error);
        if (retries > 0) {
            setTimeout(() => sendCapturedData(response_uuid, pageContent, retries - 1, delay * 2), delay);
        }
    }
}

let eventSource;
let retryTimeout = 100;
async function startSSEListener() {
    if (eventSource) eventSource.close();
    const uuid = await getUUID();
    const eventsUrl = `${baseURL}/events/${uuid}`;

    eventSource = new EventSource(eventsUrl);
    eventSource.onmessage = (event) => {
        retryTimeout = 100;
        handleIncomingMessage(event.data);
    }

    eventSource.onerror = () => {
        console.error('SSE connection lost. Retrying...');
        eventSource.close();
        setTimeout(startSSEListener, retryTimeout);
        retryTimeout = Math.min(retryTimeout * 2, 30000); // Exponential backoff, max 30s
    };
}

chrome.runtime.onInstalled.addListener(async () => {
    await getUUID();
    startSSEListener();
});

chrome.runtime.onStartup.addListener(() => {
    startSSEListener();
});

chrome.runtime.onMessage.addListener(async (request, sender, sendResponse) => {
    if (request.action === 'resetUUID') {
        try {
            await generateUUID();
            sendResponse({ success: true });
            startSSEListener();
        } catch (error) {
            sendResponse({ success: false, error: error.message });
        }
    }
    return true;
});
