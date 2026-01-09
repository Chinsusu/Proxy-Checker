import React, { useState } from 'react';
import Home from './components/Home';
import WhoisResults from './components/WhoisResults';
import IPQualityResults from './components/IPQualityResults';
import * as api from './services/api';
import './index.css';

function App() {
    const [activeTab, setActiveTab] = useState('home');
    const [whoisResults, setWhoisResults] = useState([]);
    const [ipQualityResults, setIpQualityResults] = useState([]);
    const [loading, setLoading] = useState(false);

    const runWhoisCheck = async (input) => {
        setLoading(true);
        try {
            const parseRes = await api.parseInput(input);
            const ips = parseRes.ips.map(i => i.ip);
            if (ips.length === 0) {
                alert("No IPs found in input!");
                setLoading(false);
                return;
            }
            setActiveTab('whois');
            const results = await api.checkWhois(ips);
            console.log("Whois Results:", results);
            setWhoisResults(results);
        } catch (err) {
            console.error(err);
            alert("Error checking Whois: " + err.message);
        } finally {
            setLoading(false);
        }
    };

    const runQualityCheck = async (input) => {
        setLoading(true);
        try {
            const parseRes = await api.parseInput(input);
            const proxies = parseRes.proxies.map(p => {
                let s = `${p.host}:${p.port}`;
                if (p.username && p.password) {
                    s += `:${p.username}:${p.password}`;
                }
                return s;
            });
            if (proxies.length === 0) {
                alert("No proxies found in input!");
                setLoading(false);
                return;
            }
            setActiveTab('quality');
            const apiKey = localStorage.getItem('ipquality_api_key') || '';
            console.log("Preparing Quality check with API Key:", apiKey ? "PRESENT" : "MISSING");
            const results = await api.checkIPQuality(proxies, apiKey);
            console.log("Quality Results:", results);

            // Map results to match UI expectations
            const mappedResults = results.map(r => ({
                ...r,
                proxy: `${r.ip}:${r.port}`,
                vpn: r.vpn ? 'Yes' : 'No',
                proxyFlag: r.proxy ? 'Yes' : 'No',
                fraudScore: r.fraud_score || 'N/A'
            }));

            setIpQualityResults(mappedResults);
        } catch (err) {
            console.error(err);
            alert("Error checking Quality: " + err.message);
        } finally {
            setLoading(false);
        }
    };

    const exportWhoisCSV = () => {
        if (!whoisResults || whoisResults.length === 0) return;
        const headers = ["IP", "Country", "Region", "City", "ISP", "ASN", "Timezone", "Status"];
        const rows = whoisResults.map(r => [
            r.ip, r.country, r.region, r.city, r.isp, r.asn, r.timezone, r.status
        ]);
        const csvContent = [headers, ...rows].map(e => e.join(",")).join("\n");
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement("a");
        link.href = URL.createObjectURL(blob);
        link.download = "whois_results.csv";
        link.click();
    };

    const exportQualityCSV = () => {
        if (!ipQualityResults || ipQualityResults.length === 0) return;
        const headers = ["IP:Port", "Status", "Country", "City", "VPN", "Proxy", "ISP", "Organization"];
        const rows = ipQualityResults.map(r => [
            r.proxy, r.status, r.country, r.city, r.vpn, r.proxyFlag, r.isp, r.organization
        ]);
        const csvContent = [headers, ...rows].map(e => e.join(",")).join("\n");
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement("a");
        link.href = URL.createObjectURL(blob);
        link.download = "quality_results.csv";
        link.click();
    };

    return (
        <div className="app-wrapper">
            {/* Sidebar */}
            <aside className="sidebar">
                <div className="sidebar-header">
                    <div className="brand-text">IP Checker</div>
                </div>
                <ul className="menu-list">
                    <li
                        className={`menu-item ${activeTab === 'home' ? 'active' : ''}`}
                        onClick={() => setActiveTab('home')}
                    >
                        <span className="menu-icon">🏠</span>
                        <span>Home</span>
                    </li>
                    <li
                        className={`menu-item ${activeTab === 'whois' ? 'active' : ''}`}
                        onClick={() => setActiveTab('whois')}
                    >
                        <span className="menu-icon">🌍</span>
                        <span>Whois Results</span>
                    </li>
                    <li
                        className={`menu-item ${activeTab === 'quality' ? 'active' : ''}`}
                        onClick={() => setActiveTab('quality')}
                    >
                        <span className="menu-icon">🛡️</span>
                        <span>IP Quality</span>
                    </li>
                </ul>

                {loading && (
                    <div style={{ padding: '1rem', textAlign: 'center' }}>
                        <div className="badge badge-info" style={{ width: '100%' }}>Checking...</div>
                    </div>
                )}
            </aside>

            {/* Main Content */}
            <main className="main-content">
                <div className="fade-in">
                    {activeTab === 'home' && (
                        <Home
                            runWhois={runWhoisCheck}
                            runQuality={runQualityCheck}
                            loading={loading}
                        />
                    )}
                    {activeTab === 'whois' && (
                        <WhoisResults
                            results={whoisResults}
                            onExportCSV={exportWhoisCSV}
                        />
                    )}
                    {activeTab === 'quality' && (
                        <IPQualityResults
                            results={ipQualityResults}
                            onExportCSV={exportQualityCSV}
                        />
                    )}
                </div>
            </main>
        </div>
    );
}

export default App;
