import React, { useState } from 'react';
import Home from './components/Home';
import WhoisTab from './components/WhoisTab';
import IPQualityTab from './components/IPQualityTab';
import './App.css';

function App() {
    const [activeTab, setActiveTab] = useState('home');
    const [whoisResults, setWhoisResults] = useState([]);
    const [ipQualityResults, setIpQualityResults] = useState([]);

    return (
        <div className="main-container">
            <nav>
                <div 
                    className={`nav-item ${activeTab === 'home' ? 'active' : ''}`}
                    onClick={() => setActiveTab('home')}
                >
                    Home
                </div>
                <div 
                    className={`nav-item ${activeTab === 'whois' ? 'active' : ''}`}
                    onClick={() => setActiveTab('whois')}
                >
                    Whois Results
                </div>
                <div 
                    className={`nav-item ${activeTab === 'quality' ? 'active' : ''}`}
                    onClick={() => setActiveTab('quality')}
                >
                    IPQuality Results
                </div>
            </nav>

            <div className="content">
                {activeTab === 'home' && <Home setWhois={setWhoisResults} setQuality={setIpQualityResults} switchTab={setActiveTab} />}
                {activeTab === 'whois' && <WhoisTab results={whoisResults} />}
                {activeTab === 'quality' && <IPQualityTab results={ipQualityResults} />}
            </div>
        </div>
    );
}

export default App;
