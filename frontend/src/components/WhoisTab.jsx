import React from 'react';

function WhoisTab({ results }) {
    const exportToCSV = () => {
        if (results.length === 0) return;
        const headers = ["IP", "Country", "Region", "City", "ISP", "ASN", "Timezone", "Status"];
        const rows = results.map(r => [
            r.ip, r.country, r.region, r.city, r.isp, r.asn, r.timezone, r.status
        ]);

        const csvContent = [headers, ...rows].map(e => e.join(",")).join("\n");
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement("a");
        const url = URL.createObjectURL(blob);
        link.setAttribute("href", url);
        link.setAttribute("download", "whois_results.csv");
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <div className="fade-in">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
                <h2>Whois Results</h2>
                <button className="btn-secondary" onClick={exportToCSV} disabled={results.length === 0}>
                    Export CSV
                </button>
            </div>
            <div className="glass-card" style={{ padding: '0' }}>
                <table>
                    <thead>
                        <tr>
                            <th>IP</th>
                            <th>Country</th>
                            <th>Region</th>
                            <th>City</th>
                            <th>ISP</th>
                            <th>ASN</th>
                            <th>Timezone</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        {results.map((res, i) => (
                            <tr key={i}>
                                <td>{res.ip}</td>
                                <td>
                                    {res.flag && <img src={res.flag} alt="flag" style={{ width: '20px', marginRight: '8px', verticalAlign: 'middle' }} />}
                                    {res.country}
                                </td>
                                <td>{res.region}</td>
                                <td>{res.city}</td>
                                <td>{res.isp}</td>
                                <td>{res.asn}</td>
                                <td>{res.timezone}</td>
                                <td>
                                    <span style={{ color: res.status === 'success' ? 'var(--accent-green)' : 'var(--accent-red)' }}>
                                        {res.status}
                                    </span>
                                </td>
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

export default WhoisTab;
