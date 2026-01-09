const API_BASE = '/api';

export const parseInput = async (input) => {
    const response = await fetch(`${API_BASE}/parse`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ input })
    });
    return response.json();
};

export const checkWhois = async (ips) => {
    const response = await fetch(`${API_BASE}/check/whois`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ips })
    });
    return response.json();
};

export const checkIPQuality = async (proxies, apiKey) => {
    const response = await fetch(`${API_BASE}/check/quality`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ proxies, api_key: apiKey })
    });
    return response.json();
};
