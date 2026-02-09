import React from 'react';
import { motion } from 'framer-motion';
import { CheckCircle2 } from 'lucide-react';
import Button from '../components/UI/Button';
import Card from '../components/UI/Card';

const Login = ({ onLogin }) => {
    return (
        <div style={{
            padding: '0 24px',
            height: '100%',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'space-between',
            paddingBottom: '40px',
            paddingTop: '60px'
        }}>
            <motion.div
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ delay: 0.1 }}
            >
                <h1 style={{
                    fontSize: '2.5rem',
                    fontWeight: 800,
                    letterSpacing: '-0.02em',
                    marginBottom: '8px',
                    background: 'linear-gradient(to right, #fff, #a3a3a3)',
                    WebkitBackgroundClip: 'text',
                    WebkitTextFillColor: 'transparent',
                    textAlign: 'center'
                }}>
                    Select
                </h1>
                <p style={{
                    textAlign: 'center',
                    color: 'var(--text-secondary)',
                    fontSize: '1rem'
                }}>
                    Control your workflow anywhere.
                </p>
            </motion.div>

            <motion.div
                initial={{ scale: 0.95, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ delay: 0.2 }}
                className="flex-center"
                style={{ flexDirection: 'column', gap: '24px', flex: 1, justifyContent: 'center' }}
            >
                <div style={{ position: 'relative' }}>
                    <div style={{
                        position: 'absolute',
                        inset: '-20px',
                        background: 'radial-gradient(circle, var(--accent-glow) 0%, transparent 70%)',
                        zIndex: -1,
                        opacity: 0.5
                    }} />
                    <Card className="login-card" style={{ width: '100%' }}>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', alignItems: 'center' }}>
                            <div
                                style={{
                                    width: '64px',
                                    height: '64px',
                                    borderRadius: '50%',
                                    overflow: 'hidden',
                                    border: '2px solid rgba(255,255,255,0.1)'
                                }}
                            >
                                <img
                                    src="https://api.dicebear.com/7.x/avataaars/svg?seed=rohan&backgroundColor=262626"
                                    alt="User Avatar"
                                    style={{ width: '100%', height: '100%' }}
                                />
                            </div>
                            <div style={{ textAlign: 'center' }}>
                                <p style={{ fontSize: '0.875rem', color: 'var(--text-secondary)' }}>Logged in as</p>
                                <p style={{ fontWeight: 600 }}>rohan@example.com</p>
                            </div>
                            <div style={{
                                display: 'flex',
                                alignItems: 'center',
                                gap: '8px',
                                background: 'rgba(34, 197, 94, 0.1)',
                                padding: '8px 16px',
                                borderRadius: '20px',
                                color: 'var(--status-success)'
                            }}>
                                <div style={{
                                    width: '8px',
                                    height: '8px',
                                    background: 'var(--status-success)',
                                    borderRadius: '50%',
                                    boxShadow: '0 0 10px var(--status-success)'
                                }} />
                                <span style={{ fontSize: '0.875rem', fontWeight: 600 }}>Agent Online</span>
                                <CheckCircle2 size={16} />
                            </div>
                        </div>
                    </Card>
                </div>
            </motion.div>

            <motion.div
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ delay: 0.3 }}
            >
                <Button fullWidth onClick={onLogin} variant="primary">
                    Continue to Dashboard
                </Button>
            </motion.div>
        </div>
    );
};

export default Login;
