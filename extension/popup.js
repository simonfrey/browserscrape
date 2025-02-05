document.addEventListener('DOMContentLoaded', async () => {
    document.getElementById('resetUUID').addEventListener('click', resetUUID);
    await displayUUID();
});

async function resetUUID() {
    try {
        await chrome.runtime.sendMessage({ action: 'resetUUID' });
        await displayUUID();
    } catch (error) {
        console.error('Error resetting UUID:', error);
    }
}

async function displayUUID() {
    try {
        if (chrome.storage && chrome.storage.local) {
            const { uuid } = await chrome.storage.local.get('uuid');
            document.querySelectorAll('.uuid').forEach((element) => {
              element.textContent = uuid || 'UUID not found';
            })
        } else {
            console.warn('chrome.storage.local is not available.');
            document.getElementById('uuid').textContent = 'Storage API unavailable';
        }
    } catch (error) {
        console.error('Error accessing storage:', error);
        document.getElementById('uuid').textContent = 'Error loading UUID';
    }
}
