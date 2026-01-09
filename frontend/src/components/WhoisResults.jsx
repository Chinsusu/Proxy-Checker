import React, { useState } from 'react';
import './WhoisResults.css';

const WhoisResults = ({ results = [], onExportCSV }) => {
    const [sortConfig, setSortConfig] = useState({ key: null, direction: 'asc' });
    const [searchTerm, setSearchTerm] = useState('');

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

    const filteredResults = sortedResults.filter(result =>
        Object.values(result).some(value =>
            String(value).toLowerCase().includes(searchTerm.toLowerCase())
        )
    );

    return (
        <div className="whois-results">
            <div className="results-header">
                <h2>Whois Results</h2>
                <p className="results-subtitle">Detailed geographical and ISP information</p>
            </div>

            <div className="results-actions">
                <input
                    type="text"
                    className="search-input"
                    placeholder="Search results..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                />
                <button className="btn-export" onClick={onExportCSV}>
                    Export CSV
                </button>
            </div>

            <div className="table-container">
                <table className="results-table">
                    <thead>
                        <tr>
                            <th className="col-ip" onClick={() => handleSort('ip')}>
                                IP
                                {sortConfig.key === 'ip' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-country" onClick={() => handleSort('country')}>
                                COUNTRY
                                {sortConfig.key === 'country' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-region" onClick={() => handleSort('region')}>
                                REGION
                                {sortConfig.key === 'region' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-city" onClick={() => handleSort('city')}>
                                CITY
                                {sortConfig.key === 'city' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-isp" onClick={() => handleSort('isp')}>
                                ISP
                                {sortConfig.key === 'isp' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-asn" onClick={() => handleSort('asn')}>
                                ASN
                                {sortConfig.key === 'asn' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-timezone" onClick={() => handleSort('timezone')}>
                                TIMEZONE
                                {sortConfig.key === 'timezone' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                            <th className="col-status" onClick={() => handleSort('status')}>
                                STATUS
                                {sortConfig.key === 'status' && (
                                    <span className="sort-icon">{sortConfig.direction === 'asc' ? '↑' : '↓'}</span>
                                )}
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {filteredResults.length > 0 ? (
                            filteredResults.map((result, index) => (
                                <tr key={index} className={index % 2 === 0 ? 'row-even' : 'row-odd'}>
                                    <td className="col-ip">{result.ip}</td>
                                    <td className="col-country">{result.country}</td>
                                    <td className="col-region">{result.region}</td>
                                    <td className="col-city">{result.city}</td>
                                    <td className="col-isp">{result.isp}</td>
                                    <td className="col-asn">{result.asn}</td>
                                    <td className="col-timezone">{result.timezone}</td>
                                    <td className="col-status">
                                        <span className={`status-badge status-${result.status?.toLowerCase()}`}>
                                            {result.status}
                                        </span>
                                    </td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan="8" className="no-results">
                                    {(!results || results.length === 0) ? 'No results yet. Check some IPs from the Home tab.' : 'No results match your search.'}
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>

            <div className="results-footer">
                <p>Total Results: {filteredResults.length} of {results?.length ?? 0}</p>
            </div>
        </div>
    );
};

export default WhoisResults;
