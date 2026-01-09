import React from 'react';
import './IPQualityTab.css';

function IPQualityTab({ results }) {
    const total = results.length;
    const live = results.filter(r => r.status === 'Live').length;
    const dead = total - live;
    const liveRate = total > 0 ? ((live / total) * 100).toFixed(1) : 0;

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
        <div className="quality-container">
            <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                    <h1 className="page-title">IP Quality Analysis</h1>
                    <p style={{ color: 'var(--text-muted)' }}>Proxy and VPN detection results</p>
                </div>
                <button className="btn-secondary" onClick={exportToCSV} disabled={results.length === 0}>
                    Export CSV
                </button>
            </div>

            {total > 0 && (
                <div className="stats-bar fade-in">
                    <div className="stat-item">
                        <span className="stat-label">Total Proxies</span>
                        <span className="stat-value">{total}</span>
                    </div>
                    <div className="stat-item">
                        <span className="stat-label">Live</span>
                        <span className="stat-value" style={{ color: 'var(--success)' }}>{live}</span>
                    </div>
                    <div className="stat-item">
                        <span className="stat-label">Dead</span>
                        <span className="stat-value" style={{ color: 'var(--danger)' }}>{dead}</span>
                    </div>
                    <div className="stat-item">
                        <span className="stat-label">Live Rate</span>
                        <span className="stat-value" style={{ color: 'var(--primary)' }}>{liveRate}%</span>
                    </div>
                </div>
            )}

            <div className="table-container fade-in">
                <table className="results-table">
                    <thead>
                        <tr>
                            <th className="col-proxy">IP:Port</th>
                            <th className="col-status">Status</th>
                            <th className="col-country">Country</th>
                            <th className="col-city">City</th>
                            <th className="col-vpn">VPN</th>
                            <th className="col-proxy-flag">Proxy</th>
                            <th className="col-isp">ISP</th>
                            <th className="col-org">Organization</th>
                        </tr>
                    </thead>
                    <tbody>
                        {results.map((res, i) => {
                            let rowClass = "row-clean";
                            if (res.status !== 'Live') rowClass = "row-dead";
                            else if (res.vpn || res.proxy) rowClass = "row-risk";

                            return (
                                <tr key={i} className={rowClass}>
                                    <td className="col-proxy" style={{ fontWeight: 600, color: '#000' }}>{res.ip}:{res.port}</td>
                                    <td className="col-status">
                                        <span
                                            className={`status-badge ${res.status === 'Live' ? 'status-success' : 'status-failed'}`}
                                            title={res.error || ''}
                                        >
                                            {res.status}
                                        </span>
                                    </td>
                                    <td className="col-country">{res.country}</td>
                                    <td className="col-city">{res.city}</td>
                                    <td className="col-vpn">
                                        {res.vpn ? <span className="badge badge-vpn">VPN</span> : '-'}
                                    </td>
                                    <td className="col-proxy-flag">
                                        {res.proxy ? <span className="badge badge-danger">PROXY</span> : '-'}
                                    </td>
                                    <td className="col-isp" title={res.isp}>{res.isp}</td>
                                    <td className="col-org" title={res.organization}>{res.organization}</td>
                                </tr>
                            );
                        })}
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

export default IPQualityTab;
