import React, { useState } from 'react';
import './IPQualityResults.css';

const IPQualityResults = ({ results = [], onExportCSV }) => {
    const [sortConfig, setSortConfig] = useState({ key: null, direction: 'asc' });
    const [searchTerm, setSearchTerm] = useState('');
    const [filterStatus, setFilterStatus] = useState('all');

    const handleSort = (key) => {
        let direction = 'asc';
        if (sortConfig.key === key && sortConfig.direction === 'asc') {
            direction = 'desc';
        }
        setSortConfig({ key, direction });
    };

    const sortedResults = React.useMemo(() => {
        let sortableResults = [...results];
        if (sortConfig.key) {
            sortableResults.sort((a, b) => {
                if (a[sortConfig.key] < b[sortConfig.key]) {
                    return sortConfig.direction === 'asc' ? -1 : 1;
                }
                if (a[sortConfig.key] > b[sortConfig.key]) {
                    return sortConfig.direction === 'asc' ? 1 : -1;
                }
                return 0;
            });
        }
        return sortableResults;
    }, [results, sortConfig]);

    const filteredResults = sortedResults.filter(result => {
        const matchesSearch = Object.values(result).some(value =>
            String(value).toLowerCase().includes(searchTerm.toLowerCase())
        );

        const matchesStatus = filterStatus === 'all' ||
            (filterStatus === 'live' && result.status?.toLowerCase() === 'live') ||
            (filterStatus === 'dead' && result.status?.toLowerCase() === 'dead');

        return matchesSearch && matchesStatus;
    });

    const liveCount = results.filter(r => r.status?.toLowerCase() === 'live').length;
    const deadCount = results.filter(r => r.status?.toLowerCase() === 'dead').length;
    const liveRate = results.length > 0 ? ((liveCount / results.length) * 100).toFixed(1) : 0;

    return (
        <div className="ipquality-results">
            <div className="results-header">
                <h2>IP Quality Results</h2>
                <p className="results-subtitle">Proxy validation and VPN detection</p>
            </div>

            <div className="stats-bar">
                <div className="stat-item">
                    <span className="stat-label">Total</span>
                    <span className="stat-value">{results.length}</span>
                </div>
                <div className="stat-item stat-live">
                    <span className="stat-label">Live</span>
                    <span className="stat-value">{liveCount}</span>
                </div>
                <div className="stat-item stat-dead">
                    <span className="stat-label">Dead</span>
                    <span className="stat-value">{deadCount}</span>
                </div>
                <div className="stat-item">
                    <span className="stat-label">Live Rate</span>
                    <span className="stat-value">{liveRate}%</span>
                </div>
            </div>

            <div className="results-actions">
                <input
                    type="text"
                    className="search-input"
                    placeholder="Search results..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                />

                <div className="filter-buttons">
                    <button
                        className={`filter-btn ${filterStatus === 'all' ? 'active' : ''}`}
                        onClick={() => setFilterStatus('all')}
                    >
                        All
                    </button>
                    <button
                        className={`filter-btn ${filterStatus === 'live' ? 'active' : ''}`}
                        onClick={() => setFilterStatus('live')}
                    >
                        Live
                    </button>
                    <button
                        className={`filter-btn ${filterStatus === 'dead' ? 'active' : ''}`}
                        onClick={() => setFilterStatus('dead')}
                    >
                        Dead
                    </button>
                </div>

                <button className="btn-export" onClick={onExportCSV}>
                    Export CSV
                </button>
            </div>

            <div className="table-container">
                <table className="results-table">
                    <thead>
                        <tr>
                            <th className="col-proxy" onClick={() => handleSort('proxy')}>
                                PROXY (IP:PORT)
                                {sortConfig.key === 'proxy' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-status" onClick={() => handleSort('status')}>
                                STATUS
                                {sortConfig.key === 'status' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-country" onClick={() => handleSort('country')}>
                                COUNTRY
                                {sortConfig.key === 'country' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-city" onClick={() => handleSort('city')}>
                                CITY
                                {sortConfig.key === 'city' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-region" onClick={() => handleSort('region')}>
                                REGION
                                {sortConfig.key === 'region' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-vpn" onClick={() => handleSort('vpn')}>
                                VPN
                                {sortConfig.key === 'vpn' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-proxy-flag" onClick={() => handleSort('proxyFlag')}>
                                PROXY
                                {sortConfig.key === 'proxyFlag' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-isp" onClick={() => handleSort('isp')}>
                                ISP
                                {sortConfig.key === 'isp' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-org" onClick={() => handleSort('organization')}>
                                ORGANIZATION
                                {sortConfig.key === 'organization' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {filteredResults.length > 0 ? (
                            filteredResults.map((result, index) => {
                                const isLive = result.status?.toLowerCase() === 'live';
                                const isDead = result.status?.toLowerCase() === 'dead';
                                const hasVPN = result.vpn?.toLowerCase() === 'yes';
                                const hasProxy = result.proxyFlag?.toLowerCase() === 'yes';

                                let rowClass = index % 2 === 0 ? 'row-even' : 'row-odd';
                                if (isLive && !hasVPN && !hasProxy) rowClass += ' row-clean';
                                if (isLive && (hasVPN || hasProxy)) rowClass += ' row-warning';
                                if (isDead) rowClass += ' row-dead';

                                return (
                                    <tr key={index} className={rowClass}>
                                        <td className="col-proxy">{result.proxy}</td>
                                        <td className="col-status">
                                            <span className={`status-badge status-${result.status?.toLowerCase()}`}>
                                                {result.status}
                                            </span>
                                        </td>
                                        <td className="col-country">{result.country}</td>
                                        <td className="col-city">{result.city}</td>
                                        <td className="col-region">{result.region}</td>
                                        <td className="col-vpn">
                                            <span className={`flag-badge ${hasVPN ? 'flag-yes' : 'flag-no'}`}>
                                                {result.vpn}
                                            </span>
                                        </td>
                                        <td className="col-proxy-flag">
                                            <span className={`flag-badge ${hasProxy ? 'flag-yes' : 'flag-no'}`}>
                                                {result.proxyFlag}
                                            </span>
                                        </td>
                                        <td className="col-isp">{result.isp}</td>
                                        <td className="col-org">{result.organization}</td>
                                    </tr>
                                );
                            })
                        ) : (
                            <tr>
                                <td colSpan="9" className="no-results">
                                    {results.length === 0 ? 'No results yet. Check some proxies from the Home tab.' : 'No results match your filters.'}
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>

            <div className="results-footer">
                <p>Showing {filteredResults.length} of {results.length} results</p>
            </div>
        </div>
    );
};

export default IPQualityResults;
