import React, { useState } from 'react';
import * as api from '../services/api';

function Home({ runWhois, runQuality, loading }) {
    const [input, setInput] = useState('');
    const [stats, setStats] = useState({ total: 0, ips: 0, proxies: 0 });

    const handleParse = async () => {
        if (!input.trim()) return;
        try {
            const res = await api.parseInput(input);
            setStats({
                total: res.total || 0,
                ips: res.ips?.length ?? 0,
                proxies: res.proxies?.length ?? 0
            });
        } catch (err) {
            console.error(err);
        }
    };

    return (
        <div className="home-container">
            <div className="page-header">
                <h1 className="page-title">IP & Proxy Checker</h1>
                <p style={{ color: 'var(--text-muted)' }}>Parse and check IP information from multiple sources</p>
            </div>

            <div style={{ background: '#fff', padding: '1.5rem', borderRadius: '0.5rem', boxShadow: '0 0.125rem 0.25rem rgba(165, 163, 174, 0.3)' }}>
                <textarea
                    placeholder="Enter IPs (one per line) or Proxies (IP:Port:User:Pass)..."
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onBlur={handleParse}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                            // Ensure Enter (and Ctrl+Enter) always creates a newline
                            // Default behavior for textarea is newline, but we stop propagation
                            // to avoid triggering any parent form actions or general key handlers.
                            e.stopPropagation();
                        }
                    }}
                    style={{ resize: 'vertical' }}
                />
                <div style={{ display: 'flex', gap: '15px', alignItems: 'center', flexWrap: 'wrap' }}>
                    <button className="btn-primary" onClick={() => runWhois(input)} disabled={loading}>
                        {loading ? 'Processing...' : 'Check with IPWho'}
                    </button>
                    <button className="btn-secondary" onClick={() => runQuality(input)} disabled={loading}>
                        {loading ? 'Processing...' : 'Check with IPQuality'}
                    </button>

                    <div style={{ marginLeft: 'auto', display: 'flex', gap: '20px', fontSize: '0.875rem' }}>
                        <div><span style={{ color: 'var(--text-muted)' }}>Total:</span> <span style={{ fontWeight: 600 }}>{stats.total}</span></div>
                        <div><span style={{ color: 'var(--text-muted)' }}>IPs:</span> <span style={{ fontWeight: 600 }}>{stats.ips}</span></div>
                        <div><span style={{ color: 'var(--text-muted)' }}>Proxies:</span> <span style={{ fontWeight: 600 }}>{stats.proxies}</span></div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Home;
