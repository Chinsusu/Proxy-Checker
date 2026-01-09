import React, { useState } from 'react';
import * as GreetService from '../../wailsjs/go/main/App';

function Home({ setWhois, setQuality, switchTab }) {
    const [input, setInput] = useState('');
    const [stats, setStats] = useState({ total: 0, ips: 0, proxies: 0 });

    const handleParse = async () => {
        const res = await GreetService.ParseInput(input);
        setStats({
            total: res.total,
            ips: res.ips.length,
            proxies: res.proxies.length
        });
    };

    const handleCheckWhois = async () => {
        const res = await GreetService.ParseInput(input);
        const ips = res.ips.map(i => i.ip);
        switchTab('whois');
        const results = await GreetService.CheckWhois(ips);
        setWhois(results);
    };

    const handleCheckQuality = async () => {
        const res = await GreetService.ParseInput(input);
        const proxies = res.proxies.map(p => `${p.ip}:${p.port}`);
        switchTab('quality');
        const results = await GreetService.CheckIPQuality(proxies);
        setQuality(results);
    };

    return (
        <div className="fade-in">
            <h1 style={{ marginBottom: '20px', color: 'var(--accent-blue)' }}>Proxy Information Checker</h1>
            <div className="glass-card">
                <textarea
                    placeholder="Enter IPs (one per line) or Proxies (IP:Port:User:Pass)..."
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onBlur={handleParse}
                />
                <div style={{ display: 'flex', gap: '15px', alignItems: 'center' }}>
                    <button className="btn-primary" onClick={handleCheckWhois}>Check with IPWho</button>
                    <button className="btn-secondary" onClick={handleCheckQuality}>Check with IPQuality</button>
                    <div style={{ marginLeft: 'auto', display: 'flex', gap: '20px', color: 'var(--text-dim)' }}>
                        <span>Total: {stats.total}</span>
                        <span>IPs: {stats.ips}</span>
                        <span>Proxies: {stats.proxies}</span>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Home;
