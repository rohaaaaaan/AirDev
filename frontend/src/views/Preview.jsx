import React from 'react';
import { motion } from 'framer-motion';
import { ExternalLink, Share2, CheckCircle2 } from 'lucide-react';
import Header from '../components/UI/Header';
import Button from '../components/UI/Button';
import Card from '../components/UI/Card';

const Preview = ({ onBack }) => {
    return (
        <>
            <Header title="Preview Ready" onBack={onBack} />

            <div style={{ padding: '0 20px', display: 'flex', flexDirection: 'column', alignItems: 'center', height: '100%', paddingTop: '40px' }}>
                <motion.div
                    initial={{ scale: 0.8, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    transition={{ type: 'spring', duration: 0.6 }}
                    style={{ marginBottom: '24px' }}
                >
                    <CheckCircle2 size={80} color="var(--status-success)" />
                </motion.div>

                <h2 style={{ fontSize: '2rem', fontWeight: 700, marginBottom: '8px' }}>Preview Ready!</h2>
                <a
                    href="#"
                    onClick={(e) => e.preventDefault()}
                    style={{ color: 'var(--accent-primary)', marginBottom: '40px', fontSize: '1rem', textDecoration: 'none' }}
                >
                    https://preview.app/demo
                </a>

                <div style={{ width: '100%', display: 'flex', flexDirection: 'column', gap: '16px' }}>
                    <Button fullWidth onClick={() => { }} className="flex-center" style={{ gap: '8px' }}>
                        <ExternalLink size={18} /> Open Preview
                    </Button>
                    <Button fullWidth variant="outline" onClick={() => { }} className="flex-center" style={{ gap: '8px' }}>
                        <Share2 size={18} /> Share
                    </Button>
                </div>
            </div>
        </>
    );
};

export default Preview;
