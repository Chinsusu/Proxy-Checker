import React from 'react';
import './WhoisTab.css';

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
        <div className="whois-container">
            <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                    <h1 className="page-title">Whois Results</h1>
                    <p style={{ color: 'var(--text-muted)' }}>Detailed geographical and ISP information</p>
                </div>
                <button className="btn-secondary" onClick={exportToCSV} disabled={results.length === 0}>
                    Export CSV
                </button>
            </div>

            <div className="table-container fade-in">
                <table className="results-table">
                    <thead>
                        <tr>
                            <th className="col-ip">IP</th>
                            <th className="col-country">Country</th>
                            <th className="col-region">Region</th>
                            <th className="col-city">City</th>
                            <th className="col-isp">ISP</th>
                            <th className="col-asn">ASN</th>
                            <th className="col-timezone">Timezone</th>
                            <th className="col-status">Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        {results.map((res, i) => (
                            <tr key={i}>
                                <td className="col-ip" style={{ fontWeight: 600, color: '#000' }}>{res.ip}</td>
                                <td className="col-country">
                                    <div style={{ display: 'flex', alignItems: 'center' }}>
                                        {res.flag && <img src={res.flag} alt="flag" style={{ width: '18px', marginRight: '8px', borderRadius: '2px' }} />}
                                        {res.country}
                                    </div>
                                </td>
                                <td className="col-region">{res.region}</td>
                                <td className="col-city">{res.city}</td>
                                <td className="col-isp" title={res.isp}>{res.isp}</td>
                                <td className="col-asn">
                                    <span className="badge" style={{ backgroundColor: '#f2f2f3', color: '#5d596c' }} title={res.asn}>
                                        {res.asn}
                                    </span>
                                </td>
                                <td className="col-timezone">{res.timezone}</td>
                                <td className="col-status">
                                    <span
                                        className={`status-badge ${res.status === 'success' ? 'status-success' : 'status-failed'}`}
                                        title={res.error || ''}
                                    >
                                        {res.status}
                                    </span>
                                </td>
                            </tr>
                        ))}
                        {results.length === 0 && (
                            <tr>
                                <td colSpan="8" style={{ textAlign: 'center', padding: '100px', color: 'var(--text-muted)' }}>
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
