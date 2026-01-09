import React from 'react';

function IPQualityTab({ results }) {
    const exportToCSV = () => {
        if (results.length === 0) return;
        const headers = ["IP:Port", "Status", "Country", "City", "VPN", "Proxy", "ISP", "Organization"];
        const rows = results.map(r => [
            `${r.ip}:${r.port}`, r.status, r.country, r.city, r.vpn, r.proxy, r.isp, r.organization
        ]);

        const csvContent = [headers, ...rows].map(e => e.join(",")).join("\n");
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement("a");
        const url = URL.createObjectURL(blob);
        link.setAttribute("href", url);
        link.setAttribute("download", "ipquality_results.csv");
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <div className="fade-in">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
                <h2>IPQuality Results</h2>
                <button className="btn-secondary" onClick={exportToCSV} disabled={results.length === 0}>
                    Export CSV
                </button>
            </div>
            <div className="glass-card" style={{ padding: '0' }}>
                <table>
                    <thead>
                        <tr>
                            <th>IP:Port</th>
                            <th>Status</th>
                            <th>Country</th>
                            <th>City</th>
                            <th>VPN</th>
                            <th>Proxy</th>
                            <th>ISP</th>
                            <th>Organization</th>
                        </tr>
                    </thead>
                    <tbody>
                        {results.map((res, i) => (
                            <tr key={i}>
                                <td>{res.ip}:{res.port}</td>
                                <td className={res.status === 'Live' ? 'status-live' : 'status-dead'}>{res.status}</td>
                                <td>{res.country}</td>
                                <td>{res.city}</td>
                                <td>{res.vpn ? '✅' : '❌'}</td>
                                <td>{res.proxy ? '✅' : '❌'}</td>
                                <td>{res.isp}</td>
                                <td>{res.organization}</td>
                            </tr>
                        ))}
                        {results.length === 0 && (
                            <tr>
                                <td colSpan="8" style={{ textAlign: 'center', padding: '40px', color: 'var(--text-dim)' }}>
                                    No results yet. Start checking from the Home tab.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}

export default IPQualityTab;
