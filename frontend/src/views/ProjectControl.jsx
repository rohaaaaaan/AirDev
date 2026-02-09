import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Play, Pause, AlertCircle, X, Lightbulb } from 'lucide-react';
import Header from '../components/UI/Header';
import Button from '../components/UI/Button';
import Card from '../components/UI/Card';

const ProjectControl = ({ project, onBack, onStartBuild }) => {
    const [showInstructions, setShowInstructions] = useState(true);

    if (!project) return null;

    return (
        <>
            <Header title={project.name} onBack={onBack} />

            <div style={{ padding: '0 20px', display: 'flex', flexDirection: 'column', gap: '24px' }}>
                <Card>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                            <span style={{ color: 'var(--text-secondary)' }}>Status</span>
                            <span style={{ fontWeight: 600, color: 'var(--status-success)' }}>Online</span>
                        </div>
                        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                            <span style={{ color: 'var(--text-secondary)' }}>Branch</span>
                            <span style={{ fontFamily: 'monospace', background: 'var(--bg-tertiary)', padding: '2px 6px', borderRadius: '4px' }}>main</span>
                        </div>
                    </div>
                </Card>

                <div style={{ display: 'flex', gap: '12px' }}>
                    <Button fullWidth onClick={onStartBuild} className="flex-center" style={{ gap: '8px' }}>
                        <Play size={18} fill="currentColor" /> Start Build
                    </Button>
                    <Button variant="secondary" className="flex-center">
                        <Pause size={18} />
                    </Button>
                </div>

                <Card title="Remote Control">
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                        <select
                            id="action-select"
                            style={{
                                background: 'var(--bg-tertiary)',
                                color: 'var(--text-primary)',
                                border: '1px solid var(--text-muted)',
                                padding: '8px',
                                borderRadius: '4px'
                            }}
                        >
                            <option value="FIND">Launch / Focus (FIND)</option>
                            <option value="TYPE">Type Text (TYPE)</option>
                            <option value="CLICK">Press Key (CLICK)</option>
                        </select>
                        <input
                            id="target-input"
                            type="text"
                            placeholder="Target (e.g. Notepad, Calculator)"
                            style={{
                                background: 'var(--bg-tertiary)',
                                color: 'var(--text-primary)',
                                border: '1px solid var(--text-muted)',
                                padding: '8px',
                                borderRadius: '4px'
                            }}
                        />
                        <input
                            id="value-input"
                            type="text"
                            placeholder="Value (e.g. Hello World, enter)"
                            style={{
                                background: 'var(--bg-tertiary)',
                                color: 'var(--text-primary)',
                                border: '1px solid var(--text-muted)',
                                padding: '8px',
                                borderRadius: '4px'
                            }}
                        />
                        <Button fullWidth onClick={() => {
                            const action = document.getElementById('action-select').value;
                            const target = document.getElementById('target-input').value;
                            const value = document.getElementById('value-input').value;

                            fetch(`http://127.0.0.1:8080/api/projects/${project.id}/ai-command`, {
                                method: 'POST',
                                headers: { 'Content-Type': 'application/json' },
                                body: JSON.stringify({
                                    type: 'UI_ACTION',
                                    action,
                                    target,
                                    value
                                })
                            })
                                .then(res => res.json())
                                .then(data => alert(`Command Sent! Job ID: ${data.id}`))
                                .catch(err => alert(`Error: ${err.message}`));
                        }}>
                            Send Command
                        </Button>
                    </div>
                </Card>

                <div style={{ marginTop: 'auto' }}>
                    <p style={{ fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '8px' }}>Last Result</p>
                    <div style={{
                        padding: '12px',
                        background: 'rgba(34, 197, 94, 0.1)',
                        border: '1px solid rgba(34, 197, 94, 0.2)',
                        borderRadius: '8px',
                        color: 'var(--status-success)'
                    }}>
                        Success (2h ago)
                    </div>
                </div>
            </div>

            <AnimatePresence>
                {showInstructions && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        style={{
                            position: 'absolute',
                            inset: 0,
                            background: 'rgba(0,0,0,0.6)',
                            backdropFilter: 'blur(4px)',
                            zIndex: 50,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            padding: '20px'
                        }}
                    >
                        <motion.div
                            initial={{ scale: 0.9, y: 20 }}
                            animate={{ scale: 1, y: 0 }}
                            exit={{ scale: 0.9, y: 20 }}
                            style={{ width: '100%', maxWidth: '320px' }}
                        >
                            <Card style={{ background: 'var(--bg-secondary)', border: '1px solid var(--text-muted)' }}>
                                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                                    <h3 style={{ fontSize: '1.25rem', fontWeight: 700 }}>Instructions</h3>
                                </div>

                                <ul style={{
                                    display: 'flex',
                                    flexDirection: 'column',
                                    gap: '12px',
                                    paddingLeft: '0',
                                    listStyle: 'none',
                                    fontSize: '0.95rem',
                                    color: 'var(--text-primary)'
                                }}>
                                    {[
                                        "Select a project",
                                        "Start or stop a build",
                                        "Monitor progress",
                                        "Approve deploy previews",
                                        "Open preview link"
                                    ].map((item, i) => (
                                        <li key={i} style={{ display: 'flex', gap: '12px', opacity: 0.9 }}>
                                            <span style={{ color: 'var(--text-muted)', fontWeight: 600 }}>{i + 1}.</span>
                                            {item}
                                        </li>
                                    ))}
                                </ul>

                                <div style={{
                                    marginTop: '20px',
                                    padding: '12px',
                                    background: 'rgba(255,255,255,0.05)',
                                    borderRadius: '8px',
                                    display: 'flex',
                                    gap: '12px',
                                    alignItems: 'start'
                                }}>
                                    <Lightbulb size={20} color="var(--status-warning)" style={{ flexShrink: 0, marginTop: '2px' }} />
                                    <p style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', lineHeight: '1.4' }}>
                                        You can control your workflow remotely from your phone with just a few taps.
                                    </p>
                                </div>

                                <div style={{ marginTop: '20px' }}>
                                    <Button fullWidth onClick={() => setShowInstructions(false)} variant="secondary">
                                        Got it
                                    </Button>
                                </div>
                            </Card>
                        </motion.div>
                    </motion.div>
                )}
            </AnimatePresence>
        </>
    );
};

export default ProjectControl;
