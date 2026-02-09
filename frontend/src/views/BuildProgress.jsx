import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { Check, Loader2 } from 'lucide-react';
import Header from '../components/UI/Header';
import Button from '../components/UI/Button';
import Card from '../components/UI/Card';

const steps = [
    { id: 1, text: '12 files changed', delay: 500 },
    { id: 2, text: 'Tests passed', delay: 1500 },
    { id: 3, text: 'Preview ready', delay: 2500 },
];

const BuildProgress = ({ project, onBack, onApprove }) => {
    const [logs, setLogs] = useState([]);
    const [status, setStatus] = useState('Connecting...');
    const [aiAnalysis, setAiAnalysis] = useState(null);
    const [analyzing, setAnalyzing] = useState(false);

    const handleAskAI = async () => {
        setAnalyzing(true);
        setAiAnalysis(null);
        try {
            const res = await fetch('http://localhost:8080/api/ai/analyze', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ logs: logs.join('\n') })
            });
            const data = await res.json();
            setAiAnalysis(data);
        } catch (err) {
            console.error(err);
        } finally {
            setAnalyzing(false);
        }
    };

    useEffect(() => {
        if (!project) return;

        // Reset logs when project changes or connection restarts
        setLogs(['>> Initializing connection...']);
        setStatus('Connecting...');

        const ws = new WebSocket('ws://localhost:8080/ws');
        let isMounted = true;

        ws.onopen = () => {
            if (isMounted) {
                setStatus('Connected');
                setLogs(prev => [...prev, '>> Connected to Real-time Gateway.', '>> Waiting for Agent...']);
                ws.send(JSON.stringify({
                    type: 'IDENTIFY',
                    payload: {
                        project_id: project.id,
                        role: 'CLIENT'
                    }
                }));

                // Trigger the actual build now that we are listening
                fetch(`http://localhost:8080/api/projects/${project.id}/build`, { method: 'POST' })
                    .then(res => res.json())
                    .then(data => {
                        if (isMounted) {
                            setStatus('Build Started');
                            setLogs(prev => [...prev, `>> Build Job Created: ${data.id}`]);
                        }
                    })
                    .catch(err => {
                        if (isMounted) {
                            setStatus('Error');
                            setLogs(prev => [...prev, `>> Failed to trigger build: ${err.message}`]);
                        }
                    });
            }
        };

        ws.onmessage = (event) => {
            if (!isMounted) return;
            try {
                const msg = JSON.parse(event.data);
                if (msg.type === 'LOG_CHUNK') {
                    const rawText = msg.payload.chunk || '';
                    // Robust ANSI strip regex
                    const cleanText = rawText.replace(/[\u001b\u009b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]/g, '');
                    setLogs(prev => [...prev, cleanText]);
                } else if (msg.type === 'JOB_UPDATE') {
                    setStatus(msg.payload.status);
                    setLogs(prev => [...prev, `\n>> Job Status: ${msg.payload.status}`]);
                }
            } catch (e) {
                console.error("WS Parse Error", e);
            }
        };

        ws.onerror = (e) => {
            if (isMounted) setLogs(prev => [...prev, '>> Connection Error']);
        };

        ws.onclose = () => {
            if (isMounted) setLogs(prev => [...prev, '>> Disconnected']);
        };

        return () => {
            isMounted = false;
            ws.close();
        };
    }, [project]);

    return (
        <>
            <Header title={`Build: ${project?.name || 'Unknown'}`} onBack={onBack} />

            <div style={{ padding: '0 20px', display: 'flex', flexDirection: 'column', height: '100%' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                    <h2 style={{ fontSize: '1.2rem', fontWeight: 600 }}>Live Logs</h2>
                    <span style={{ fontSize: '0.9rem', color: 'var(--text-secondary)' }}>Status: {status}</span>
                </div>

                <Card style={{ flex: 1, marginBottom: '20px', overflow: 'hidden', display: 'flex', flexDirection: 'column', background: '#1e1e1e', padding: 0 }}>
                    <div style={{
                        padding: '16px',
                        fontFamily: 'monospace',
                        fontSize: '0.9rem',
                        color: '#ddd',
                        overflowY: 'auto',
                        height: '100%',
                        whiteSpace: 'pre-wrap'
                    }}>
                        {logs.map((log, i) => (
                            <span key={i}>{log}</span>
                        ))}
                    </div>
                </Card>

                {/* AI Analysis Result */}
                {aiAnalysis && (
                    <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        style={{ marginBottom: '20px' }}
                    >
                        <Card style={{ background: 'linear-gradient(135deg, #2a2a3e 0%, #1e1e2e 100%)', border: '1px solid #4a4a5e' }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
                                <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#a78bfa' }} />
                                <h3 style={{ fontSize: '1rem', fontWeight: 600, color: '#e0e0e0' }}>AI Insight</h3>
                                <span style={{ marginLeft: 'auto', fontSize: '0.8rem', color: '#a78bfa' }}>{aiAnalysis.confidence}% Conf.</span>
                            </div>
                            <p style={{ fontSize: '0.9rem', color: '#ccc', marginBottom: '12px' }}>{aiAnalysis.analysis}</p>
                            <div style={{ background: 'rgba(0,0,0,0.3)', padding: '10px', borderRadius: '6px' }}>
                                <strong style={{ display: 'block', fontSize: '0.8rem', color: '#fff', marginBottom: '4px' }}>Suggestion:</strong>
                                <code style={{ fontSize: '0.85rem', color: '#a78bfa', fontFamily: 'monospace' }}>{aiAnalysis.suggestion}</code>
                            </div>
                        </Card>
                    </motion.div>
                )}

                <div style={{ paddingBottom: '40px', display: 'flex', flexDirection: 'column', gap: '12px' }}>

                    {!aiAnalysis && (
                        <Button fullWidth onClick={handleAskAI} disabled={analyzing} style={{ background: 'var(--bg-secondary)', border: '1px solid var(--border-color)' }}>
                            {analyzing ? (
                                <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                    <Loader2 className="spin" size={16} /> Analyzing...
                                </div>
                            ) : (
                                "Analyze with AI"
                            )}
                        </Button>
                    )}

                    <Button fullWidth onClick={onApprove}>
                        Return to Dashboard
                    </Button>
                </div>
            </div>
        </>
    );
};

export default BuildProgress;
